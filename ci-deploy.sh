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

# Login to the docker registry
docker login -u gitlab-ci-token -p "$CI_BUILD_TOKEN" "$CI_REGISTRY"

# Move to CMD folder to build binaries
cd "$CI_PROJECT_DIR"/cmd || exit

## Compile API Server
if mustCompile "api"; then
  echo "Deploying Server API"
  cd ./server-api || exit
  docker build --pull -t "$CI_REGISTRY_IMAGE"/server-api:"$IMAGE_TAG" .
  docker push "$CI_REGISTRY_IMAGE"/server-api:"$IMAGE_TAG"
  cd ..
fi


