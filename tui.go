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
	index             int
	width             int
	height            int
	questions         []Question
	originalQuestions []Question
	styles            *Styles
	done              bool

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
	orig := make([]Question, len(questions))
	copy(orig, questions)
	return &model{questions: questions, styles: styles, ircIn: ircIn, ircOut: ircOut, channel: channel,originalQuestions: orig}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]

	switch v := msg.(type) {
	case IrcMsg:
		m.messages = append(m.messages, string(v))
		return m, nil
	}
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			m.questions = make([]Question, len(m.originalQuestions))
            for i, oq := range m.originalQuestions {
                var q Question
                switch oq.input.(type) {
                case *Shortta:
                    q = NewShortQuestion(oq.question)
                default:
                    q = NewLongQuestion(oq.question)
                }
                q.answer = ""
                m.questions[i] = q
            }
            m.index = 0
            m.done = false
            return m, nil
		case "enter":
			text := strings.TrimSpace(current.input.Value())
			current.answer = text

			if text == "" {
				return m, nil
			}

			if m.ircOut != nil && m.channel != nil {
				chName := <-m.channel
				m.channel <- chName
				m.ircOut <- fmt.Sprintf("PRIVMSG #%s :%s", chName, text)
			}

			m.questions = append(m.questions[:m.index], m.questions[m.index+1:]...)
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			log.Printf("Question: %s,Answer: %s", current.question, current.answer)
			return m, nil
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading...."
	}

	if m.done {
		var out string
		out += "Message has been sent\n\n"
		for _, q := range m.questions {

			out += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		out += "\nPress ctrl+c to exit.\nPress ctrl+r to Send Another Message."
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			out,
		)
	}

	if len(m.questions) == 0 {
		m.done = true
		return m.View()
	}

	current := m.questions[m.index]
	start := 0
	var msgLines string
	if len(m.messages) > 10 {
		start = len(m.messages) - 10
	}
	// if m.done {
	// 	var output string
	// 	for _, q := range m.questions {
	// 		output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
	// 	}
	// 	return output
	// }

	for _, l := range m.messages[start:] {
		msgLines += l + "\n"
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
		m.done = true
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
