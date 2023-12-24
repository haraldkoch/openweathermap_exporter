FROM golang:alpine as builder

RUN apk add --no-cache git

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG VERSION

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    GOBIN=/go/bin \
    GO111MODULE=auto

WORKDIR /workspace

COPY . .

RUN \
  export GOOS \
  && GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) \
  && export GOARCH \
  && GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) \
  && export GOARM \
  && GOARM=$(echo ${TARGETPLATFORM} | cut -d / -f3 | cut -c2-) \
  && go get -d -v \
  && go build -o /go/bin/openweathermap_exporter -ldflags="-w -s"


FROM quay.io/prometheus/busybox:glibc

ARG BUILD_DATE
ARG VCS_REF

COPY --from=builder /go/bin/openweathermap_exporter /bin/openweathermap_exporter
EXPOSE      2112
USER        nobody
ENTRYPOINT  ["/bin/openweathermap_exporter"]

LABEL maintainer="Harald Koch <harald.koch@gmail.com>" \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${VCS_REF} \
      org.opencontainers.image.source="https://github.com/haraldkoch/openweathermap_exporter" \
      org.opencontainers.image.title="openweathermap_exporter" \
      org.opencontainers.image.version="${VERSION}"
