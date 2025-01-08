package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = ""

type model struct {
	textInput textinput.Model
	results   [][]string
	viewport  viewport.Model
	err       error
	width     int
	height    int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter search query..."
	ti.Focus()
	ti.Width = 40

	vp := viewport.New(100, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return model{
		textInput: ti,
		results:   [][]string{},
		viewport:  vp,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.textInput.Value() != "" {
				results, err := searchYoutube(m.textInput.Value())
				m.results = results
				m.err = err
				if err == nil {
					m.viewport.SetContent(m.renderTable())
				}
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) renderTable() string {
	t := table.NewWriter()
	t.SetStyle(table.Style{
		Box: table.StyleBoxRounded,
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})

	t.AppendHeader(table.Row{"#", "Title", "Channel", "Video ID"})

	for i, row := range m.results {
		t.AppendRow(table.Row{
			i + 1,
			lipgloss.NewStyle().MaxWidth(60).Render(row[0]),
			lipgloss.NewStyle().MaxWidth(30).Render(row[1]),
			row[2],
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 3, WidthMax: 3},
		{Number: 2, WidthMin: 60, WidthMax: 60},
		{Number: 3, WidthMin: 30, WidthMax: 30},
		{Number: 4, WidthMin: 15, WidthMax: 15},
	})

	return t.Render()
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true).
		Render("YouTube Search CLI") + "\n\n")

	s.WriteString(m.textInput.View() + "\n\n")

	if m.err != nil {
		return s.String() + lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	if len(m.results) > 0 {
		s.WriteString(m.viewport.View() + "\n")
	}

	s.WriteString("\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: Scroll • Enter: Search • Esc: Quit"))

	return s.String()
}

func searchYoutube(query string) ([][]string, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(25)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var results [][]string
	for _, item := range response.Items {
		results = append(results, []string{
			item.Snippet.Title,
			item.Snippet.ChannelTitle,
			item.Id.VideoId,
		})
	}

	return results, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
