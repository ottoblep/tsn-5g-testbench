FROM alpine:3.19.1
COPY ./src /src
RUN apk add make build-base
RUN make -f /src/Makefile