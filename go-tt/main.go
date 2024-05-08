package main

import (
	"net"
	"time"
	"fmt"
	"flag"
	"golang.org/x/net/ipv4"
	"github.com/facebook/time/ptp/protocol"
)

func main() {
	port_interface_name := flag.String("brif", "eth1", "Interface of TT bridge outside port")
	fivegs_opponent_ip := flag.String("upip", "10.60.0.1", "IP of either UE or UPF where ptp packets will be forwarded")
	flag.Parse()
	TtListen(*port_interface_name, *fivegs_opponent_ip)
}

func TtListen(port_interface_name string, fivegs_opponent_ip string) {
	// Receives PTP messages via multicast 224.0.0.107 or 129 with port 319
	// Forward packets via 5GS or sends multicast to outside

	non_peer_msg_multicast_grp := net.IPv4(224, 0, 1, 129)
	peer_msg_multicast_grp := net.IPv4(224, 0, 0, 107)

	port_interface, err := net.InterfaceByName(port_interface_name)
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	multicast_conn, err := net.ListenPacket("udp4", "0.0.0.0:319")
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	fivegs_conn, err := net.Dial("udp", fivegs_opponent_ip + ":319")
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	// Join multicast groups
	multicast_packet_conn := ipv4.NewPacketConn(multicast_conn)
	if err := multicast_packet_conn.JoinGroup(port_interface, &net.UDPAddr{IP: non_peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := multicast_packet_conn.JoinGroup(port_interface, &net.UDPAddr{IP: peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Enable identification of multicast packets by destination
	// With this we can group the PTP messages into peer-delay and non-peer-delay types https://en.wikipedia.org/wiki/Precision_Time_Protocol
	if err := multicast_packet_conn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println(err.Error())
		return
	}

	defer multicast_conn.Close()
	defer multicast_packet_conn.Close()

	fmt.Println("TT: initialization complete")

	b := make([]byte, 1024)
	for {
		_, cm, _, err := multicast_packet_conn.ReadFrom(b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(peer_msg_multicast_grp) {
				fmt.Println("TT: received a peer delay ptp message")
				b = HandlePacket(b)
				fivegs_conn.Write(b)
				if err != nil { 
					fmt.Println(err.Error())
					return
				}
			} else if cm.Dst.Equal(non_peer_msg_multicast_grp) {
				fmt.Println("TT: received a non-peer-delay ptp message")
				b = HandlePacket(b)
				fivegs_conn.Write(b)
				if err != nil { 
					fmt.Println(err.Error())
					return
				}
			} else {
				fmt.Println("TT: received packet with unknown multicast group dropping it")
			}
		} else {
			fmt.Println("TT: received non-multicast packet")
			b = HandlePacket(b)
			_, err := fivegs_conn.Write(b)
			if err != nil { 
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func HandlePacket(raw_pkt []byte) ([]byte) {
	// Timestamp incoming PTP packets in downstream direction and calculate correction field for upstream direction
	// We assume sync packets only flow downstream and delay requests only upstream (user equipment is always ptp-slave)

	// Attempt to parse possible PTP packet
	parsed_pkt, err := protocol.DecodePacket(raw_pkt)

	// If parsing fails do nothing
	if err != nil { return raw_pkt }

	// Type switch into ptp packet types
	switch type_ptr := parsed_pkt.(type) {
		case *protocol.SyncDelayReq: {
			if parsed_pkt.MessageType() == protocol.MessageSync {
				fmt.Println("TT: parsed a sync packet")
				pkt := AppendIngressTimestamp(*type_ptr)
				raw_pkt, err = pkt.MarshalBinary()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}

	return raw_pkt 
}

func AppendIngressTimestamp(pkt protocol.SyncDelayReq) (protocol.SyncDelayReq) {
	// We hijack the 64bit correction factor field for temporarily storing the ingress time
	// In the correction field we store the time elapsed since the origin timestamp
	// then we calculate the difference at the egress port for the residence time
	// TODO: This makes it impossible to chain different bridges and accumulate corrections
	ns_since_origin_timestamp := time.Since(pkt.SyncDelayReqBody.OriginTimestamp.Time()).Nanoseconds()
	pkt.Header.CorrectionField = protocol.NewCorrection(float64(ns_since_origin_timestamp))
	return pkt
}