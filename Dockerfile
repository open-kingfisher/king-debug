FROM golang:1.14.3 as builder
ARG NAME="king-debug"
ARG GIT_URL="https://github.com/open-kingfisher/$NAME.git"
RUN git clone $GIT_URL /$NAME && cd /$NAME && make

FROM alpine:3.10

ARG NAME="king-debug"
COPY --from=builder /$NAME/entrypoint.sh /entrypoint.sh
COPY --from=builder /$NAME/bin/$NAME /usr/local/bin
COPY --from=builder /$NAME/bin/debug /usr/local/bin


ENTRYPOINT ["/bin/sh","/entrypoint.sh"]
