# syntax=docker/dockerfile:1

ARG ALPINE_IMAGE=alpine:3.17.0

FROM ${ALPINE_IMAGE} as args
ARG HELM_VERSION=v3.10.2
ARG KUSTOMIZE_VERSION=v4.5.7

FROM args as args-amd64
ARG HELM_SHA256SUM="2315941a13291c277dac9f65e75ead56386440d3907e0540bf157ae70f188347"
ARG KUSTOMIZE_SHA256SUM="701e3c4bfa14e4c520d481fdf7131f902531bfc002cb5062dcf31263a09c70c9"

FROM args as args-arm64
ARG HELM_SHA256SUM="57fa17b6bb040a3788116557a72579f2180ea9620b4ee8a9b7244e5901df02e4"
ARG KUSTOMIZE_SHA256SUM="65665b39297cc73c13918f05bbe8450d17556f0acd16242a339271e14861df67"

# FROM args as args-arm
# ARG HELM_SHA256SUM="25af344f46348958baa1c758cdf3b204ede3ddc483be1171ed3738d47efd0aae"

FROM --platform=$BUILDPLATFORM args-${TARGETARCH} as downloader

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETARCH

RUN apk update && apk add curl

RUN curl --location "https://get.helm.sh/helm-${HELM_VERSION}-linux-${TARGETARCH}.tar.gz" > helm.tar.gz && \
    echo "${HELM_SHA256SUM}  helm.tar.gz" | sha256sum && \
    tar xzvf helm.tar.gz && \
    chmod +x /linux-${TARGETARCH}/helm && \
    mv  /linux-${TARGETARCH}/helm /usr/bin/helm

RUN curl --location "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_${TARGETARCH}.tar.gz" > kustomize.tar.gz && \
    echo "${KUSTOMIZE_SHA256SUM}  kustomize.tar.gz" | sha256sum && \
    tar xzvf kustomize.tar.gz && \
    chmod +x kustomize && \
    mv kustomize /usr/bin/kustomize

FROM golang:1.21-alpine as builder
ARG TARGETARCH
ARG TARGETPLATFORM
WORKDIR /kube-score
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG KUBE_SCORE_VERSION
ARG KUBE_SCORE_COMMIT
ARG KUBE_SCORE_DATE

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build \
    -ldflags="-X 'main.version=${KUBE_SCORE_VERSION}-docker-${TARGETPLATFORM}' -X 'main.commit=${KUBE_SCORE_COMMIT}' -X 'main.date=${KUBE_SCORE_DATE}'" \
    -o /usr/bin/kube-score \
    ./cmd/kube-score

FROM ${ALPINE_IMAGE} as runner
RUN apk update && \
    apk upgrade && \
    apk add bash ca-certificates git

COPY --from=downloader /usr/bin/helm /usr/bin/helm
COPY --from=downloader /usr/bin/kustomize /usr/bin/kustomize
COPY --from=builder /usr/bin/kube-score /usr/bin/kube-score

# Symlink to /kube-score for backwards compatibility (with kube-score v1.15.0 and earlier)
RUN ln -s /usr/bin/kube-score /kube-score

# Dry runs
RUN /kube-score version && \
    /usr/bin/kube-score version && \
    kube-score version && \
    helm version && kustomize version

WORKDIR /project
ENTRYPOINT ["/kube-score"]
