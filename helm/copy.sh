#! /bin/bash

set -euo pipefail
set -x
dir="$(dirname $0)"
rm -r $dir/includes/nats
cp -r $dir/../scripts/nats $dir/includes/nats