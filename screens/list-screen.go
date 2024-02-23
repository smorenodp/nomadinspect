package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listKeyMap struct {
	viewJob key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		viewJob: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "View Job Info"),
		),
	}
}

func NewListScreen(matches []Job, width, height int) ListScreen {
	var (
		delegateKeys = &delegateKeyMap{}
		listKeys     = newListKeyMap()
	)
	delegate := newItemDelegate(delegateKeys)
	items := []list.Item{}
	for _, r := range matches {
		aux := r
		items = append(items, &aux)
	}
	matchesList := list.New(items, delegate, 0, 0)
	matchesList.Title = "Matches"
	matchesList.Styles.Title = titleStyle
	matchesList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.viewJob,
		}
	}
	h, v := appStyle.GetFrameSize()
	matchesList.SetSize(width-h, height-v)

	return ListScreen{List: matchesList, keys: listKeys, delegateKeys: delegateKeys}
}

type ListScreen struct {
	List         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

type ListCommandLaunched struct {
	Type string
}

func (l ListScreen) Start() tea.Cmd {
	return nil
}

func (l ListScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		l.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if l.List.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, l.keys.viewJob):
			return l, func() tea.Msg { return JobView(*l.List.SelectedItem().(*Job)) }
		}
	case JobView:
	}
	newListModel, cmd := l.List.Update(msg)
	l.List = newListModel
	return l, cmd
}

func (l ListScreen) View() string {
	return appStyle.Render(l.List.View())
}
