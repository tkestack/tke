#!/usr/bin/env bash

set -o xtrace

pwd

make build.web.console

cp -rv web/console/build "$DST_DIR/assets"
