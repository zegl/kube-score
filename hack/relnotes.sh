#!/usr/bin/env bash

set -euo pipefail

# Dependencies: rg (ripgrep), jq

if [ "$#" -ne 1 ]; then
    echo "Unexpected number of arguments"
    echo
    echo "Usage:"
    echo "  $0 previous-release-tag"
    exit 1;
fi

PREV_RELEASE=$1
CURRENT_TAG=$(git tag -l --points-at HEAD);

if [ "$(echo "$CURRENT_TAG" | wc -l)" -ne 1 ]; then
    echo "Could not detect version. HEAD should have exactly one tag.";
    exit 1;
fi

#
# Generate list of changes based on RELNOTES in commits
#
echo "# Changes";
RELNOTE_MERGES=$(git log ${PREV_RELEASE}...HEAD --grep RELNOTE --oneline --merges)
while read -r line; do
    COMMIT=$(echo "$line" | awk '{print $1}')
    git show "$COMMIT" | rg -o '^\s+([0-9]+):(.*?)\s+RELNOTE:(.*?)\s+```' --multiline-dotall --multiline --replace "* #\$1 \$3" || true;
done <<< "$RELNOTE_MERGES"

#
# Authors secrion
#
echo
echo -n "This release contains contributions from: "
git log ${PREV_RELEASE}...HEAD | rg -o "Co-authored-by: (.*?) <" --replace "\$1" | sort |  uniq | awk 'ORS=", "' | sed 's/, $//'

#
# Download instructions
#
echo
echo "# Download"
echo "* Download the binaries from the GitHub release page"
echo "* Download the image from Docker Hub: \`zegl/kube-score:${CURRENT_TAG}\`"
echo "* Download the image from Docker Hub with Helm pre-installed: \`zegl/kube-score:${CURRENT_TAG}-helm\`"
echo "* Download from homebrew: \`brew install kube-score/tap/kube-score\`"
