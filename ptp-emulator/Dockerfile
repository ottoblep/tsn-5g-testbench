FROM ubuntu:focal
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade --yes && \
    DEBIAN_FRONTEND=noninteractive apt-get install --yes \
       inetutils-ping \
       iproute2 \
       build-essential \
       git
       # python3-pip && \
    # pip3 install --ignore-installed pyyaml

RUN git clone -b dev https://github.com/bestvibes/IEEE1588-PTP /bestvibes-ptp
COPY ./*.patch /bestvibes-ptp

WORKDIR /bestvibes-ptp
RUN git reset --hard 6f3b38a2b293d771d241ba3e22c892b9347c9f0e
RUN git apply *.patch
RUN make master_c slave_c