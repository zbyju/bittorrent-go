package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/mybencode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

type TorrentFileInfo struct {
	Length      int    `json:"length"`
	Name        string `json:"name"`
	PieceLength int    `json:"piece length"`
	Pieces      string `json:"pieces"`
}

type TorrentFile struct {
	Announce string          `json:"announce"`
	Info     TorrentFileInfo `json:"info"`
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := mybencode.DecodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		if len(os.Args) != 3 {
			fmt.Println("Usage: info [path/to/file.torrent]")
			return
		}

		file_path := os.Args[2]
		f, err := os.ReadFile(file_path)

		if err != nil {
			fmt.Println("test")
			return
		}

		decoded, err := mybencode.DecodeBencode(string(f))
		if err != nil {
			fmt.Println("Error decoding bencode:", err)
			return
		}

		decodedMap := decoded.(map[string]interface{})

		jsonString, _ := json.Marshal(decodedMap)

		var tf TorrentFile
		json.Unmarshal(jsonString, &tf)

		encoded, err := mybencode.Encode(decodedMap["info"])

		if err != nil {
			return
		}

		h := sha1.New()
		h.Write([]byte(encoded))
		bs := h.Sum(nil)

		fmt.Printf("Tracker URL: %s\n", tf.Announce)
		fmt.Printf("Length: %d\n", tf.Info.Length)
		fmt.Printf("Info Hash: %x\n", bs)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
