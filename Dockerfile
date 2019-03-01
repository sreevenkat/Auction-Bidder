FROM alpine:3.7

RUN apk add --no-cache ca-certificates

RUN apk add --update \
        bash \
        curl \
    && rm -rf /var/cache/apk/*

ENV GOLANG_VERSION 1.9.3

ADD . .

CMD ["bash"]