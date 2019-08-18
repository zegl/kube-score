FROM debian:stretch as downloader

ARG HELM_VERSION=v2.14.3
ARG HELM_SHA256SUM="38614a665859c0f01c9c1d84fa9a5027364f936814d1e47839b05327e400bf55"

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
