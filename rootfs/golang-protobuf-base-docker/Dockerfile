FROM golang:1.9
ENV PROTOBUF_TAG=v3.5.1
ADD build.sh /tmp/build.sh
RUN chmod +x /tmp/build.sh \
    && apt-get update \
    && apt-get upgrade -y \
    && /tmp/build.sh \
