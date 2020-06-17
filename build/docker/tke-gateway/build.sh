#!/usr/bin/env bash

set -o xtrace

pwd

make web.build.console

cp -rv web/console/build "$DST_DIR/assets"
