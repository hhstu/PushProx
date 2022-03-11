FROM --platform=linux/arm64 golang:1.17.6 AS builder
WORKDIR /build
ARG TARGETARCH
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && cd cmd/proxy && GOOS=linux GOARCH=${TARGETARCH} go build  -ldflags "-s -w" . &&\
    cd ../client && GOOS=linux GOARCH=${TARGETARCH} go build  -ldflags "-s -w" .

FROM busybox:1.35.0-glibc

COPY --from=builder /build/cmd/proxy/proxy .
COPY --from=builder /build/cmd/client/client .
