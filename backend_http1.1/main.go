package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//go func() {
	err := http.ListenAndServe(":8081", http.HandlerFunc(h))
	if err != nil {
		log.Fatal(err)
	}
	//}()
	//err := http.ListenAndServeTLS(":8443", "../server.pem", "../server-key.pem", http.HandlerFunc(h))
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func h(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	fmt.Println(req.ContentLength)
	b := make([]byte, 1)
	i := 0

	for {
		n, err := req.Body.Read(b)
		if err != nil {
			if n > 0 {
				i++
				fmt.Printf("Byte %d : %q\n", i, b)
			}
			fmt.Println(err)
			break
		}
		i++
		fmt.Printf("Byte %d : %q\n", i, b)
	}
}
