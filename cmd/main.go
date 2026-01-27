package main

import (
	"fmt"

	"github.com/dafraer/sentence-gen-grpc-server/server"
)

func main() {
	fmt.Println("Hello World")
	srv := server.NewServer()
	srv.Run()
}
