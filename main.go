package main

import (
	"bufio"
	"fmt"
	"irc/irc"
	"log"
	"os"
	"strings"
)

const (
	domain = "irc.oftc.net"
	port   = "6667"
	user   = "building101"
	nick   = "building101"
)

func main() {
	client := irc.Init(domain, port, "1223", user, nick)
	c := &client

	c.Connect()
	c.Disconnect()

	c.Join("testchannel")
	c.SayToNick(nick, "hello self test")
	res,err := c.GetResponse()
	fmt.Println("Response:",res)
	irc.Handle_error(err)

	go func() {
		for {
			test := c.GetData()
			fmt.Println(test)
		}
	}()
	
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter Your Message to send to irc server.")
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "/quit" {
			fmt.Println("Exiting.")
			os.Exit(0)
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Println("Testing:",line)
		c.Say(line)
		res,err := c.GetResponse()
		fmt.Println("Response:",res)
		irc.Handle_error(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}