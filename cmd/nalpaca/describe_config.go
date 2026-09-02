package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/caarlos0/env/v11"
)

// descTag is the struct tag fields can carry, alongside their `env` tag, to document what
// the variable does. caarlos0/env has no concept of this itself - it only ever reads
// `env`/`envPrefix`/`envDefault` (see collectFieldMeta below) - so it's a convention of
// this codebase, not the library's.
const descTag = "desc"

// sensitiveTag marks a field's value as something that shouldn't be echoed back in plain
// text by default (API keys, secrets, passwords) - another codebase convention layered on
// top of the same fields, unrelated to anything caarlos0/env itself understands.
const sensitiveTag = "sensitive"

// describeConfig walks cfg (a pointer to a struct, generally the app's config struct)
// using the same struct tags that env.Parse uses to actually populate it, and writes a
// human-readable table of every environment variable the application understands to w.
//
// This is a purely static report driven by struct tags - it never reads the real
// environment, so it's safe to call with a zero-value config and doesn't require any of
// the variables it reports on to be set.
func describeConfig(w io.Writer, cfg any) error {
	fields, err := env.GetFieldParams(cfg)
	if err != nil {
		return fmt.Errorf("walking config struct: %w", err)
	}

	meta := collectFieldMeta(cfg)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ENV VAR\tREQUIRED\tDEFAULT\tDESCRIPTION\tNOTES")

	for _, f := range fields {
		def := "-"
		switch {
		case !f.HasDefaultValue:
		case f.DefaultValue == "":
			def = `""`
		default:
			def = f.DefaultValue
		}

		desc := meta[f.Key].desc
		if desc == "" {
			desc = "-"
		}

		required := "✅"
		if !f.Required {
			required = "❌"
		}

		fmt.Fprintf(tw, "%s\t%v\t%s\t%s\t%s\n", f.Key, required, def, desc, describeNotes(f))
	}

	return tw.Flush()
}

// fieldMeta holds the codebase-convention metadata attached to a config field alongside its
// `env` tag - see descTag and sensitiveTag.
type fieldMeta struct {
	desc      string
	sensitive bool
}

// collectFieldMeta walks cfg's type the same way env.Parse/env.GetFieldParams walk it -
// accumulating `envPrefix` as it descends into nested structs - but instead of collecting
// parse state, it collects each leaf field's desc/sensitive tags, keyed by the same
// fully-qualified env var name (prefix + `env` key) that GetFieldParams reports as
// FieldParams.Key. That shared key is what lets callers join the two together.
//
// env.GetFieldParams can't be reused for this directly: it only ever tracks the tags it
// already knows about (env/envPrefix/envDefault and the env tag's own comma-options), and
// has no extension point for an arbitrary custom tag like desc or sensitive, nor does it
// expose the reflect.StructField a caller would need to read one back. So this is a small
// parallel walk, driven by the type alone (no instance/env values needed).
func collectFieldMeta(cfg any) map[string]fieldMeta {
	t := reflect.TypeOf(cfg)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	out := make(map[string]fieldMeta)
	walkFieldMeta(t, "", out)
	return out
}

func walkFieldMeta(t reflect.Type, prefix string, out map[string]fieldMeta) {
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		ownKey, _, _ := strings.Cut(f.Tag.Get("env"), ",")
		if ownKey != "" && ownKey != "-" {
			if desc, sensitive := f.Tag.Get(descTag), f.Tag.Get(sensitiveTag) == "true"; desc != "" || sensitive {
				out[prefix+ownKey] = fieldMeta{desc: desc, sensitive: sensitive}
			}
		}

		ft := f.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			walkFieldMeta(ft, prefix+f.Tag.Get("envPrefix"), out)
		}
	}
}

// describeNotes renders the non-required/non-default field options (env tags like
// notEmpty, file, expand, and unset) as a human-readable summary.
func describeNotes(f env.FieldParams) string {
	var n []string
	if f.NotEmpty {
		n = append(n, "must not be empty")
	}
	if f.LoadFile {
		n = append(n, "value is a path to a file whose contents are used")
	}
	if f.Expand {
		n = append(n, "supports ${VAR} expansion")
	}
	if f.Unset {
		n = append(n, "process unsets this from its environment after reading it")
	}

	if len(n) == 0 {
		return "-"
	}

	return strings.Join(n, "; ")
}
