# FROM golang:latest as gobuilder
# WORKDIR /app
# COPY test.go ./
# COPY go.mod ./
# COPY setup.sh ./
# RUN uname -m
# RUN bash setup.sh
# RUN CGO_ENABLED=0 GOOS=linux  go build -ldflags="-w -s" -o server

# FROM ubuntu:latest AS build
# WORKDIR /app
# COPY upx-3.96-amd64_linux.tar.xz ./
# COPY --from=gobuilder /app/server ./
# RUN apt-get update
# RUN apt-get install xz-utils
# RUN tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz
# RUN ls -l
# RUN /usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./server

# FROM scratch as main
# COPY --from=build /app/server /
# ADD ca-certificates.crt /etc/ssl/certs/
# EXPOSE 8082
# ENTRYPOINT [ "/server" ]

FROM golang:1.18 as gobuilder
WORKDIR /app
COPY test.go ./
COPY go.mod ./
COPY setup.sh ./
COPY upx-3.96-amd64_linux.tar.xz ./
RUN bash setup.sh && \
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server && \
apt-get update && \
apt-get install xz-utils && \
tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz && \
/usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./server

FROM scratch as main
COPY --from=gobuilder /app/server /
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8082
ENTRYPOINT [ "/server" ]