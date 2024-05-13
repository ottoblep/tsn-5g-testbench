FROM ubuntu:focal
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade --yes && \
    DEBIAN_FRONTEND=noninteractive apt-get install --yes \
       inetutils-ping \
       iproute2 \
       tcpreplay \
       xxd

COPY ./*.pcap /