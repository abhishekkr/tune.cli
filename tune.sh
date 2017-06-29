#!/usr/bin/env bash

RUN_AT_DIR=$(pwd)
cd $(dirname $0)
THIS_DIR=$(pwd)
cd "${RUN_AT_DIR}"

[[ -z "${TUNE_SH_DIR}" ]] && export TUNE_SH_DIR="${THIS_DIR}"

source "${THIS_DIR}/helpers.sh"
source "${THIS_DIR}/tunefind.com/tunefind.sh"

tunefind-search-songs luci

