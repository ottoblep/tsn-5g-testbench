package main

import (
	"net"
	"time"
	"fmt"
	"golang.org/x/net/ipv4"
	"github.com/facebook/time/ptp/protocol"
)

func main() {
	TtListen()
}

func TtListen() {
	// Receives PTP messages via multicast 224.0.0.107 or 129 with port 319
	// Forward packets via 5GS or sends multicast to outside

	non_peer_msg_multicast_grp := net.IPv4(224, 0, 1, 129)
	peer_msg_multicast_grp := net.IPv4(224, 0, 0, 107)

	uplink_interface, err := net.InterfaceByName("eth1") // TODO: interface must not be hard-coded
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	uplink_multicast_conn, err := net.ListenPacket("udp4", "0.0.0.0:319")
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	downlink_conn, err := net.Dial("udp", "10.60.0.1:319") // TODO: UE ip must not be hard-coded
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	// Join multicast groups
	uplink_multicast_packet_conn := ipv4.NewPacketConn(uplink_multicast_conn)
	if err := uplink_multicast_packet_conn.JoinGroup(uplink_interface, &net.UDPAddr{IP: non_peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := uplink_multicast_packet_conn.JoinGroup(uplink_interface, &net.UDPAddr{IP: peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Enable identification of multicast packets by destination
	// With this we can group the PTP messages into peer-delay and non-peer-delay types https://en.wikipedia.org/wiki/Precision_Time_Protocol
	if err := uplink_multicast_packet_conn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println(err.Error())
		return
	}

	defer uplink_multicast_conn.Close()
	defer uplink_multicast_packet_conn.Close()

	fmt.Println("TT: initialization complete")

	b := make([]byte, 1024)
	for {
		_, cm, _, err := uplink_multicast_packet_conn.ReadFrom(b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(peer_msg_multicast_grp) {
				fmt.Println("TT: received a peer delay ptp message")
				b = HandlePacket(b)
				downlink_conn.Write(b)
				if err != nil { 
					fmt.Println(err.Error())
					return
				}
			} else if cm.Dst.Equal(non_peer_msg_multicast_grp) {
				fmt.Println("TT: received a non-peer-delay ptp message")
				b = HandlePacket(b)
				downlink_conn.Write(b)
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
			_, err := downlink_conn.Write(b)
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