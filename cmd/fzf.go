package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func fzfRun(ctx context.Context, state adapter.Adapter) (string, error) {
	servers, err := state.ListServers(ctx)
	if err != nil {
		return "", fmt.Errorf("list servers: %w", err)
	}
	sessions := make([]session, len(servers))
	for i, s := range servers {
		sessions[i] = session{
			label:   s.Label,
			workdir: s.Workdir,
		}
	}
	return fzf{}.Run(sessions)
}

type fzf struct {
}

const (
	SetInverse   = "\033[7m"
	ResetInverse = "\033[27m"
)

type session struct {
	label   string
	workdir string
}

func (s session) display(tab int) string {
	var filler strings.Builder
	for i := 0; i < tab-len(s.label); i++ {
		filler.WriteRune(' ')
	}
	return fmt.Sprintf("%s%s  %s", s.label, filler.String(), s.workdir)
}

func (f fzf) Run(sessions []session) (string, error) {
	choices := make([]string, len(sessions))
	dirs := make(map[string]session)
	maxima := 0
	for i, s := range sessions {
		choices[i] = s.label
		dirs[s.label] = s
		if len(s.label) > maxima {
			maxima = len(s.label)
		}
	}
	m := &model{
		choices: choices,
		dirs:    dirs,
		search:  "",
		matches: make([]string, len(sessions)),
		tab:     maxima,
	}
	_, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", fmt.Errorf("tea fzf: %w", err)
	}
	match := strings.TrimSpace(strings.Split(m.matches[0], "[")[0])
	return match, nil
}

type model struct {
	choices []string
	dirs    map[string]session
	matches []string
	search  string
	index   int
	tab     int
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			os.Exit(1)
		case tea.KeyEnter.String():
			if m.matches[0] != "" {
				return m, tea.Quit
			}
		case tea.KeyBackspace.String():
			if len(m.search) >= 1 {
				m.search = m.search[:len(m.search)-1]
			}
			m.index = 0
		case tea.KeyUp.String():
			m.index++
			if m.index > len(m.matches)-1 || m.matches[m.index] == "" {
				m.index = 0
			}
		case tea.KeyDown.String():
			m.index--
			if m.index < 0 {
				i := len(m.matches) - 1
				for m.matches[i] != "" {
					i--
				}
				m.index = i
			}
		default:
			msg := msg.String()
			if len(msg) == 1 && ('!' <= msg[0] || msg[0] <= '~') {
				m.search += msg
			}
			m.index = 0
		}
	}
	matches := fuzzy.RankFind(m.search, m.choices)
	sort.Sort(matches)
	for i, match := range matches {
		m.matches[i] = m.dirs[match.Target].display(m.tab)
	}
	for i := len(matches); i < len(m.matches); i++ {
		m.matches[i] = ""
	}
	return m, nil
}
func (m model) View() string {
	var s strings.Builder
	s.WriteString("Select server:\n")
	for j := range m.matches {
		i := len(m.matches) - 1 - j
		s.WriteString(SetInverse)
		if i == m.index {
			s.WriteRune('>')
		} else {
			s.WriteRune(' ')
			s.WriteString(ResetInverse)
		}
		s.WriteRune(' ')
		s.WriteString(m.matches[i])
		if i == m.index {
			s.WriteString(ResetInverse)
		}
		s.WriteRune('\n')
	}
	s.WriteString("Filter: [")
	s.WriteString(m.search)
	s.WriteRune(']')
	return s.String()
}
