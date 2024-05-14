# Unfortunately we need a newer go version than free5gc
FROM golang:1.21.10-bullseye AS ptp-server 

# Build PTP-Server 
COPY . /ptp-server
WORKDIR /ptp-server
RUN go build