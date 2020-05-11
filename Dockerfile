FROM golang:alpine as build
WORKDIR /opt/k8s-ops-agent
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o app

FROM alpine:latest
COPY --from=build /opt/k8s-ops-agent/app /
CMD ["/app"]

