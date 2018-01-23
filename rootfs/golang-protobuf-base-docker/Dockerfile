FROM golang:1.9-alpine
ENV PROTOBUF_TAG=v3.1.0
ADD build.sh /tmp/build.sh
RUN chmod +x /tmp/build.sh \
    && /tmp/build.sh \
    && /tmp/build.sh \
    && apk update \
    && apk add git