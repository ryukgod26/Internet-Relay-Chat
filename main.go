package main

import (
	"bufio"
	"fmt"
	"irc/irc"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type model struct {
	questions   []string
	width       int
	height      int
	index       int
	answerField textinput.Model
	styles      *Styles
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")

	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)

	return s
}

func New(questions []string) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Enter Your Answer Here"
	answerField.Focus()         // ensure placeholder / cursor appears
	answerField.CharLimit = 512 // allow longer input
	answerField.Width = 60      // ensure
	return &model{questions: questions, answerField: answerField, styles: styles}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.index < 5{
			m.index ++
			return m,nil
			}else{
			m.answerField.SetValue("Done!")
			return m,nil
			}
		}

	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading...."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
		lipgloss.Center,
		m.questions[m.index],
		m.styles.InputField.Render(m.answerField.View()),
	),
	)
	
	
}

const (
	domain = "irc.oftc.net"
	port   = "6667"
	user   = "building101"
	nick   = "building101"
)

func main() {

	questions := []string{
		"What is Your Name?",
		"What is Your Username?",
		"What is Your Nickname?",
		"What is the Domain of the Server You Want to connect to?",
		"What is the Serevr Port?",
		"What is the Channnel name You wnat to enter?"}

	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	irc.Handle_error(err)
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		irc.Handle_error(err)
	}

	client := irc.Init(domain, port, "1223", user, nick)
	c := &client

	c.Connect()
	c.Disconnect()

	c.Join("testchannel")
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
		fmt.Println("Testing:", line)
		c.Say(line)
		res, err := c.GetResponse()
		fmt.Println("Response:", res)
		irc.Handle_error(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
