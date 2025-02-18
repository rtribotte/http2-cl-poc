package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/quic-go/qpack"
	"github.com/quic-go/quic-go"
)

func main() {
	contentLength := 50
	dataFrames := []int{30, 10}

	// Configuration TLS
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h3"},
	}

	// Configuration QUIC
	quicConfig := &quic.Config{
		Versions: []quic.Version{quic.Version1},
	}

	// Établir une connexion QUIC
	ctx := context.Background()
	conn, err := quic.DialAddr(ctx, "localhost:8080", tlsConf, quicConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Ouvrir un stream QUIC
	stream, err := conn.OpenStream()
	if err != nil {
		log.Fatal(err)
	}

	// Préparer les en-têtes HTTP/3
	var headerBuf bytes.Buffer
	encoder := qpack.NewEncoder(&headerBuf)

	headers := []qpack.HeaderField{
		{Name: ":method", Value: "POST"},
		{Name: ":scheme", Value: "https"},
		{Name: ":authority", Value: "localhost"},
		{Name: ":path", Value: "/"},
		{Name: "content-length", Value: fmt.Sprintf("%d", contentLength)},
	}

	for _, header := range headers {
		encoder.WriteField(header)
	}

	// Écrire l'en-tête HTTP/3
	frameHeader := make([]byte, 2)
	frameHeader[0] = 0x01 // Type HEADERS
	frameHeader[1] = byte(headerBuf.Len())

	_, err = stream.Write(frameHeader)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stream.Write(headerBuf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	for i, frameSize := range dataFrames {
		data := randomBytes(frameSize)

		frameHeader := make([]byte, 2)
		frameHeader[0] = 0x00 // Type DATA
		frameHeader[1] = byte(len(data))

		_, err = stream.Write(frameHeader)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stream.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)

		if i == len(dataFrames)-1 {
			stream.Close()
		}
	}

	// Lire la réponse
	decoder := qpack.NewDecoder(func(f qpack.HeaderField) {
		fmt.Printf("Header reçu: %v\n", f)
	})

	for {
		// Lire l'en-tête de frame
		frameHeader := make([]byte, 2)
		_, err := io.ReadFull(stream, frameHeader)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		frameType := frameHeader[0]
		frameLen := int(frameHeader[1])

		// Lire le contenu de la frame
		frameData := make([]byte, frameLen)
		_, err = io.ReadFull(stream, frameData)
		if err != nil {
			log.Fatal(err)
		}

		switch frameType {
		case 0x01: // HEADERS
			decoder.Write(frameData)
		case 0x00: // DATA
			fmt.Printf("Données reçues: %s\n", string(frameData))
		}
	}
}

func randomBytes(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return b
}
