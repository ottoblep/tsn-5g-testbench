# We reuse the base image from free5gc since it already includes a go build env and networking tools
FROM free5gc-base as ptp-server 

# Build PTP-Server 
COPY . /ptp-server
WORKDIR /ptp-server
RUN go build