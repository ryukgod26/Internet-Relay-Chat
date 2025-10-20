package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type IrcMsg string

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type model struct {
	index     int
	width     int
	height    int
	questions []Question
	styles    *Styles
	done      bool

	ircIn    chan string
	ircOut   chan string
	messages []string
	channel  chan string
}

type Question struct {
	question string
	answer   string
	input    Input
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

func New(questions []Question, ircIn chan string, ircOut chan string, channel chan string) *model {
	styles := DefaultStyles()
	// answerField := textinput.New()
	// answerField.Placeholder = "Enter Your Answer Here"
	// answerField.Focus()
	// answerField.CharLimit = 512
	// answerField.Width = 60
	return &model{questions: questions, styles: styles, ircIn: ircIn, ircOut: ircOut, channel: channel}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]

	switch v := msg.(type){
	case IrcMsg:
		m.messages = append(m.messages, string(v))
		return m,nil
	}
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			text := strings.TrimSpace(current.input.Value())
			if text == ""{
				return m,nil
			}
			if m.ircOut != nil{
				m.ircOut <- fmt.Sprintf("PRIVMSG #%s :%s", m.channel, text)
			}

			m.Next()
			current.answer = text
			log.Printf("Question: %s,Answer: %s", current.question, current.answer)
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading...."
	}
	current := m.questions[m.index]
	start := 0
	var msgLines string
	if len(m.messages) > 10{
		start = len(m.messages) - 10
	}
	// if m.done {
	// 	var output string
	// 	for _, q := range m.questions {
	// 		output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
	// 	}
	// 	return output
	// }
	
	for _,l := range m.messages[start:]{
		msgLines += l+ "\n"
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

func NewShortQuestion(question string) Question {
	q := NewQuestion(question)
	field := NewShortAnswerField()
	q.input = field
	return q
}

func NewLongQuestion(question string) Question {
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
