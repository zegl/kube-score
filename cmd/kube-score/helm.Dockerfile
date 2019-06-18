FROM debian:stretch as downloader

ARG HELM_VERSION=v2.14.1
ARG HELM_SHA256SUM="804f745e6884435ef1343f4de8940f9db64f935cd9a55ad3d9153d064b7f5896"

RUN apt-get update && \
    apt-get install -y curl && \
    curl --location "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" > helm.tar.gz && \
    echo "${HELM_SHA256SUM}  helm.tar.gz" | sha256sum --check && \
    tar xzvf helm.tar.gz && \
    chmod +x /linux-amd64/helm

FROM alpine:3.4
RUN apk update && \
    apk upgrade && \
    apk add bash
COPY --from=downloader /linux-amd64/helm /usr/bin/helm
COPY kube-score /usr/bin/kube-score
ENTRYPOINT ["/kube-score"]
