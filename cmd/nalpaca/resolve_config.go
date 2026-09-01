package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/caarlos0/env/v11"
)

// resolvedConfig runs cfg through the exact same env.Parse codepath the app itself uses at
// startup, against the real process environment, then prints what the fully resolved
// configuration looks like - so you can preview a real run's config without actually
// starting anything.
//
// env.Parse errors out the instant a required variable is missing, which is exactly what
// you don't want from a preview tool: you may be checking config on a box that doesn't have
// secrets wired up yet. So this tolerates a failed parse - env.Parse still populates every
// field it successfully can (see the doc comment on the loop in the caarlos0/env source:
// one field failing doesn't stop the others from being parsed) - and reports what did fail
// as warnings below the table instead of bailing out.
func resolvedConfig(w io.Writer, cfg any, showSecrets bool) error {
	fields, err := env.GetFieldParams(cfg)
	if err != nil {
		return fmt.Errorf("walking config struct: %w", err)
	}
	meta := collectFieldMeta(cfg)

	var warnings []string
	if perr := env.Parse(cfg); perr != nil {
		var agg env.AggregateError
		if !errors.As(perr, &agg) {
			// Not the kind of per-field error collection we know how to tolerate
			// (e.g. cfg wasn't a pointer to a struct at all) - this one's fatal.
			return fmt.Errorf("parsing config: %w", perr)
		}
		for _, e := range agg.Errors {
			warnings = append(warnings, e.Error())
		}
	}

	values := make(map[string]string)
	v := reflect.ValueOf(cfg)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	collectValues(v, "", values)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ENV VAR\tREQUIRED\tVALUE\tDEFAULT\tDESCRIPTION")

	for _, f := range fields {
		m := meta[f.Key]

		val := values[f.Key]
		switch {
		case f.Required && !isEnvSet(f.Key):
			val = "⚠ REQUIRED, NOT SET"
		case val == "":
			val = "-"
		case m.sensitive && !showSecrets:
			val = "•••• (redacted; pass -show-secrets to reveal)"
		}

		def := "-"
		switch {
		case !f.HasDefaultValue:
		case f.DefaultValue == "":
			def = `""`
		default:
			def = f.DefaultValue
		}

		desc := m.desc
		if desc == "" {
			desc = "-"
		}

		fmt.Fprintf(tw, "%s\t%v\t%s\t%s\t%s\n", f.Key, f.Required, val, def, desc)
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	if len(warnings) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "The app would fail to start as currently configured - these couldn't be resolved from the environment:")
		for _, warning := range warnings {
			fmt.Fprintln(w, "  -", warning)
		}
	}

	return nil
}

// isEnvSet reports whether key is actually present in the process environment, as opposed
// to being filled in from an envDefault. Used to tell "required and genuinely unset" apart
// from "required, but only reachable through env.Parse's aggregate error" - we already have
// FieldParams.Required, so we don't need to pick that fact back out of the parse error.
func isEnvSet(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// collectValues walks v (a populated config struct - the same instance env.Parse just ran
// against) and records each leaf field's resolved value as a display string, keyed the same
// way collectFieldMeta keys its map: fully-qualified env var name, accumulated via
// `envPrefix` as it descends.
func collectValues(v reflect.Value, prefix string, out map[string]string) {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := v.Field(i)

		ownKey, _, _ := strings.Cut(f.Tag.Get("env"), ",")
		if ownKey != "" && ownKey != "-" {
			out[prefix+ownKey] = formatValue(fv)
		}

		if fv.Kind() == reflect.Struct {
			collectValues(fv, prefix+f.Tag.Get("envPrefix"), out)
		}
	}
}

// formatValue renders a parsed field's value the way you'd want to read it in a config
// preview: symbol lists as a comma-joined string (matching how you'd actually set
// SYMBOLS=... on the way in), anything with a String() method (e.g. time.Duration) via that,
// and everything else via fmt's default verb.
func formatValue(v reflect.Value) string {
	switch vv := v.Interface().(type) {
	case []string:
		return strings.Join(vv, ",")
	case fmt.Stringer:
		return vv.String()
	default:
		return fmt.Sprintf("%v", vv)
	}
}
