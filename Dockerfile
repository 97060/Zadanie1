FROM ubuntu:latest AS build
WORKDIR /app
COPY test.go ./
COPY go.mod ./
COPY upx-3.96-amd64_linux.tar.xz ./
COPY setup.sh ./
RUN apt-get update
RUN apt-get install wget
RUN wget --no-check-certificate https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
RUN apt-get install xz-utils
RUN ls -l
RUN tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz
RUN bash setup.sh
RUN export CGO_ENABLED=0 && /usr/local/go/bin/go build -ldflags="-s -w" test.go
RUN /usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./test

FROM scratch as main
COPY --from=build /app/test /test
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8082
CMD [ "/test" ]