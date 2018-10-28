FROM golang:1.11-stretch as builder
WORKDIR /go/app
RUN apt-get install -y git
ADD . /go/app
RUN go build github.com/zegl/kube-score/cmd/kube-score

FROM scratch
COPY --from=builder /go/app/kube-score /kube-score