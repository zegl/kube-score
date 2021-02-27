FROM zegl/kube-score:v1.10.0 as kube-score

FROM alpine:3.13.2
COPY --from=kube-score /kube-score /
