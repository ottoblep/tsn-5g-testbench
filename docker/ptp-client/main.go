package main

import (
	"flag"
	"fmt"
	"github.com/facebook/time/ptp/simpleclient"
	"time"
)

func main() {
	server_addr := flag.String("serv_ip", "10.200.202.100", "PTP-Server Address (in this case forwarded by the bridge)")
	ifname := flag.String("if", "eth0", "interface name to send/receive packets")
	flag.Parse()

	cfg := &simpleclient.Config{
		Address: *server_addr,
		Iface: *ifname,
		Timeout: 1 * time.Minute,
		Duration: 30 * time.Minute,
		Timestamping: 0, // = Software https://pkg.go.dev/github.com/facebook/time/timestamp#Timestamp
	}

	client := simpleclient.New(cfg, displayResult)
	err := client.Run()
	fmt.Println(err.Error())
}

func displayResult(result *simpleclient.MeasurementResult) {
	fmt.Println("Delay: %s", result.Delay)
	fmt.Println("Offset: %s", result.Offset)
	fmt.Println("ServerToClientDiff: %s", result.ServerToClientDiff)
	fmt.Println("ClientToServerDiff: %s", result.ClientToServerDiff)
	fmt.Println("Timestamp: %s", result.Timestamp)
}