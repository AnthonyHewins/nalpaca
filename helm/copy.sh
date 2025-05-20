#! /bin/bash

set -euo pipefail
set -x
dir="$(dirname $0)"
cp -r $dir/../scripts/nats $dir/includes/nats