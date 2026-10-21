package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/rpc"
	"os"

	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)

	client, err := rpc.Dial("tcp", *server)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	file, err := os.Open("wordlist")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		request := stubs.Request{Message: scanner.Text()}
		response := new(stubs.Response)
		if err := client.Call(stubs.PremiumReverseHandler, request, response); err != nil {
			panic(err)
		}
		fmt.Println(request.Message, "->", response.Message)
	}
}
