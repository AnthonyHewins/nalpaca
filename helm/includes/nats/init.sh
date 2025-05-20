#! /bin/bash

set -euo pipefail

dir="$(dirname $0)"
ctx=""

help() {
  echo "usage: $(basename $0) [FLAGS]"
  echo
  echo "Initializes NATS for nalpaca. You'll need the NATS cli."
  echo "Editing any of the JSON files this script uses in $dir to fit your"
  echo "needs is encouraged, as long as you keep the stream names the same and don't"
  echo "edit retention policies, as that changes behavior"
  echo 
  echo "Editing subjects is not supported and will break your installation"
  echo 
  echo "  -h    Display help"
  echo "  -c    Set the NATS context to use"
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
  stream=nalpaca-$(basename $i)-stream-v0 
  $n stream add $stream --config $i/stream.json # has to be ran first
  find $i -iname "*-consumer.json" -exec $n consumer add $stream  --config {} \;
done

nats kv add nalpaca