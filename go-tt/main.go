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
	port_interface_name := flag.String("portif", "eth1", "Interface of TT bridge outside port")
	gtp_tun_ip := flag.String("tunip", "10.100.200.137", "IP of this endpoint of the gtp tunnel (in upstream direction there is no interface ip just the packet destination matters)")
	gtp_tun_ip_opponent := flag.String("tunopip", "10.60.0.1", "IP of the other endpoint of the gtp tunnel where ptp packets will be forwarded to (in upstream direction there is no interface ip just the packet destination matters)")
	flag.Parse()
	TtListen(*port_interface_name, *gtp_tun_ip, *gtp_tun_ip_opponent)
}

func TtListen(port_interface_name string, gtp_tun_ip string, gtp_tun_ip_opponent string) {
	// Receives PTP messages via multicast 224.0.0.107 or 129 with ip port 319
	// Forward packets via 5GS or sends multicast to outside
	// port such as in "port_interface" refers to the outside connection points to the TSN bridge which normally are ethernet "ports"

	non_peer_msg_multicast_grp := net.IPv4(224, 0, 1, 129)
	peer_msg_multicast_grp := net.IPv4(224, 0, 0, 107)

	port_interface, err := net.InterfaceByName(port_interface_name)
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	port_conn, err := net.ListenPacket("udp4", ":319")
	if err != nil { 
		fmt.Println(err.Error())
		return
	}
	port_packet_conn := ipv4.NewPacketConn(port_conn)

	fivegs_send_conn, err := net.Dial("udp", gtp_tun_ip_opponent + ":319")
	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	// Join multicast groups
	if err := port_packet_conn.JoinGroup(port_interface, &net.UDPAddr{IP: non_peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := port_packet_conn.JoinGroup(port_interface, &net.UDPAddr{IP: peer_msg_multicast_grp}); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Enable identification of multicast packets by destination
	// With this we can group the PTP messages into peer-delay and non-peer-delay types https://en.wikipedia.org/wiki/Precision_Time_Protocol
	if err := port_packet_conn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println(err.Error())
		return
	}

	defer port_conn.Close()
	defer port_packet_conn.Close()
	defer fivegs_send_conn.Close()

	fmt.Println("TT: initialization complete")

	b := make([]byte, 1024)
	for {
		// Multicast packets come from the bridge port and are forwarded to UPF/UE via unicast
		_, cm, _, err := port_packet_conn.ReadFrom(b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(peer_msg_multicast_grp) {
				fmt.Println("TT: received a peer delay ptp message")
				_, b = HandlePacket(true, b)
				fivegs_send_conn.Write(b)
				if err != nil { 
					fmt.Println(err.Error())
					return
				}
			} else if cm.Dst.Equal(non_peer_msg_multicast_grp) {
				fmt.Println("TT: received a non-peer-delay ptp message")
				_, b = HandlePacket(true, b)
				fivegs_send_conn.Write(b)
				if err != nil { 
					fmt.Println(err.Error())
					return
				}
			} else {
				fmt.Println("TT: received packet with unknown multicast group dropping it")
			}
		} else {
			// Non-Multicast packets come from inside the 5gs and are forwarded to the outside again via multicast 

			// Prevent looping IP packets
			headerptr, err := ipv4.ParseHeader(b) 
			if err != nil { 
				fmt.Println(err.Error())
				continue
			}
			if !(*headerptr).Dst.Equal(net.ParseIP(gtp_tun_ip)) {
				fmt.Println("TT: dropping looping packet")
				continue
			}

			fmt.Println("TT: reading non-multicast packet")
			msg_type, b := HandlePacket(false, b)

			var dst *net.UDPAddr
			switch msg_type {
				case protocol.MessagePDelayReq, protocol.MessagePDelayResp, protocol.MessagePDelayRespFollowUp: {
					fmt.Println("TT: sending out packet coming from 5gs bridge")
					dst = &net.UDPAddr{IP: peer_msg_multicast_grp, Port: 319}
				}
				case 255: {
					fmt.Println("TT: dropping non-PTP packet")
					return
				}
				default: {
					fmt.Println("TT: sending out packet coming from 5gs bridge")
					dst = &net.UDPAddr{IP: non_peer_msg_multicast_grp, Port: 319}
				}
			}

			if err := port_packet_conn.SetMulticastInterface(port_interface); err != nil {
				fmt.Println(err.Error())
				return
			}

			port_packet_conn.SetMulticastTTL(2)
			if _, err := port_packet_conn.WriteTo(b, nil, dst); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func HandlePacket(incoming bool, raw_pkt []byte) (protocol.MessageType, []byte) {
	// Act as transparent clock 
	// We hijack the 64bit correction field for temporarily storing the ingress time
	// In the correction field we store the time elapsed since the origin timestamp
	// Then overwrite the elapsed time with the residence time at the egress port
	// TODO: This makes it impossible to chain different bridges and accumulate corrections

	// Attempt to parse possible PTP packet
	parsed_pkt, err := protocol.DecodePacket(raw_pkt)

	// If parsing fails do nothing
	if err != nil { return 255, raw_pkt }

	// Type switch into ptp packet types
	switch type_ptr := parsed_pkt.(type) {
		case *protocol.SyncDelayReq: {
			pkt := *type_ptr
			if parsed_pkt.MessageType() == protocol.MessageSync {
				if incoming {
					fmt.Println("TT: adding ingress timestamp to sync packet")
					ingress_ns_since_origin_timestamp := float64(time.Since(pkt.SyncDelayReqBody.OriginTimestamp.Time()).Nanoseconds())
					pkt.Header.CorrectionField = protocol.NewCorrection(ingress_ns_since_origin_timestamp)
				} else {
					fmt.Println("TT: calculating residence time for sync packet")
					egress_ns_since_origin_timestamp := float64(time.Since(pkt.SyncDelayReqBody.OriginTimestamp.Time()).Nanoseconds())
					ingress_ns_since_origin_timestamp := pkt.Header.CorrectionField.Nanoseconds()
					residence_time := egress_ns_since_origin_timestamp - ingress_ns_since_origin_timestamp
					if residence_time <= 0 {
						fmt.Println("TT: computed negative residence time, are the tt's clocks synchronized?")
						residence_time = 0
					}
					pkt.Header.CorrectionField = protocol.NewCorrection(residence_time)
				}
			}
			raw_pkt, err = pkt.MarshalBinary()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		case *protocol.PDelayReq: {
			pkt := *type_ptr
			if parsed_pkt.MessageType() == protocol.MessagePDelayReq {
				if incoming {
					fmt.Println("TT: adding ingress timestamp to peer delay packet")
					ingress_ns_since_origin_timestamp := float64(time.Since(pkt.PDelayReqBody.OriginTimestamp.Time()).Nanoseconds())
					pkt.Header.CorrectionField = protocol.NewCorrection(ingress_ns_since_origin_timestamp)

				} else {
					fmt.Println("TT: calculating residence time for peer delay packet")
					egress_ns_since_origin_timestamp := float64(time.Since(pkt.PDelayReqBody.OriginTimestamp.Time()).Nanoseconds())
					ingress_ns_since_origin_timestamp := pkt.Header.CorrectionField.Nanoseconds()
					residence_time := egress_ns_since_origin_timestamp - ingress_ns_since_origin_timestamp
					if residence_time <= 0 {
						fmt.Println("TT: computed negative residence time, are the tt's clocks synchronized?")
						residence_time = 0
					}
					pkt.Header.CorrectionField = protocol.NewCorrection(residence_time)
				}
			}
			raw_pkt, err = protocol.Bytes(&pkt) // TODO: generates wrong length field
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return parsed_pkt.MessageType(), raw_pkt 
}