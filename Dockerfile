FROM alpine
MAINTAINER Julor <julor@qq.com>

RUN mkdir -p /tmp/gocode/src/github.com/julor/godns

COPY .  /tmp/gocode/src/github.com/julor/godns

RUN set -x \
  && echo "https://mirrors.ustc.edu.cn/alpine/v3.3/main" > /etc/apk/repositories \
  && echo "https://mirrors.ustc.edu.cn/alpine/v3.3/community" >> /etc/apk/repositories \
  && apk update \
  && buildDeps='go git bzr' \
  && apk add --update $buildDeps \
  && GOPATH=/tmp/gocode GO15VENDOREXPERIMENT=1 go install github.com/julor/godns \
  && mkdir -p /usr/local/godns/log \
  && mkdir -p /usr/local/godns/conf \
  && mv /tmp/gocode/bin/godns /usr/local/godns/ \
  && mv /tmp/gocode/src/github.com/julor/godns/godns.conf /usr/local/godns/ \
  && chmod +x /usr/local/godns/godns \
  && apk del $buildDeps \
  && rm -rf /var/cache/apk/* /tmp/*

WORKDIR /usr/local/godns/

EXPOSE 53

ENTRYPOINT ["./godns","-c","godns.conf"]