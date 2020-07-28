FROM debian:stretch as downloader

ARG KUSTOMIZE_VERSION=v3.8.1
ARG KUSTOMIZE_SHA256SUM="9d5b68f881ba89146678a0399469db24670cba4813e0299b47cb39a240006f37"

RUN apt-get update && \
    apt-get install -y curl && \
    curl --location "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_amd64.tar.gz" > kustomize.tar.gz && \
    echo "${KUSTOMIZE_SHA256SUM}  kustomize.tar.gz" | sha256sum --check && \
    tar xzvf kustomize.tar.gz && \
    chmod +x kustomize

FROM alpine:3.10.1
RUN apk update && \
    apk upgrade && \
    apk add bash ca-certificates git
COPY --from=downloader kustomize /usr/bin/kustomize
COPY kube-score /usr/bin/kube-score
