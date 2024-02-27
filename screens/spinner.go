package screens

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/nomad/api"
	"github.com/smorenodp/nomadinspect/nomad"
)

type SpinnerScreen struct {
	actualNs    string
	matches     []regexp.Regexp
	spinner     spinner.Model
	client      nomad.NomadClient
	err         error
	finished    bool
	resources   []string
	nsChan      chan string
	jobChan     chan string
	endChan     chan string
	matchedChan chan Job
	Matched     []Job
	namespaces  []string
	and         bool
}

type Job struct {
	ID        string
	Namespace string
	Content   string
	Raw       *api.Job
	Matches   []regexp.Regexp
}

func (r *Job) Title() string       { return r.ID }
func (r *Job) Description() string { return r.Namespace }
func (r *Job) FilterValue() string { return r.ID }

type NamespaceMsg string
type JobMsg string
type MatchedMsg Job
type EndMsg string

type Namespace struct {
	Name string
	Jobs []string
}

type OutputMessage struct {
	Jobs []Job
}

func (s SpinnerScreen) retrieveInfo() {
	var err error
	namespaces := s.namespaces
	if len(s.namespaces) == 0 {
		namespaces, err = s.client.ListNs()
		if err != nil {
			return
		}
	}

	for _, ns := range namespaces {
		s.nsChan <- ns
		jobs, err := s.client.ListJobs(ns)
		if err != nil {
			return
		}
		for _, j := range jobs {
			s.jobChan <- j
			content, raw, err := s.client.InspectJob(ns, j)
			if err != nil {
				s.endChan <- "Error inspecting job"
			}
			matches := []regexp.Regexp{}
			for _, m := range s.matches {
				// if strings.Contains(content, m) {
				// 	matches = append(matches, m)
				// } else if s.and {
				// 	break
				// }
				if m.MatchString(content) {
					// f := m.FindAllString(content, -1)
					// s.file.WriteString(fmt.Sprintf("The matches are %+v with findings %s\n", s.matches, f))
					// matches = append(matches, f...)
					matches = append(matches, m)
				} else if s.and {
					break
				}
			}
			if (len(matches) > 0 && !s.and) || (s.and && len(matches) == len(s.matches)) {
				s.matchedChan <- Job{j, ns, content, raw, matches}
			}
		}
	}
	s.endChan <- ""
}

func (s *SpinnerScreen) CheckSelect() tea.Msg {
	select {
	case ns := <-s.nsChan:
		return NamespaceMsg(ns)
	case job := <-s.jobChan:
		return JobMsg(job)
	case job := <-s.matchedChan:
		return MatchedMsg(job)
	case <-s.endChan:
		return EndMsg("Check all info")
	}
}

func (s *SpinnerScreen) init() tea.Msg {
	go s.retrieveInfo()
	return s.CheckSelect()
}

func NewSpinnerScreen(namespaces, matches []string, and bool) SpinnerScreen {
	var r []regexp.Regexp
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	client, _ := nomad.New()
	for _, m := range matches {
		r = append(r, *regexp.MustCompile(m))
	}
	return SpinnerScreen{spinner: s, client: client, matches: r, namespaces: namespaces, and: and, nsChan: make(chan string), matchedChan: make(chan Job), endChan: make(chan string), jobChan: make(chan string), resources: make([]string, 5)}
}

func (s SpinnerScreen) Start() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.init)
}

func (s SpinnerScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		if s.finished {
			cmd = nil
		} else {
			s.spinner, cmd = s.spinner.Update(msg)
		}
		return s, cmd

	case NamespaceMsg:
		s.actualNs = string(msg)
		return s, s.CheckSelect
	case JobMsg:
		s.resources = append(s.resources[1:], string(msg))
		return s, s.CheckSelect
	case MatchedMsg:
		s.Matched = append(s.Matched, Job(msg))
		return s, s.CheckSelect
	case EndMsg:
		s.finished = true
		return s, func() tea.Msg { return OutputMessage{s.Matched} }
	default:
		return s, nil
	}
}

func (s SpinnerScreen) View() string {
	var str string
	if s.err != nil {
		str = s.err.Error()
	} else {
		str += fmt.Sprintf("%s\n\n", s.spinner.View())
		str += fmt.Sprintf("Checking namespace %s\n\n", s.actualNs)
		for _, resource := range s.resources {
			if resource != "" {
				str += fmt.Sprintf(" ... Checking %s\n", resource)
			}
		}
	}
	return str
}
