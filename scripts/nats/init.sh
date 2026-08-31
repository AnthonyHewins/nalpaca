#! /bin/bash

set -euo pipefail

dir="$(dirname $0)"
ctx=""

help() { cat <<EOF
usage: $(basename $0) [FLAGS]

Initializes NATS for nalpaca. You'll need the NATS cli.
Editing any of the JSON files this script uses in $dir to fit your
needs can be okay if you are editing things that don't bother nalpaca.
Settings that should not be touched are the stream names and the subjects
  
Editing subjects is not supported and will break your installation
  
  -h    Display help
  -c    Set the NATS context to use
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

for i in $(find $dir -mindepth 1 -type d); do
  stream=nalpaca-$(basename $i)-stream 
  $n stream add $stream --config $i/stream.json # has to be ran first
  find $i -iname "*-consumer.json" -exec $n consumer add $stream  --config {} \;
done

$n kv add nalpaca