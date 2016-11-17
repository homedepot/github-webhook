package web

import (
	"strconv"
	"net"
	"log"
	"net/http"
)

func Listen(port int) {

	sock, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatal(err)
	}

	//start listening
	go func() {
		err := http.Serve(sock, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
