package main
import ("net/rpc"
	"flag"
	"time"
	"fmt"
	"uk.ac.bris.cs/distributed3/pairbroker/stubs"
	"math/rand")

func main(){
	brokerAddr := flag.String("broker","127.0.0.1:8030", "Address of broker instance")
	topic := flag.String("topic","multiply", "Topic this miner will publish to")
	flag.Parse()
	//Dial broker address.
	client, _ := rpc.Dial("tcp", *brokerAddr)
	status := new(stubs.StatusReport)
	//Create a new buffered channel
	client.Call(stubs.CreateChannel, stubs.ChannelRequest{Topic: *topic, Buffer: 10}, status)
	//Random seed for 'mining' work.
	rand.Seed(time.Now().UnixNano())
	for {
		//Create three new random integers (smaller range so X*Y*Z fits in an int)
		newtriplet := stubs.Triplet{X: rand.Intn(999999), Y: rand.Intn(999999), Z: rand.Intn(999999)}
		//Form a request to publish it in 'multiply'
		towork := stubs.PublishRequest{Topic: *topic, Triplet: newtriplet}
		//Call the broker
		err := client.Call(stubs.Publish, towork, status)
		if err != nil {
			fmt.Println("RPC client returned error:")
			fmt.Println(err)
			fmt.Println("Shutting down miner.")
			break
		}
		time.Sleep(1*time.Second)
	}
}
