package tui

import (
	"fmt"
	"os"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 1)

type item struct {
	title, desc, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.id }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.id
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

// SelectRequest shows a list of received requests and lets the user select one.
func SelectRequest() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	fmt.Println("Fetching received requests...")
	token, err := api.EnsureValidToken()
	if err != nil {
		return "", err
	}
	resp, err := api.QueryReceivedRequests(cfg.Endpoint, token, api.ListFilter{PageIndex: 1, PageSize: 50})
	if err != nil {
		return "", err
	}

	items := make([]list.Item, len(resp.DataList))
	for i, req := range resp.DataList {
		items[i] = item{
			title: req.Title,
			desc:  fmt.Sprintf("ID: %s | From: %s | Time: %s", req.RequestId, req.UserEmail, req.SendTime),
			id:    req.RequestId,
		}
	}

	m := model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Select a request to download"

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	return finalModel.(model).choice, nil
}

// Multi-select model for files
type fileItem struct {
	name     string
	selected bool
}

type fileModel struct {
	requestID string
	files     []fileItem
	cursor    int
	offset    int // Scroll offset
	height    int // Window height
	selected  map[int]struct{}
	quitting  bool
	done      bool
}

func (m fileModel) Init() tea.Cmd {
	return nil
}

func (m fileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height - 6 // Reserve less space (header 2, footer 2, margins)
		if m.height < 5 {
			m.height = 5
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				if m.cursor >= m.offset+m.height {
					m.offset = m.cursor - m.height + 1
				}
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "a": // Select all
			for i := range m.files {
				m.selected[i] = struct{}{}
			}
		case "n": // Unselect all (None)
			m.selected = make(map[int]struct{})
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m fileModel) View() string {
	if m.quitting {
		return ""
	}

	header := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render(fmt.Sprintf("Files for request %s:", m.requestID)) + "\n"
	header += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Use [space] to select, [a] all, [n] none, [enter] to confirm") + "\n"

	s := ""
	end := m.offset + m.height
	if end > len(m.files) {
		end = len(m.files)
	}

	for i := m.offset; i < end; i++ {
		file := m.files[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if _, ok := m.selected[i]; ok {
			checked = "[x]"
		}

		line := fmt.Sprintf("%s %s %s", cursor, checked, file.name)
		if m.cursor == i {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render(line)
		}
		s += line + "\n"
	}

	footer := ""
	if len(m.files) > m.height {
		footer = fmt.Sprintf("--- [%d-%d / %d] ---", m.offset+1, end, len(m.files))
	}

	return docStyle.Render(header + s + footer)
}

// SelectFiles shows a list of files in a request and lets the user select which ones to download.
func SelectFiles(requestID string) ([]string, error) {
	fmt.Printf("Fetching files for request %s...\n", requestID)
	detail, err := api.GetRequestDetail(requestID)
	if err != nil {
		return nil, err
	}

	files := make([]fileItem, len(detail.Files))
	for i, f := range detail.Files {
		files[i] = fileItem{name: f.FileName}
	}

	m := fileModel{
		requestID: requestID,
		files:     files,
		selected:  make(map[int]struct{}),
		height:    20, // Default if WindowSizeMsg hasn't fired
	}
	// Default to all selected
	for i := range files {
		m.selected[i] = struct{}{}
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	fm := finalModel.(fileModel)
	if fm.quitting || !fm.done {
		os.Exit(0)
	}

	var selected []string
	for i := range fm.selected {
		selected = append(selected, fm.files[i].name)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no files selected")
	}

	return selected, nil
}
