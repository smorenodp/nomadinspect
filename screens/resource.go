package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#60DEA9", Dark: "#60DEA9"}).BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return headerStyle.Copy().BorderStyle(b)
	}()
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#60DEA9", Dark: "#60DEA9"}).PaddingBottom(1)
)

type EndView string

type ResourceScreen struct {
	Content       string
	Name          string
	Namespace     string
	Matches       []string
	MatchedLines  []int
	matchedIndex  int
	matchedLoaded bool
	viewPort      viewport.Model
}

func NewResourceScreen(j Job, w, h int) ResourceScreen {
	r := ResourceScreen{
		Content:   j.Content,
		Name:      j.ID,
		Namespace: j.Namespace,
		Matches:   j.Matches,
	}

	headerHeight := lipgloss.Height(r.headerView())
	footerHeight := lipgloss.Height(r.footerView())
	helpHeight := lipgloss.Height(r.helpView())
	verticalMarginHeight := headerHeight + footerHeight + helpHeight
	r.viewPort = viewport.New(w, h-verticalMarginHeight)

	r.Content, _ = PrettyString(r.Content)
	r.highlightMatched()
	r.viewPort.SetContent(r.Content)
	r.viewPort.YPosition = 0
	return r
}

func (r *ResourceScreen) Start() tea.Cmd {
	return nil
}

func (r *ResourceScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.viewPort.Width = msg.Width
		r.viewPort.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return r, func() tea.Msg { return EndView("Viewing resource ended") }
		case "m":
			if !r.matchedLoaded {
				r.getMatchedLines()
			}
			r.move2Matched()
		}
	}
	r.viewPort, cmd = r.viewPort.Update(msg)

	return r, cmd
}

func (r *ResourceScreen) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", r.headerView(), r.viewPort.View(), r.footerView(), r.helpView())
}

func (r *ResourceScreen) headerView() string {
	title := headerStyle.Render(fmt.Sprintf("%s > %s", r.Namespace, r.Name))
	line := strings.Repeat("─", max(0, r.viewPort.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (r *ResourceScreen) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", r.viewPort.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, r.viewPort.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (r *ResourceScreen) helpView() string {
	help := helpStyle.Render("Press enter to return to the list, m to move to the next match")
	return lipgloss.JoinHorizontal(lipgloss.Center, help)
}

func (r *ResourceScreen) highlightMatched() {
	for _, m := range r.Matches {
		r.Content = strings.ReplaceAll(r.Content, m, fmt.Sprintf("\x1B[31m%s\x1B[0m", m))
	}
}

func (r *ResourceScreen) getMatchedLines() {
	lines := strings.Split(r.Content, "\n")
	for index, line := range lines {
		for _, m := range r.Matches {
			if strings.Contains(line, m) {
				r.MatchedLines = append(r.MatchedLines, index)
			}
		}
	}
}

func (r *ResourceScreen) move2Matched() {
	index := r.MatchedLines[r.matchedIndex] - (r.viewPort.VisibleLineCount() / 2)
	r.viewPort.SetYOffset(index)
	r.setBold()
	r.matchedIndex++
	if r.matchedIndex == len(r.MatchedLines) {
		r.matchedIndex = 0
	}
}

func (r *ResourceScreen) setBold() {
	index := r.MatchedLines[r.matchedIndex]
	lines := strings.Split(r.Content, "\n")
	line := lines[index]
	for _, m := range r.Matches {
		line = strings.ReplaceAll(line, fmt.Sprintf("\x1B[31m%s\x1B[0m", m), fmt.Sprintf("\x1B[31;1;4m%s\x1B[0m", m))
	}
	lines[index] = line
	r.viewPort.SetContent(strings.Join(lines, "\n"))
}
