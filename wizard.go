package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type Shortta struct {
	textinput textinput.Model
}

type Longta struct {
	textarea textarea.Model
}

func (sa *Shortta) View() string {
	return sa.textinput.View()
}

func (sa *Shortta) Update(msg tea.Msg) (Input,tea.Cmd){
	var cmd tea.Cmd
	sa.textinput,cmd = sa.textinput.Update(msg)
	return sa,cmd
}

func (sa *Shortta) Value() string {
	return sa.textinput.Value()
}

func (sa *Shortta) Blur() tea.Msg {
	return	sa.textinput.Blur
}

func NewShortAnswerField() *Shortta {
	ti := textinput.New()
	ti.Placeholder = "Enter Your Answer Here"
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 60
	return &Shortta{ti}
}

func (la *Longta) View() string {
	return la.textarea.View()
}

func (la *Longta) Update(msg tea.Msg) (Input,tea.Cmd) {
	var cmd tea.Cmd
	la.textarea,cmd = la.textarea.Update(msg)
	return la,cmd
}

func (la *Longta) Value() string {
	return la.textarea.Value()
}

func (la *Longta) Blur() tea.Msg {
	return	la.textarea.Blur
}

func NewLongAnswerField() *Longta {
	ta := textarea.New()
	ta.Placeholder = "Enter Your Answer Here"
	ta.Focus()
	ta.CharLimit = 512

	return &Longta{ta}
}

