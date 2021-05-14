package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jeonjonghyeok/chat-mqtt/mqtt"
	"github.com/webview/webview"
)

const (
	broker = "broker.emqx.io"
	port   = 1883
)

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("mqtt chat application")
	w.SetSize(800, 600, webview.HintNone)

	//topic := flag.String("room", "", "chatroom")
	name := flag.String("name", "", "username")
	flag.Parse()
	/*
		w.Bind("noop", func() string {
			log.Println("hello")
			return "hello"
		})
		w.Bind("add", func(a, b int) int {
			return a + b
		})
		w.Bind("quit", func() {
			w.Terminate()
		})
		//w.Navigate("https://en.m.wikipedia.org/wiki/Main_page")
		w.Navigate(`data:text/html,
			<!doctype html>
			<html>
				<body>hello</body>
				<script>
					window.onload = function() {
						document.body.innerText = ` + "`hello, ${navigator.userAgent}`" + `;
						noop().then(function(res) {
							console.log('noop res', res);
							add(1, 2).then(function(res) {
								console.log('add res', res);
							});
						});
					};
				</script>
			</html>
		)`)
		w.Run()
	*/

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
	//fmt.Printf("\troom:     %s\n", *topic)

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
