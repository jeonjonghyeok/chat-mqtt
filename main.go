package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/AllenDang/gform"
	"github.com/jeonjonghyeok/chat-mqtt/mqtt"
)

const (
	broker = "broker.emqx.io"
	port   = 1883
)

func main() {
	gform.Init()
	topic := flag.String("room", "", "chatroom")
	name := flag.String("name", "", "username")
	flag.Parse()

	if *topic == "" {
		fmt.Println("Invalid setting for -room, must not be empty")
		return
	}
	if *name == "" {
		fmt.Println("Invalid setting for -name, must not be empty")
		return
	}
	fmt.Printf("Info:\n")
	fmt.Printf("\tname:  %s\n", *name)
	fmt.Printf("\troom:     %s\n", *topic)

	scanner := bufio.NewScanner(os.Stdin)

	client := mqtt.NewClient(broker, *name, port)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqtt.Sub(client, *topic)
	fmt.Println("msg를 입력하세요")
	var msg = ""
	for msg != "end" {
		fmt.Print(*name + ": ")
		scanner.Scan()
		msg = scanner.Text()
		mqtt.Pub(client, *topic, msg, *name)

	}

	client.Disconnect(250)
}
