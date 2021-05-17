package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jeonjonghyeok/chat-mqtt/mqtt"
)

const (
	broker = "broker.emqx.io"
	port   = 1883
)

func main() {
	topic := make([]string, 5)

	//flag
	name := flag.String("name", "", "username")
	flag.Parse()
	if *name == "" {
		fmt.Println("Invalid setting for -name, must not be empty")
		return
	}
	fmt.Printf("Info:\n")
	fmt.Printf("\tname:  %s\n", *name)

	//room
	fmt.Println("구독할 채팅방을 입력하세요:")
	room_scanner := bufio.NewScanner(os.Stdin)
	room_scanner.Scan()
	slice := strings.Split(room_scanner.Text(), " ")
	for i, str := range slice {
		topic[i] = str
	}
	fmt.Printf("\troom: ")
	for j := 0; j < len(topic); j++ {
		fmt.Printf(" %s", topic[j])
	}
	fmt.Println()

	//mqtt
	client := mqtt.NewClient(broker, *name, port)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for i := 0; i < len(topic); i++ {
		if topic[i] == "" {
			continue
		}
		mqtt.Sub(client, topic[i])

	}

	//message
	fmt.Println("msg를 입력하세요")
	var msg = ""
	for msg != "end" {
		fmt.Print(*name + ": ")
		message_scanner := bufio.NewScanner(os.Stdin)
		message_scanner.Scan()
		msg = message_scanner.Text()
		for i := 0; i < len(topic); i++ {
			if topic[i] == "" {
				continue
			}
			mqtt.Pub(client, topic[i], msg, *name)
		}
	}

	client.Disconnect(250)
}
