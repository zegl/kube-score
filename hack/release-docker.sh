#!/bin/bash

set -euo pipefail

function echoinfo() {
	LIGHT_GREEN='\033[1;32m'
	NC='\033[0m' # No Color
	printf "${LIGHT_GREEN}%s${NC}\n" "$1"
}


VERSION="$(git describe --tags --abbrev=0)"
TAG_ARGS="-t zegl/kube-score:${VERSION}"
PLATFORM_ARGS=""
PUSH_OR_LOAD_ARG="--load"

if [ -z ${PUSH_LATEST+x} ]; then
    echoinfo "[x] Dry run. (Set PUSH_LATEST if you want to push and tag latest)"
else 
    echoinfo "[x] Making production build. Will push to Docker Hub!"
    TAG_ARGS="${TAG_ARGS} -t zegl/kube-score:latest"
    PLATFORM_ARGS="--platform linux/arm64 --platform linux/amd64"
    PUSH_OR_LOAD_ARG="--push"
    PUSH_OR_LOAD_ARG=""
fi

docker buildx build \
    --build-arg KUBE_SCORE_VERSION=${VERSION} \
    --build-arg "KUBE_SCORE_COMMIT=$(git rev-parse HEAD)" \
    --build-arg "KUBE_SCORE_DATE=$(date -Iseconds)" \
    ${PLATFORM_ARGS} \
    ${TAG_ARGS} \
    ${PUSH_OR_LOAD_ARG} \
    --target runner \
    .
