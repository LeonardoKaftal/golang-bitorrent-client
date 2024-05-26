package peer

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func UnmarshallPeer(peerBinary []byte) (*[]Peer, error) {
	const peerSize = 6
	if len(peerBinary)%peerSize != 0 {
		return nil, fmt.Errorf("received malformatted peer, invalid peer number")
	}
	peerNumber := len(peerBinary) / peerSize
	peers := make([]Peer, peerNumber)
	for i := 0; i < len(peerBinary); i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peerBinary[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peerBinary[offset+4 : offset+6])
	}
	return &peers, nil
}
