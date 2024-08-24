package main

import (
	"flag"
	"fmt"
	"github.com/facebook/time/ptp/protocol"
	"github.com/mdlayher/packet"
	"net"
	"sync"
	"time"
	"unsafe"
)

// Variables shared by the listeners
var (
	gtp_tun_opponent_addr_string string
	enable_unicast               bool
	enable_twostep               bool
	port_interface_name          string
	unicast_addr_string          string

	last_sync_residence_time           protocol.Correction
	last_sync_residence_time_mutex     sync.Mutex
	last_delayreq_residence_time       protocol.Correction
	last_delayreq_residence_time_mutex sync.Mutex
)

func main() {
	// The term "port" such as in "port_interface" refers to the outside connections of TSN bridge which normally are ethernet ports
	gtp_tun_opponent_addr_string_flag := flag.String("tunopip", "10.60.0.1", "IP of the other endpoint of the gtp tunnel where ptp packets will be forwarded to (in upstream direction there is no interface ip just the routing matters)")
	enable_twostep_flag := flag.Bool("twostep", false, "Switch operation from one step to two step")
	port_interface_name_flag := flag.String("portif", "eth1", "Interface of TT bridge outside port (only used with multicast)")
	flag.Parse()

	gtp_tun_opponent_addr_string = *gtp_tun_opponent_addr_string_flag
	enable_unicast = *enable_unicast_flag
	enable_twostep = *enable_twostep_flag
	port_interface_name = *port_interface_name_flag
	last_sync_residence_time = 0
	last_delayreq_residence_time = 0

	InitializeTT()
}

func InitializeTT() {
	// Setup Internal 5GS connection
	// IP port 38495 is chosen to communicate between UE and UPF because the multicast is bound to 319 and 320
	fivegs_addr, err := net.ResolveUDPAddr("udp", gtp_tun_opponent_addr_string+":38495")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fivegs_listen_addr, _ := net.ResolveUDPAddr("udp", ":38495")

	fivegs_conn, err := net.ListenUDP("udp4", fivegs_listen_addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer fivegs_conn.Close()

	port_interface, err := net.InterfaceByName(port_interface_name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 0x22F0 is TSN protocol type https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_ether.h
	raw_conn, err := packet.Listen(port_interface, 0, 0x22F0, &packet.Config{} )

	fmt.Println("TT: initialization complete")

	// TODO: Could use a WaitGroup instead of loop
	for {
		time.Sleep(5 * time.Second)
	}
}

func ListenIncoming() {

}

func ListenOutgoing() {

}

func HandlePacket(incoming bool, raw_pkt []byte) (protocol.MessageType, []byte) {
	// Act as transparent clock
	// We want to support both two step and one step transparent clock operation
	// so we both update the Sync/DelayRequest correction fields directly (1-step) and store the residence for a possible FollowUp or DelayResponse (2-step)
	// Peer to peer mode is not supported

	// Attempt to parse possible PTP packet
	parsed_pkt, err := protocol.DecodePacket(raw_pkt)
	zero_correction := protocol.NewCorrection(0)
	if err != nil {
		fmt.Println(err.Error())
		return 255, raw_pkt
	}

	// Type switch into ptp packet types
	switch pkt_ptr := parsed_pkt.(type) {
	case *protocol.SyncDelayReq:
		{
			correction := CalculateCorrection(incoming, (*pkt_ptr).Header.CorrectionField)

			if enable_twostep && !incoming {
				// In two step mode the follow up / delay response communicate the residence time
				(*pkt_ptr).Header.CorrectionField = zero_correction
			} else {
				(*pkt_ptr).Header.CorrectionField = correction
			}

			if !incoming {
				if (*pkt_ptr).Header.MessageType() == protocol.MessageSync {
					last_sync_residence_time_mutex.Lock()
					last_sync_residence_time = correction
					last_sync_residence_time_mutex.Unlock()
				} else {
					last_delayreq_residence_time_mutex.Lock()
					last_delayreq_residence_time = correction
					last_delayreq_residence_time_mutex.Unlock()
				}
			}
			raw_pkt, err = (*pkt_ptr).MarshalBinary()
		}
	case *protocol.FollowUp:
		{
			if !incoming {
				(*pkt_ptr).Header.CorrectionField = last_sync_residence_time
				raw_pkt, err = (*pkt_ptr).MarshalBinary()
			}
		}
	case *protocol.DelayResp:
		{
			if incoming {
				(*pkt_ptr).Header.CorrectionField = last_delayreq_residence_time
				raw_pkt, err = (*pkt_ptr).MarshalBinary()
			}
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return parsed_pkt.MessageType(), raw_pkt
}

func CalculateCorrection(incoming bool, correctionField protocol.Correction) protocol.Correction {
	// We hijack the 64bit correction field for temporarily storing the ingress time
	// Then we overwrite the elapsed time with the residence time at the egress port
	// Normally this is done by appending a suffix to the ptp message
	// TODO: This makes it impossible to chain different bridges and accumulate corrections

	if incoming {
		return UnixNanoToCorrection(time.Now().UnixNano())
	} else {
		residence_time := float64(time.Now().UnixNano() - CorrectionToUnixNano(correctionField))
		if residence_time <= 0 {
			fmt.Println("TT: calculated nonsense residence time ", residence_time, ", are the tt's clocks synchronized?")
			residence_time = 0
		}
		return protocol.NewCorrection(residence_time)
	}
}

// These functions do not convert between the two types! They are used to store a unix nanosecond in the header of the ptp message which has type protocol.Correction
// This unsafe casting works because the protocol.Correction type and unix nanoseconds types are both 64bit
func UnixNanoToCorrection(f int64) protocol.Correction {
	return *(*protocol.Correction)(unsafe.Pointer(&f))
}

func CorrectionToUnixNano(f protocol.Correction) int64 {
	return *(*int64)(unsafe.Pointer(&f))
}
