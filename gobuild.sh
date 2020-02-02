#!/usr/bin/env bash

cmd_args=$@
cmd_args_cnt=$#
mustCompile() {
  if [[ 0 -eq ${cmd_args_cnt} ]]; then
    return 0
  fi
  for i in ${cmd_args}; do
    if [[ ${i} == $1 ]]; then
      return 0
    fi
  done
  return 1
}

here=$(pwd)

## Compile API Server
if mustCompile "api"; then
  echo "Building Server API"
  mkdir -p ./cmd/server-api/_build
  docker run \
    -v "${here}":/ronak \
    registry.ronaksoftware.com/base/docker/vips \
    /bin/bash -c "cd /ronak/cmd/server-api && go build -mod=vendor -a -ldflags '-s -w' -o ./_build/server-api ./"
  cd ./cmd/server-api/ || exit
  docker build -t registry.ronaksoftware.com/customers/ronakvision/server-blip/server-api:dev .

fi

