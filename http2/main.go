package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func main() {
	contentLength := 50
	dataFrames := []int{30, 10}

	// Open a connection to inspect SettingsFrame.
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2"},
	})
	if err != nil {
		log.Fatal(err)
	}

	framer := http2.NewFramer(conn, conn)

	_, err = conn.Write([]byte(http2.ClientPreface))
	if err != nil {
		panic(err)
	}

	//_, err = framer.ReadFrame()
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = framer.WriteSettings()
	if err != nil {
		panic(err)
	}

	// Lire les réponses des paramètres du serveur (SETTINGS)
	settingsFrame, err := framer.ReadFrame()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received settings frame: %v\n", settingsFrame)

	// En-têtes de la requête POST avec Content-Length
	headerBuf := &bytes.Buffer{}
	hpackEncoder := hpack.NewEncoder(headerBuf)
	hpackEncoder.WriteField(hpack.HeaderField{Name: ":method", Value: "POST"})
	hpackEncoder.WriteField(hpack.HeaderField{Name: ":scheme", Value: "http"})
	hpackEncoder.WriteField(hpack.HeaderField{Name: ":authority", Value: "localhost"})
	hpackEncoder.WriteField(hpack.HeaderField{Name: ":path", Value: "/"})
	hpackEncoder.WriteField(hpack.HeaderField{Name: "content-length", Value: fmt.Sprintf("%d", contentLength)})

	dec := hpack.NewDecoder(4<<10, func(f hpack.HeaderField) {
		fmt.Println("header:", f.String())
	})

	// Envoyer les HEADERS
	err = framer.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		BlockFragment: headerBuf.Bytes(),
		EndHeaders:    true,
	})
	if err != nil {
		panic(err)
	}

	for i, frameSize := range dataFrames {
		end := false
		if i == len(dataFrames)-1 {
			end = true
		}
		err = framer.WriteData(1, end, randomBytes(frameSize))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	// Lire la réponse
	for {
		frame, err := framer.ReadFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fr, ok := frame.(*http2.HeadersFrame)
		if ok {
			dec.Write(fr.HeaderBlockFragment())
		}

		if _, ok := frame.(*http2.RSTStreamFrame); ok {
			fmt.Printf("Received RST frame: %v\n", frame)
			break
		}

		fmt.Printf("Received frame: %v\n", frame)
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
