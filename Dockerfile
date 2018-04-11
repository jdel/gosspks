FROM jdel/alpine:edge

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}
ARG GOSSPKS_VERSION=${GOSSPKS_VERSION:-master}
ARG GOSSPKS_COMMIT=

LABEL maintainer=julien@del-piccolo.com

USER root

RUN apk add --update curl \
 && apk add --virtual build-dependencies go gcc build-base glide git openssh-client \
 && adduser gosspks -D \
 && mkdir -p /home/user/gosspks/packages /home/user/gosspks/cache \
 && chown -R user:user /tmp /home/user \
 && curl -sL https://github.com/jdel/gosspks/archive/${GOSSPKS_VERSION}.zip -o gosspks.zip \
 && mkdir -p ${GOPATH}/src/github.com/jdel/ \
 && unzip gosspks.zip -d ${GOPATH}/src/github.com/jdel/ \
 && rm -f gosspks.zip \
 && mv ${GOPATH}/src/github.com/jdel/gosspks-* ${GOPATH}/src/github.com/jdel/gosspks \
 && go get -v github.com/golang/dep/cmd/dep \
 && cd $GOPATH/src/github.com/golang/dep/cmd/dep \
 && git checkout tags/v0.4.1 && go install \
 && cd ${GOPATH}/src/github.com/jdel/gosspks/ \
 && dep ensure -v -vendor-only \
 && go build -o /usr/local/bin/gosspks -ldflags "-X github.com/jdel/gosspks/cfg.Version=${GOSSPKS_VERSION}-${GOSSPKS_COMMIT}" \
 && apk del build-dependencies \
 && rm -rf /var/cache/apk/* \
 && rm -rf /root/.glide/ \
 && rm -rf ${GOPATH}
 
USER user

WORKDIR /home/user/

EXPOSE 8080

VOLUME ["/tmp/", "/home/gosspks/gosspks/packages", "/home/gosspks/gosspks/cache"]
 
CMD ["/usr/local/bin/gosspks"]