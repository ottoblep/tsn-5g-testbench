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

	// Internal 5GS connection
	// IP port 50000 is arbitrarily chosen to communicate between UE and UPF because the multicast is bound to 319
	fivegs_opponent_addr, err := net.ResolveUDPAddr("udp", gtp_tun_ip_opponent+":50000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fivegs_addr, _ := net.ResolveUDPAddr("udp", ":50000")

	fivegs_conn, err := net.ListenUDP("udp4", fivegs_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unicast Outside Connections
	// event_unicast_addr, err := net.ResolveUDPAddr("udp", ":319")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// general_unicast_addr, err := net.ResolveUDPAddr("udp", ":320")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// event_unicast_conn , err := net.ListenUDP("udp", event_unicast_addr)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// general_unicast_conn , err := net.ListenUDP("udp", general_unicast_addr)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// Multicast Outside Connections
	peer_event_addr, _ := net.ResolveUDPAddr("udp", "224.0.0.107:319")
	peer_general_addr, _ := net.ResolveUDPAddr("udp", "224.0.0.107:320")
	non_peer_event_addr, _ := net.ResolveUDPAddr("udp", "224.0.1.129:319")
	non_peer_general_addr, _ := net.ResolveUDPAddr("udp", "224.0.1.129:320")

	peer_general_multicast_conn, err := net.ListenMulticastUDP("udp", port_interface, peer_general_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	peer_event_multicast_conn, err := net.ListenMulticastUDP("udp", port_interface, peer_event_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	non_peer_general_multicast_conn, err := net.ListenMulticastUDP("udp", port_interface, non_peer_general_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	non_peer_event_multicast_conn, err := net.ListenMulticastUDP("udp", port_interface, non_peer_event_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer fivegs_conn.Close()
	defer peer_general_multicast_conn.Close()
	defer peer_event_multicast_conn.Close()
	defer non_peer_event_multicast_conn.Close()
	defer non_peer_general_multicast_conn.Close()
	// defer event_unicast_conn.Close()
	// defer general_unicast_conn.Close()

	fmt.Println("TT: initialization complete")

	// For some reason both multicast connections pick up all multicast packets (peer and non-peer) instead of only their group as specified in https://pkg.go.dev/net#ListenMulticastUDP
	// As such one listener is sufficient per port
	go ListenIncoming(non_peer_general_multicast_conn, fivegs_conn, fivegs_opponent_addr)
	go ListenIncoming(non_peer_event_multicast_conn, fivegs_conn, fivegs_opponent_addr)
	go ListenOutgoing(fivegs_conn,
		peer_general_multicast_conn, peer_event_multicast_conn,
		non_peer_general_multicast_conn, non_peer_event_multicast_conn,
		peer_general_addr, peer_event_addr,
		non_peer_general_addr, non_peer_event_addr,
	)

	// TODO: Could use a WaitGroup instead of loop
	for {
		time.Sleep(5 * time.Second)
	}
}

func ListenIncoming(multicast_conn *net.UDPConn, fivegs_conn *net.UDPConn, fivegs_opponent_addr *net.UDPAddr) {
	b := make([]byte, 1024)
	for {
		_, _, err := multicast_conn.ReadFrom(b)
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

func ListenOutgoing(fivegs_conn *net.UDPConn,
	peer_general_multicast_conn *net.UDPConn,
	peer_event_multicast_conn *net.UDPConn,
	non_peer_general_multicast_conn *net.UDPConn,
	non_peer_event_multicast_conn *net.UDPConn,
	peer_general_addr *net.UDPAddr,
	peer_event_addr *net.UDPAddr,
	non_peer_general_addr *net.UDPAddr,
	non_peer_event_addr *net.UDPAddr) {

	b := make([]byte, 1024)

	for {
		_, _, err := fivegs_conn.ReadFromUDP(b)

		msg_type, b := HandlePacket(false, b)

		switch msg_type {
		// Outgoing split by: Multicast or Unicast, port 320 or 319, multicast 0.107 or 1.129
		case protocol.MessageSync, protocol.MessageDelayReq: // Port 319 event, 224.0.1.129 non-peer
			{
				fmt.Println("TT: sending out non-peer event packet coming from 5gs bridge")
				_, err = non_peer_event_multicast_conn.WriteToUDP(b, non_peer_event_addr)
			}
		case protocol.MessagePDelayReq, protocol.MessagePDelayResp: // Port 319 event, 224.0.0.107 peer
			{
				fmt.Println("TT: sending out peer event packet coming from 5gs bridge")
				_, err = peer_event_multicast_conn.WriteToUDP(b, peer_event_addr)
			}
		case protocol.MessageAnnounce, protocol.MessageFollowUp, protocol.MessageDelayResp, protocol.MessageSignaling, protocol.MessageManagement: // Port 320 general, 224.0.1.129 non-peer
			{
				fmt.Println("TT: sending out non-peer general packet coming from 5gs bridge")
				_, err = non_peer_general_multicast_conn.WriteToUDP(b, non_peer_general_addr)
			}
		case protocol.MessagePDelayRespFollowUp: // Port 320 general, 224.0.0.107 peer
			{
				fmt.Println("TT: sending out peer general packet coming from 5gs bridge")
				_, err = peer_general_multicast_conn.WriteToUDP(b, peer_general_addr)
			}
		case 255:
			{
				fmt.Println("TT: dropping non-PTP packet")
			}
		default:
			{
				fmt.Println("TT: dropping unknown PTP packet type")
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
