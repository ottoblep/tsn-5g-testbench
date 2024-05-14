# We reuse the base image from free5gc since it already includes a go build env and networking tools
FROM free5gc-base as ptp-client 

# Build PTP-Client 
COPY . /ptp-client
WORKDIR /ptp-client
RUN go build