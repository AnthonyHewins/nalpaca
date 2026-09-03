#! /bin/bash

set -euo pipefail

dir="$(dirname $0)"
ctx=""

help() { cat <<EOF
usage: $(basename $0) [FLAGS]

Initializes NATS for nalpaca. You'll need the NATS cli.
All flags valid for the NATS cli will work here:

NATS server urls (\$NATS_URL)
Username or Token (\$NATS_USER)
Password (\$NATS_PASSWORD)
Token (\$NATS_TOKEN)
User credentials (\$NATS_CREDS)
User NKEY (\$NATS_NKEY)
User JWT (\$NATS_JWT)
User seed (\$NATS_SEED)
TLS public certificate (\$NATS_CERT)
TLS private key (\$NATS_KEY)
TLS certificate authority chain (\$NATS_CA)
Time to wait on responses from NATS (\$NATS_TIMEOUT)
SOCKS5 proxy for connecting to NATS server (\$NATS_SOCKS_PROXY)
Sets a color scheme to use (\$NATS_COLOR)
Configuration context (\$NATS_CONTEXT)

Editing any of the JSON files this script uses in $dir to fit your
needs can be okay if you are editing things that don't bother nalpaca.
Settings that should not be touched are the stream names and the subjects
  
Editing subjects is not supported and will break your installation
  
  -h    Display help
  -c    Set the NATS context to use (or use env \$NATS_CONTEXT)
EOF
  exit $1
}

while getopts "hc:" flag; do
  case $flag in
  h) help 0;;
  c) ctx="$OPTARG";;
  \?) echo "Option not defined: $flag"; help 1;;
  esac
done

if [[ $(which nats) == "" ]]; then
  echo "Missing nats cli"
  help 1
fi

n="nats"
if [[ "$ctx" != "" ]]; then
  n+=" --context=$ctx"
fi

l="$dir/stream.json"
$n stream add nalpaca --config $l
find $dir -iname "*consumer.json" -exec $n consumer add nalpaca --config {} \;
$n kv add nalpaca