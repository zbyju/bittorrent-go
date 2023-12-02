package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

type TorrentFile struct {
	Announce string `json:"announce"`
	Info     struct {
		Length      int    `json:"length"`
		Name        string `json:"name"`
		PieceLength int    `json:"piece length"`
		Pieces      string `json:"pieces"`
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decode.DecodeBencode(bencodedValue)
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
			return
		}

		decoded, err := decode.DecodeBencode(string(f))
		if err != nil {
			fmt.Println("Error decoding bencode:", err)
			return
		}

		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			return
		}

		var output TorrentFile
		err = json.Unmarshal([]byte(jsonOutput), &output)
		if err != nil {
			fmt.Println("Error unmarshalling from JSON:", err)
			return
		}

		fmt.Println("Tracker URL:", output.Announce)
		fmt.Println("Length:", output.Info.Length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
