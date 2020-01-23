FROM golang:1.13-alpine3.10 as builder

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}
ARG GOSSPKS_VERSION=${GOSSPKS_VERSION:-master}
ARG GOSSPKS_COMMIT=

COPY . /src

WORKDIR /src

RUN apk add --update curl gcc build-base \
 && go get -v ./... \
 && go test -v ./... \
 && go build -ldflags "-s -w -X jdel.org/gosspks/cfg.Version=${GOSSPKS_VERSION}" \
 && chmod +x /src/gosspks

FROM jdel/alpine:3.10
LABEL maintainer=julien@del-piccolo.com

COPY --from=builder /src/gosspks /usr/local/bin/gosspks

RUN mkdir -p /home/user/gosspks/cache /home/user/gosspks/packages \
 && chown user:user /home/user/gosspks/cache /home/user/gosspks/packages

EXPOSE 8080

VOLUME ["/tmp/", "/home/user/gosspks/cache", "/home/user/gosspks/packages"]
 
CMD ["/usr/local/bin/gosspks"]
