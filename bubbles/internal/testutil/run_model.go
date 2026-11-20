package testutil

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
)

func RunModel(_ testing.TB, m tea.Model, iterations int, message tea.Msg) string {
	if iterations == 0 {
		iterations = 1
	}
	m.Init()
	var cmd tea.Cmd = func() tea.Msg {
		return message
	}

	for i := 0; cmd != nil && i < iterations; i++ {
		msgs := flatten(cmd())
		var nextCmds []tea.Cmd
		var next tea.Cmd
		for _, msg := range msgs {
			fmt.Printf("Message: %+v %+v\n", reflect.TypeOf(msg), msg)
			m, next = m.Update(msg)
			nextCmds = append(nextCmds, next)
		}
		cmd = tea.Batch(nextCmds...)
	}
	return render(m.View().Content)
}

// render downsamples the view the same way a bubbletea v2 program does when its output is not a terminal. Unlike
// lipgloss v1, lipgloss v2 always emits colors and leaves it to the program's renderer to strip them.
func render(s string) string {
	var buf bytes.Buffer
	w := colorprofile.Writer{Forward: &buf, Profile: colorprofile.NoTTY}
	_, _ = w.WriteString(s)
	return buf.String()
}

func flatten(p tea.Msg) (msgs []tea.Msg) {
	if p == nil {
		return nil
	}
	if batch, ok := p.(tea.BatchMsg); ok {
		for _, c := range batch {
			if c == nil {
				continue
			}
			msgs = append(msgs, flatten(c())...)
		}
		return msgs
	}
	return []tea.Msg{p}
}
