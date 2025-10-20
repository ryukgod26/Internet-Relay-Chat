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
	index       int
	width       int
	height      int
	questions   []Question
	styles      *Styles
}

type Question struct {
	question string
	answer   string
	input 	 Input
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")

	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.ASCIIBorder()).
		Padding(1).
		Width(80)

	return s
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Enter Your Answer Here"
	answerField.Focus()
	answerField.CharLimit = 512
	answerField.Width = 60
	return &model{questions: questions, styles: styles}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(strings.TrimSpace(current.input.Value())) > 0 {
				m.Next()
				current.answer = strings.TrimSpace(current.input.Value())
				log.Printf("Question: %s,Answer: %s", current.question, current.answer)
				return m, current.input.Blur()
			} else {
				return m, nil
			}
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	current := m.questions[m.index]
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
			m.questions[m.index].question,
			m.styles.InputField.Render(current.input.View()),
		),
	)

}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}

}

func NewQuestion(question string) Question {
	return Question{question: question}
}

func NewShortQuestion(question string) Question{
	q := NewQuestion(question)
	field := NewShortAnswerField()
	q.input = field
	return q
}

func NewLongQuestion(question string) Question{
	q := NewQuestion(question)
	field := NewLongAnswerField()
	q.input = field
	return q
}

const (
	domain = "irc.oftc.net"
	port   = "6667"
	user   = "building101"
	nick   = "building101"
)

func main() {

	questions := []Question{
		NewShortQuestion("What is Your Name?"),
		NewShortQuestion("What is Your Username?"),
		NewShortQuestion("What is Your Nickname?"),
		NewShortQuestion("What is the Domain of the Server You Want to connect to?"),
		NewShortQuestion("What is the Serevr Port?"),
		NewShortQuestion("What is the Channnel name You wnat to enter?")}
		NewLongQuestion("Type the message you want to send?")

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
