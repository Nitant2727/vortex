package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"

	"vortex/pkg/models"
	"vortex/pkg/ui"
	"vortex/pkg/utils"
	"vortex/pkg/youtube"
)

type model struct {
	models.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter search query..."
	ti.Focus()
	ti.Width = 40

	vp := viewport.New(100, 20)
	vp.Style = ui.ViewportStyle

	return model{
		Model: models.Model{
			TextInput:   ti,
			Results:     [][]string{},
			Viewport:    vp,
			Loading:     false,
			SelectedRow: 0,
			ShowHelp:    false,
		},
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
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width - 4
		m.Viewport.Height = msg.Height - 8

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.TextInput.Value() != "" {
				m.Loading = true
				m.SelectedRow = 0
				return m, searchYoutubeCmd(m.TextInput.Value())
			}
		case "up", "k":
			if len(m.Results) > 0 {
				m.SelectedRow = utils.Max(0, m.SelectedRow-1)
				m.Viewport.SetContent(renderTable(m.Model))
			}
		case "down", "j":
			if len(m.Results) > 0 {
				m.SelectedRow = utils.Min(len(m.Results)-1, m.SelectedRow+1)
				m.Viewport.SetContent(renderTable(m.Model))
			}
		case "o":
			if len(m.Results) > 0 && m.SelectedRow < len(m.Results) {
				videoID := m.Results[m.SelectedRow][2]
				url := youtube.GetVideoURL(videoID)
				go utils.OpenBrowser(url)
			}
		case "?":
			m.ShowHelp = !m.ShowHelp
		}
	case models.SearchResultMsg:
		m.Loading = false
		m.Results = msg.Results
		m.Err = msg.Err
		if msg.Err == nil {
			m.Viewport.SetContent(renderTable(m.Model))
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func renderTable(m models.Model) string {
	t := table.NewWriter()
	t.SetStyle(table.Style{
		Box: table.StyleBoxRounded,
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})

	t.AppendHeader(table.Row{"#", "Title", "Channel", "Video ID", "URL"})

	for i, row := range m.Results {
		videoURL := youtube.GetVideoURL(row[2])
		style := lipgloss.NewStyle()
		if i == m.SelectedRow {
			style = style.Background(lipgloss.Color("62")).Foreground(lipgloss.Color("0"))
		}

		t.AppendRow(table.Row{
			style.Render(fmt.Sprint(i + 1)),
			style.Copy().MaxWidth(50).Render(row[0]),
			style.Copy().MaxWidth(25).Render(row[1]),
			style.Render(row[2]),
			style.Copy().Foreground(lipgloss.Color("45")).Render(videoURL),
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 3, WidthMax: 3},
		{Number: 2, WidthMin: 50, WidthMax: 50},
		{Number: 3, WidthMin: 25, WidthMax: 25},
		{Number: 4, WidthMin: 12, WidthMax: 12},
		{Number: 5, WidthMin: 45, WidthMax: 45},
	})

	return t.Render()
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(ui.TitleStyle.Render("YouTube Search CLI"))
	s.WriteString("\n\n")
	s.WriteString(m.TextInput.View())
	s.WriteString("\n\n")

	if m.Loading {
		return s.String() + ui.LoadingStyle.Render("Searching...")
	}

	if m.Err != nil {
		return s.String() + ui.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.Err))
	}

	if len(m.Results) > 0 {
		s.WriteString(m.Viewport.View())
		s.WriteString("\n")
	}

	s.WriteString("\n")
	if m.ShowHelp {
		s.WriteString(ui.HelpStyle.Render("Controls:\n"))
		s.WriteString(ui.HelpStyle.Render("↑/k, ↓/j: Navigate • Enter: Search • o: Open in browser • ?: Toggle help • Esc: Quit"))
	} else {
		s.WriteString(ui.HelpStyle.Render("↑/↓: Navigate • Enter: Search • o: Open • ?: Help • Esc: Quit"))
	}

	return s.String()
}

func searchYoutubeCmd(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := youtube.Search(query)
		return models.SearchResultMsg{Results: results, Err: err}
	}
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
