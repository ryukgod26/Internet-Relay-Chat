package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	PRIVMSG = iota
	NONE
)

type MSG struct {
	username string
	chatType int
	message  string
}

type Client struct {
	address      string
	port         string
	username     string
	password     string
	nick         string
	chat         MSG
	conn         net.Conn
	reader       *textproto.Reader
	disconnected bool
	server       string
}

func (msg *MSG) IsPrivateMessage() bool {
	return msg.chatType == PRIVMSG
}

func send_data(conn net.Conn, msg string) {
	fmt.Fprintf(conn, "%s\r\n", msg)
}

func Init(address string, port string, password string, username string, nick string) Client {
	var client Client
	client.address = address
	client.port = port
	client.nick = nick
	client.password = password
	client.username = username
	return client
}

func ParseMessage(msg string) MSG {
	var message MSG
	if strings.Contains(msg, "PRIVMSG") {
		message.chatType = PRIVMSG
	} else {
		message.chatType = NONE
	}

	index := strings.Index(msg, ":") //try " :"
	if index != -1 {
		message.message = msg[index:1] //if " :" then use 2 instead of 1
	} else {
		message.message = msg
	}

	if message.chatType == PRIVMSG {
		index := strings.Index(msg, "!")
		if index != -1 {
			username := msg[:index]
			message.username = username[1:]
		}
	}

	return message
}

func (client *Client) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", client.address, client.port))
	handle_error(err)
	client.conn = conn
	client.Auth()
}

func (client *Client) Disconnect() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range c {
			client.disconnected = true
			client.conn.Close()
			os.Exit(0)
		}
	}()
}

func (client *Client) Auth() {

	if len(client.password) > 0 {
		send_data(client.conn, "PASS "+client.password)
	}

	if len(client.username) > 0 {
		send_data(client.conn, fmt.Sprintf("USER %s 0 * :%s",client.username,client.username))
	}

	if len(client.nick) > 0 {
		send_data(client.conn, "NICK "+client.nick+"\r\n")
	}
}

func (client *Client) Join(server string) {
	send_data(client.conn, "JOIN #"+server)
	client.server = server

	reader := bufio.NewReader(client.conn)
	tp := textproto.NewReader(reader)
	client.reader = tp
}

func (client *Client) HandlePong(data string) {
	if client.disconnected {
		return
	}
	if strings.HasPrefix(data, "PING") {
		send_data(client.conn, fmt.Sprintf("PONG %s\n", strings.TrimPrefix(data, "PING ")))
	}
}

func (client *Client) GetData() MSG {
	if client.disconnected {
		return MSG{}
	}
	data, err := client.reader.ReadLine()
	handle_error(err)
	client.HandlePong(data)
	fmt.Println(data)
	msg := ParseMessage(data)
	client.chat = msg
	return msg
}

func (client *Client) Say(msg string) {
	send_data(client.conn, fmt.Sprintf("PRIVMSG #%s : %s\n", client.server, msg))
}


func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
