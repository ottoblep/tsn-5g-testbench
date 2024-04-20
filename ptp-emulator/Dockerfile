FROM ubuntu:focal
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade --yes && \
    DEBIAN_FRONTEND=noninteractive apt-get install --yes \
       build-essential \
       git \
       python3-pip && \
    pip3 install --ignore-installed pyyaml
RUN git clone -b dev https://github.com/bestvibes/IEEE1588-PTP /bestvibes-ptp
WORKDIR /bestvibes-ptp
RUN make master_c slave_c