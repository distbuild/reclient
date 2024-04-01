#!/bin/bash

# Install Bazel
sudo apt install apt-transport-https curl gnupg -y
curl -L https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64 -o bazelisk
chmod +x bazelisk
sudo mv bazelisk /usr/local/bin

# Install depot_tools
git clone https://chromium.googlesource.com/chromium/tools/depot_tools --depth=1
