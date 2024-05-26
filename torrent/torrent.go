package torrent

import (
	"io"
	"log"
	"main/peer"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int64
	Length      int64
	Name        string
}

const port uint16 = 6881

func bencodeToTorrentFile(bencode *Bencode) (*TorrentFile, error) {
	infoHash, err := bencode.GetInfoHash()
	if err != nil {
		return nil, err
	}
	pieceHashes, err := bencode.SplitPieceHashes()
	if err != nil {
		return nil, err
	}
	return &TorrentFile{
		Announce:    bencode.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bencode.Info.Length,
		Length:      bencode.Info.Length,
		Name:        bencode.Info.Name,
	}, nil
}

func ReadTorrentFile(torrentPath string) (*TorrentFile, error) {
	torrentBuff, err := os.ReadFile(torrentPath)
	if err != nil {
		log.Fatal("Error reading torrent file:", err)
	}
	bencode := UnmarshallBencode(torrentBuff)
	return bencodeToTorrentFile(bencode)
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.FormatInt(t.Length, 10)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) (*[]peer.Peer, error) {
	trackerURL, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Timeout: 15 * time.Second}
	bencodeTrackerResponse, err := httpClient.Get(trackerURL)
	if err != nil {
		return nil, err
	}
	defer bencodeTrackerResponse.Body.Close()
	body, err := io.ReadAll(bencodeTrackerResponse.Body)
	if err != nil {
		return nil, err
	}
	trackerResponse := UnmarshallTrackerBencodeResponse(body)
	return peer.UnmarshallPeer([]byte(trackerResponse.Peers))
}
