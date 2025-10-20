package main

import (
	"fmt"
	"irc/irc"
	"time"

	// "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	questions := []Question{
		NewShortQuestion("What is Your Name?"),
		NewShortQuestion("What is Your Username?"),
		NewShortQuestion("What is Your Nickname?"),
		NewShortQuestion("What is the Domain of the Server You Want to connect to?"),
		NewShortQuestion("What is the Serevr Port?"),
		NewShortQuestion("What is the Channnel name You wnat to enter?"),
		NewLongQuestion("Type the message you want to send?"),
	}

	ircIn := make(chan string, 64)
	ircOut := make(chan string, 64)
	channel := make(chan string, 64)
	m := New(questions, ircIn, ircOut, channel)

	f, err := tea.LogToFile("debug.log", "debug")
	irc.Handle_error(err)
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	go func() {
		_, err := p.Run()
		irc.Handle_error(err)
	}()

	// if _, err := p.Run(); err != nil {
	// 	irc.Handle_error(err)
	// }

	client := irc.Init(domain, port, "1223", user, nick)
	c := &client

	c.Connect()
	c.Disconnect()

	c.Join(<-m.channel)
	c.SayToNick(nick, "hello self test")
	res, err := c.GetResponse()
	fmt.Println("Response:", res)
	irc.Handle_error(err)

	go func() {
		for {
			test := c.GetData()
			fmt.Println(test)
		}
	}()

	go func() {
		for {
			line, err := c.GetResponse()
			irc.Handle_error(err)
			ircIn <- line
			p.Send(IrcMsg(line))
		}
	}()

	go func() {
		for out := range ircOut {
			c.SendRaw(out)
			p.Send(">>> " + out)

		}
	}()

	for {
		time.Sleep(time.Second)
	}
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println("Enter Your Message to send to irc server.")
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	if strings.TrimSpace(line) == "/quit" {
	// 		fmt.Println("Exiting.")
	// 		os.Exit(0)
	// 	}
	// 	if strings.TrimSpace(line) == "" {
	// 		continue
	// 	}
	// 	fmt.Println("Testing:", line)
	// 	c.Say(line)
	// 	res, err := c.GetResponse()
	// 	fmt.Println("Response:", res)
	// 	irc.Handle_error(err)
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
}
