package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"quic-experiment/internal"
)

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0,0,0,0),
		Port: 2800,
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Listening on port 2800")

	defer func(){
		if err := listener.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		connection, err := listener.AcceptTCP()
		log.Printf("Received connection from %s", connection.RemoteAddr())
		if err != nil {
			log.Fatal(err)
		}

		go processConnection(connection)
	}
}

func processConnection(connection *net.TCPConn) {
	defer func(){
		if err := connection.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	remoteAddr := connection.RemoteAddr()

	payloadLenBuffer := make([]byte, 8)
	currentBytesCount := uint64(0)
	for currentBytesCount != 8 {
		n, err := connection.Read(payloadLenBuffer[currentBytesCount:])
		if err != nil{
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	payloadLen := internal.ExtractPayloadLen(payloadLenBuffer)
	payload := make([]byte, payloadLen)
	currentBytesCount = 0
	for currentBytesCount != payloadLen {
		n, err := connection.Read(payload[currentBytesCount:])
		if err != nil{
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
		n, err := connection.Write(response[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	log.Printf("Responded to %s", remoteAddr)
}