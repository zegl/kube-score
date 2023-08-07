# How to make a release

1. Build a release `GITHUB_TOKEN=ghp_XXXX goreleaser`
2. Build and push to Docker Hub: `PUSH_LATEST=1 ./hack/release-docker.sh`
3. Update the krew index: `./hack/krew-index.sh`
4. Update the relnotes
5. Update homebrew
