package main

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"quic-experiment/internal"

	"github.com/lucas-clemente/quic-go"
)

func main() {
	listener, err := quic.ListenAddr("0.0.0.0:2800", generateTLSConfig(), nil)
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

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-experiment"},
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
