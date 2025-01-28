package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	TextInput   textinput.Model
	Results     [][]string
	Viewport    viewport.Model
	Err         error
	Width       int
	Height      int
	Loading     bool
	SelectedRow int
	ShowHelp    bool
}

type SearchResultMsg struct {
	Results [][]string
	Err     error
}
