#!/usr/bin/env bash

cmd_args=${PROJECTS_TO_BUILD[*]}
cmd_args_cnt=${#PROJECTS_TO_BUILD[*]}
mustCompile() {
  if [[ -z ${cmd_args} ]]; then
    return 0
  fi
  for i in ${cmd_args}; do
    if [[ ${i} == $1 ]]; then
      return 0
    fi
  done
  return 1
}

echo "${cmd_args}" "${cmd_args_cnt}"
echo "${PROJECTS_TO_BUILD}"

cd "$CI_PROJECT_DIR"/cmd || exit

## Compile API Server
if mustCompile "api"; then
  echo "Building Server API"
  mkdir -p "$CI_PROJECT_DIR"/cmd/server-api/_build
  GOOS=linux GOARCH=amd64 go build -mod=vendor -o "$CI_PROJECT_DIR"/cmd/server-api/_build/server-api ./server-api
fi

