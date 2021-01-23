FROM debian:stretch as downloader

ARG HELM_VERSION=v3.5.0
ARG HELM_SHA256SUM="3fff0354d5fba4c73ebd5db59a59db72f8a5bbe1117a0b355b0c2983e98db95b"

RUN apt-get update && \
    apt-get install -y curl && \
    curl --location "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" > helm.tar.gz && \
    echo "${HELM_SHA256SUM}  helm.tar.gz" | sha256sum --check && \
    tar xzvf helm.tar.gz && \
    chmod +x /linux-amd64/helm

FROM alpine:3.10.1
RUN apk update && \
    apk upgrade && \
    apk add bash ca-certificates
COPY --from=downloader /linux-amd64/helm /usr/bin/helm
COPY kube-score /usr/bin/kube-score
