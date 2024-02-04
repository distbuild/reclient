#!/bin/bash

export PATH=$PWD/depot_tools:$PATH

# bootstrap: bazel-bin/cmd/bootstrap/bootstrap_/bootstrap
# dumpstats: bazel-bin/cmd/dumpstats/dumpstats_/dumpstats
# reproxy: bazel-bin/cmd/reproxy/reproxy_/reproxy
# rewrapper: bazel-bin/cmd/rewrapper/rewrapper_/rewrapper
bazel build --config=clangscandeps //cmd/...
