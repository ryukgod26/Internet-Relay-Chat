package main

import (
	"fmt"
	"irc/irc"
	"log"
	"os"
	"strings"
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
		NewShortQuestion("What is your Password?"),
		NewLongQuestion("Type the message you want to send?"),
	}

	m := New(questions, nil, nil, make(chan string, 1))
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel,err := p.Run()
	irc.Handle_error(err)

	fm, ok := finalModel.(*model)
    if !ok {
        log.Fatal("unexpected model type")
    }

	f, err := tea.LogToFile("debug.log", "debug")
	irc.Handle_error(err)
	defer f.Close()


	//name := strings.TrimSpace(fm.questions[0].answer)
    username := strings.TrimSpace(fm.questions[1].answer)
    nick := strings.TrimSpace(fm.questions[2].answer)
    domain := strings.TrimSpace(fm.questions[3].answer)
    port := strings.TrimSpace(fm.questions[4].answer)
    channel := strings.TrimSpace(fm.questions[5].answer)
	password := strings.TrimSpace(fm.questions[6].answer)


	if domain == "" || port == "" || nick == "" || username == "" || channel == "" {
        fmt.Println("missing required IRC configuration, exiting")
        os.Exit(1)
    }

	// go func() {
	// 	_, err := p.Run()
	// 	irc.Handle_error(err)
	// }()

	// if _, err := p.Run(); err != nil {
	// 	irc.Handle_error(err)
	// }

	client := irc.Init(domain, port, password, username, nick)
	c := &client

	c.Connect()
	c.Disconnect()

	time.Sleep(500 * time.Millisecond)

	c.Join(channel)
	// c.SayToNick(nick, "hello self test")
	// res, err := c.GetResponse()
	// fmt.Println("Response:", res)
	// irc.Handle_error(err)

	    go func() {
        for {
            line, err := c.GetResponse()
            if err != nil {
                log.Println("irc read error:", err)
                return
            }
            if strings.TrimSpace(line) == "" {
                continue
            }
            fmt.Println(line)
        }
    }()


    c.Say("hello from TUI client")


    select {}
	// go func() {
	// 	for {
	// 		test := c.GetData()
	// 		fmt.Println(test)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		line, err := c.GetResponse()
	// 		irc.Handle_error(err)
	// 		ircIn <- line
	// 		p.Send(IrcMsg(line))
	// 	}
	// }()

	// go func() {
	// 	for out := range ircOut {
	// 		c.SendRaw(out)
	// 		p.Send(">>> " + out)

	// 	}
	// }()

	// for {
	// 	time.Sleep(time.Second)
	// }
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
