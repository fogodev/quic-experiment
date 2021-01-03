package main

import (
	"fmt"
	"log"
	"net"
	"quic-experiment/internal"
	"time"
)

func main() {
	fileSizes := []internal.FileSize{internal.File100KB, internal.File1MB, internal.File10MB, internal.File100MB, internal.File1GB}
	allTimes := make([][]time.Duration, 5)
	for run := 0; run < 5; run++ {
		elapsedTimes := make([]time.Duration, len(fileSizes))
		for i, fileSize := range fileSizes {
			elapsedTimes[i] = sendData(fileSize, net.IPv4(127, 0, 0, 1), 2800)
		}
		allTimes[run] = elapsedTimes
	}

	fmt.Println("100KB,1MB,10MB,100MB,1GB")
	for run := 0; run < 5; run++ {
		fmt.Printf(
			"%v,%v,%v,%v,%v\n",
			allTimes[run][0],
			allTimes[run][1],
			allTimes[run][2],
			allTimes[run][3],
			allTimes[run][4],
		)
	}

}

func sendData(size internal.FileSize, ip net.IP, port int) time.Duration {
	payload, hash := internal.GenerateFileBytes(size)
	startTime := time.Now()

	connection, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := connection.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	message := internal.PrepareMessage(payload)
	messageLen := uint64(len(message))
	currentBytesCount := uint64(0)
	for currentBytesCount != messageLen {
		n, err := connection.Write(message[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	responseLenBuffer := make([]byte, 8)
	currentBytesCount = 0
	for currentBytesCount != 8 {
		n, err := connection.Read(responseLenBuffer[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	responseLen := internal.ExtractPayloadLen(responseLenBuffer)
	response := make([]byte, responseLen)
	currentBytesCount = 0
	for currentBytesCount != responseLen {
		n, err := connection.Read(response[currentBytesCount:])
		if err != nil {
			log.Fatal(err)
		}
		currentBytesCount += uint64(n)
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Done; Elapsed time (size %v kb): %s", size, elapsedTime)

	responseHash := string(response)
	log.Printf("\nReceived hash: %s\nExpected hash: %s\nMatched: %v", responseHash, hash, responseHash == hash)

	return elapsedTime
}
