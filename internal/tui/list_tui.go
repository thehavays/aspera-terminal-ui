package tui

import (
	"fmt"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListType int

const (
	Received ListType = iota
	Shared
)

type listModel struct {
	list     list.Model
	listType ListType
	quitting bool
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			newType := Received
			if m.listType == Received {
				newType = Shared
			}
			// Fetch new items
			items, title := getListItems(newType)
			m.listType = newType
			m.list.Title = title
			m.list.SetItems(items)
			m.list.Select(0)
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

func getListItems(lt ListType) ([]list.Item, string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, "Error loading config"
	}

	var items []list.Item
	var title string
	filter := api.ListFilter{PageIndex: 1, PageSize: 100}

	token, err := api.EnsureValidToken()
	if err != nil {
		return nil, fmt.Sprintf("Auth Error: %v", err)
	}

	if lt == Received {
		resp, err := api.QueryReceivedRequests(cfg.Endpoint, token, filter)
		if err != nil {
			return nil, fmt.Sprintf("Error fetching received: %v", err)
		}
		title = "Received Transfers (Tab to switch)"
		items = make([]list.Item, len(resp.DataList))
		for i, req := range resp.DataList {
			dlStatus := "Not Downloaded"
			if req.Downloaded {
				dlStatus = "Downloaded"
			}
			items[i] = item{
				title: req.Title,
				desc:  fmt.Sprintf("ID: %s | From: %s | Time: %s | %s", req.RequestId, req.UserEmail, req.SendTime, dlStatus),
				id:    req.RequestId,
			}
		}
	} else {
		resp, err := api.QuerySharedRequests(cfg.Endpoint, token, filter)
		if err != nil {
			return nil, fmt.Sprintf("Error fetching shared: %v", err)
		}
		title = "Shared Transfers (Tab to switch)"
		items = make([]list.Item, len(resp.DataList))
		for i, req := range resp.DataList {
			items[i] = item{
				title: req.Subject,
				desc:  fmt.Sprintf("ID: %s | Time: %s | Status: %s | System: %s", req.RequestId, req.ShareTime, req.Status, req.System),
				id:    req.RequestId,
			}
		}
	}
	return items, title
}

// ShowList fetches the initial data and returns the model.
func ShowList(lt ListType) listModel {
	items, title := getListItems(lt)
	m := listModel{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		listType: lt,
	}
	m.list.Title = title
	return m
}

func RunListTUI(lt ListType) error {
	m := ShowList(lt)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
