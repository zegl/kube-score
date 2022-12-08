#!/bin/bash

set -euo pipefail
set -x

VERSION=$(git describe --tags --abbrev=0 | cut -c2-)
KREW_INDEX_PATH="$(mktemp -d)"
FILE="${KREW_INDEX_PATH}/plugins/score.yaml"
KUBE_SCORE_REPO_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd )"

checksum() {
    grep -E "kube-score_${VERSION}_${1}_amd64.tar.gz" ${KUBE_SCORE_REPO_ROOT}/dist/checksums.txt | awk '{print $1}'
}

gg() {
    git -C "$KREW_INDEX_PATH" "$@"
}

git clone git@github.com:kubernetes-sigs/krew-index.git "$KREW_INDEX_PATH"
git -C  "$KREW_INDEX_PATH" remote add zegl git@github.com:zegl/krew-index.git

gg checkout master
gg fetch origin
gg reset --hard origin/master
gg branch -D "kube-score-${VERSION}" || true
gg checkout -b "kube-score-${VERSION}"

yq --inplace ".spec.version = \"v${VERSION}\""  "${FILE}"

yq --inplace ".spec.platforms[0].uri = \"https://github.com/zegl/kube-score/releases/download/v${VERSION}/kube-score_${VERSION}_darwin_amd64.tar.gz\"" "$FILE"
yq --inplace ".spec.platforms[0].sha256 = \"$(checksum darwin)\"" "$FILE"

yq --inplace ".spec.platforms[1].uri = \"https://github.com/zegl/kube-score/releases/download/v${VERSION}/kube-score_${VERSION}_linux_amd64.tar.gz\"" "$FILE"
yq --inplace ".spec.platforms[1].sha256 = \"$(checksum linux)\"" "$FILE"

gg add plugins/score.yaml
gg commit -m "Update score to v${VERSION}"
gg push --force -u zegl "kube-score-${VERSION}"

open "https://github.com/kubernetes-sigs/krew-index/compare/master...zegl:kube-score-${VERSION}?expand=1"
