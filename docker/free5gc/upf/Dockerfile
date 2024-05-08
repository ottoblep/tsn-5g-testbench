# Compared to the upf dockerfile from https://github.com/free5gc/free5gc-compose 
# this does not discard the build environment including the compiler and source code
# so it can be used for development

FROM free5gc/base

ENV DEBIAN_FRONTEND noninteractive

# Install UPF dependencies
RUN apt-get update \
    && apt-get install -y libmnl0 libyaml-0-2 iproute2 iptables strace net-tools iputils-ping curl netcat tcpdump\
    && apt-get clean

# Build Free5GC
RUN git clone --recurse-submodules https://github.com/free5gc/free5gc /free5gc

# Optionally Replace select network functions with forked versions
# RUN rm -rf free5gc/NFs/upf
# ADD go-upf free5gc/NFs/upf

WORKDIR /free5gc
RUN make upf 

RUN mkdir -p cert/ public

RUN git clone https://github.com/free5gc/go-gtp5gnl.git && mkdir "go-gtp5gnl/bin" && \
    cd "go-gtp5gnl/cmd/gogtp5g-tunnel" &&  go build -o "/free5gc/gtp5g-tunnel" . && \
    cd -

# Build NW-TT 
COPY ./go-tt /go-tt
WORKDIR /go-tt
RUN go build

# Config files volume
VOLUME [ "/free5gc/config" ]

