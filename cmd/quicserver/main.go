package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"quic-experiment/internal"

	"github.com/lucas-clemente/quic-go"
)

func main() {
	listener, err := quic.ListenAddr("0.0.0.0:2800", internal.GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Listening on port 2800")
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		session, err := listener.Accept(context.Background())
		log.Printf("Received connection from %s", session.RemoteAddr())
		if err != nil {
			log.Fatal(err)
		}

		go processStream(session)
	}
}

func processStream(session quic.Session) {
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	remoteAddr := session.RemoteAddr()
	payloadLenBuffer := make([]byte, 8)
	currentBytesCount := uint64(0)
	for currentBytesCount != 8 {
		n, err := stream.Read(payloadLenBuffer[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	payloadLen := internal.ExtractPayloadLen(payloadLenBuffer)
	payload := make([]byte, payloadLen)
	currentBytesCount = 0
	for currentBytesCount != payloadLen {
		n, err := stream.Read(payload[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	hash := fmt.Sprintf("%x", md5.Sum(payload))
	log.Printf("Hash from %s: %s", remoteAddr, hash)

	response := internal.PrepareMessage([]byte(hash))
	responseLen := uint64(len(response))
	currentBytesCount = 0
	for currentBytesCount != responseLen {
		n, err := stream.Write(response[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	log.Printf("Responded to %s", remoteAddr)
}
