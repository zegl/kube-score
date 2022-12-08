#!/bin/bash

VERSION="$(git describe --tags --abbrev=0)"

# TODO(gustav): also push latest!

docker buildx build \
    --build-arg KUBE_SCORE_VERSION=${VERSION} \
    --build-arg "KUBE_SCORE_COMMIT=$(git rev-parse HEAD)" \
    --build-arg "KUBE_SCORE_DATE=$(date -Iseconds)" \
    --platform linux/arm64 \
    --platform linux/amd64 \
    -t zegl/kube-score:${VERSION} \
    --push \
    --target runner \
    .
