# Unfortunately we need a newer go version than free5gc
FROM golang:1.21.10-bullseye AS ptp-client 

# Build PTP-Client 
COPY . /ptp-client
WORKDIR /ptp-client
RUN go build