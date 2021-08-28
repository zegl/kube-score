FROM debian:stretch as downloader

ARG HELM_VERSION=v3.6.3
ARG HELM_SHA256SUM="07c100849925623dc1913209cd1a30f0a9b80a5b4d6ff2153c609d11b043e262"

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
