FROM golang:1.18 as gobuilder
ARG DOCKER_PLATFORM=$TARGETPLATFORM
WORKDIR /app
COPY test.go ./
COPY go.mod ./
COPY setup.sh ./
RUN export CGO_ENABLED=0 && \
export GOOS=linux
RUN uname -m
RUN bash setup.sh
RUN /usr/local/go/bin/go build -ldflags="-s -w" test.go

FROM ubuntu:latest AS build
WORKDIR /app
COPY upx-3.96-amd64_linux.tar.xz ./
COPY --from=gobuilder /app/test /app/
RUN apt-get update
RUN apt-get install xz-utils
RUN tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz
RUN /usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./test

FROM scratch as main
COPY --from=build /app/test /test
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8082
CMD [ "/test" ]