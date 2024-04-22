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
COPY ./set_ptp_slave_ip.patch /bestvibes-ptp/set_ptp_slave_ip.patch
WORKDIR /bestvibes-ptp
RUN git apply set_ptp_slave_ip.patch
RUN make master_c slave_c