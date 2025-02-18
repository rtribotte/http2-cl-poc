package main

import (
	"fmt"
	"net/http"
)

func main() {
	go http.ListenAndServe(":8081", http.HandlerFunc(h))
	http.ListenAndServeTLS(":8443", "", "", http.HandlerFunc(h))
}

func h(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	fmt.Println(req.ContentLength)
	b := make([]byte, 1)
	i := 0

	for {
		_, err := req.Body.Read(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		i++
		fmt.Printf("Byte %d : %q\n", i, b)
	}
}
