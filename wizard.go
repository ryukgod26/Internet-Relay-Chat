package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur()
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

type ShortAnswerField struct {
	textinput textinput.Model
}

type LongAnswerField struct {
	textarea textarea.Model
}

func (sa *ShortAnswerField) Update(msg tea.Msg) {

}

func (sa *ShortAnswerField) Value() string {
	return sa.textinput.Value()
}

func (sa *ShortAnswerField) Blur() {
	sa.textinput.Blur()
}

func NewShortAnswerField() *ShortAnswerField {
	ti := textinput.New()
	return &ShortAnswerField{ti}
}

func (la *LongAnswerField) Value() string {
	return la.textarea.Value()
}

func (la *LongAnswerField) Blur() {
	la.textarea.Blur()
}

func NewLongAnswerField() *LongAnswerField {
	ta := textarea.New()
	return &LongAnswerField{ta}
}
