package main

import (
	"flag"
	"fmt"
	"github.com/facebook/time/ptp/protocol"
	"net"
	"time"
)

func main() {
	port_interface_name := flag.String("portif", "eth1", "Interface of TT bridge outside port")
	gtp_tun_ip_opponent := flag.String("tunopip", "10.60.0.1", "IP of the other endpoint of the gtp tunnel where ptp packets will be forwarded to (in upstream direction there is no interface ip just the routing matters)")
	flag.Parse()
	TtListen(*port_interface_name, *gtp_tun_ip_opponent)
}

func TtListen(port_interface_name string, gtp_tun_ip_opponent string) {
	// Receives PTP messages via multicast 224.0.0.107 or 224.0.1.129 with ip port 319
	// Forwards packets via 5GS or sends multicast to outside
	// Updates the correction field of PTP packets passing through the 5GS
	// The term "port" such as in "port_interface" refers to the outside connections of TSN bridge which normally are ethernet ports

	port_interface, err := net.InterfaceByName(port_interface_name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// IP port 50000 is arbitrarily chosen to communicate between UE and UPF because the multicast is bound to 319
	fivegs_opponent_addr, err := net.ResolveUDPAddr("udp", gtp_tun_ip_opponent+":50000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	local_address, err := net.ResolveUDPAddr("udp", ":50000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	peer_addr, err := net.ResolveUDPAddr("udp", "224.0.0.107:319")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	non_peer_addr, err := net.ResolveUDPAddr("udp", "224.0.1.129:319")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	peer_mc_conn, err := net.ListenMulticastUDP("udp", port_interface, peer_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	non_peer_mc_conn, err := net.ListenMulticastUDP("udp", port_interface, non_peer_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fivegs_conn, err := net.ListenUDP("udp4", local_address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer fivegs_conn.Close()
	defer peer_mc_conn.Close()
	defer non_peer_mc_conn.Close()

	fmt.Println("TT: initialization complete")

	// For some reason both multicast connections pick up all multicast packets instead of only their group as specified in https://pkg.go.dev/net#ListenMulticastUDP
	// As such one of them is sufficient for listening
	go ListenIncoming(non_peer_mc_conn, fivegs_conn, fivegs_opponent_addr)
	// go ListenIncoming(peer_mc_conn, fivegs_conn, fivegs_opponent_addr) 
	go ListenOutgoing(fivegs_conn, peer_mc_conn, non_peer_mc_conn, peer_addr, non_peer_addr)

	// TODO: Could use a WaitGroup instead of loop
	for {
		time.Sleep(5 * time.Second)
	}
}

func ListenIncoming(mc_conn *net.UDPConn, fivegs_conn *net.UDPConn, fivegs_opponent_addr *net.UDPAddr) {
	b := make([]byte, 1024)
	for {
		_, _, err := mc_conn.ReadFrom(b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("TT: received a multicast ptp message")
		_, b = HandlePacket(true, b)

		_, err = fivegs_conn.WriteToUDP(b, fivegs_opponent_addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func ListenOutgoing(fivegs_conn *net.UDPConn, peer_mc_conn *net.UDPConn, non_peer_mc_conn *net.UDPConn, peer_addr *net.UDPAddr, non_peer_addr *net.UDPAddr) {
	b := make([]byte, 1024)

	for {
		_, _, err := fivegs_conn.ReadFromUDP(b)

		msg_type, b := HandlePacket(false, b)

		switch msg_type {
		// Peer messages are sent to a different broadcast address
		case protocol.MessagePDelayReq, protocol.MessagePDelayResp, protocol.MessagePDelayRespFollowUp:
			{
				fmt.Println("TT: sending out peer packet coming from 5gs bridge")
				_, err = peer_mc_conn.WriteToUDP(b, peer_addr)
			}
		case 255:
			{
				fmt.Println("TT: dropping non-PTP packet")
			}
		default:
			{
				fmt.Println("TT: sending out non-peer packet coming from 5gs bridge")
				_, err = non_peer_mc_conn.WriteToUDP(b, non_peer_addr)
			}
		}

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func HandlePacket(incoming bool, raw_pkt []byte) (protocol.MessageType, []byte) {
	// Act as transparent clock

	// Attempt to parse possible PTP packet
	parsed_pkt, err := protocol.DecodePacket(raw_pkt)
	if err != nil {
		return 255, raw_pkt
	}

	// Type switch into ptp packet types
	switch pkt_ptr := parsed_pkt.(type) {
	case *protocol.SyncDelayReq:
		{
			(*pkt_ptr).Header.CorrectionField = CalculateCorrection(incoming, (*pkt_ptr).SyncDelayReqBody.OriginTimestamp, (*pkt_ptr).Header.CorrectionField)
			raw_pkt, err = (*pkt_ptr).MarshalBinary()
		}
	case *protocol.PDelayReq:
		{
			(*pkt_ptr).Header.CorrectionField = CalculateCorrection(incoming, (*pkt_ptr).PDelayReqBody.OriginTimestamp, (*pkt_ptr).Header.CorrectionField)
			raw_pkt, err = protocol.Bytes(&(*pkt_ptr)) // TODO: sometimes generates wrong length?
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return parsed_pkt.MessageType(), raw_pkt
}

func CalculateCorrection(incoming bool, originTimestamp protocol.Timestamp, correctionField protocol.Correction) protocol.Correction {
	// We hijack the 64bit correction field for temporarily storing the ingress time
	// In the correction field we store the time elapsed since the origin timestamp
	// Then we overwrite the elapsed time with the residence time at the egress port
	// TODO: This makes it impossible to chain different bridges and accumulate corrections

	ns_since_origin_timestamp := float64(time.Since(originTimestamp.Time()).Nanoseconds())

	if incoming {
		fmt.Println("TT: adding ingress timestamp")
		return protocol.NewCorrection(ns_since_origin_timestamp)
	} else {
		fmt.Println("TT: calculating residence time")
		ns_since_origin_timestamp_at_ingress := correctionField.Nanoseconds()
		residence_time := ns_since_origin_timestamp - ns_since_origin_timestamp_at_ingress
		if residence_time <= 0 {
			fmt.Println("TT: computed negative residence time, are the tt's clocks synchronized?")
			residence_time = 0
		}
		return protocol.NewCorrection(residence_time)
	}
}
