#!/bin/sh
set -eu
bun install
moon run :test
moon run :lint
moon run :build
