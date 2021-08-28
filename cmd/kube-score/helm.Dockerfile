FROM debian:stretch as downloader

ARG HELM_VERSION=v2.17.0
ARG HELM_SHA256SUM="f3bec3c7c55f6a9eb9e6586b8c503f370af92fe987fcbf741f37707606d70296"

RUN apt-get update && \
    apt-get install -y curl && \
    curl --location "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" > helm.tar.gz && \
    echo "${HELM_SHA256SUM}  helm.tar.gz" | sha256sum --check && \
    tar xzvf helm.tar.gz && \
    chmod +x /linux-amd64/helm

FROM alpine:3.14.2
RUN apk update && \
    apk upgrade && \
    apk add bash ca-certificates
COPY --from=downloader /linux-amd64/helm /usr/bin/helm
COPY kube-score /usr/bin/kube-score
