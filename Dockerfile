FROM golang:1.11-stretch as builder
WORKDIR /go/app
ENV CGO_ENABLED=0

ADD "go.*" /go/app/
ADD "*.go" /go/app/
ADD cmd /go/app/cmd/
ADD score /go/app/score/
ADD scorecard /go/app/scorecard/

RUN go build github.com/zegl/kube-score/cmd/kube-score

FROM scratch
COPY --from=builder /go/app/kube-score /kube-score
ENTRYPOINT ["/kube-score"]
