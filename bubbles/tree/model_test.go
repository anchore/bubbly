package tree

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"

	"github.com/anchore/bubbly"
	"github.com/anchore/bubbly/bubbles/internal/testutil"
)

var _ bubbly.VisibleModel = (*dummyViewer)(nil)

type dummyViewer struct {
	hidden bool
	state  string
}

func (d dummyViewer) IsVisible() bool {
	return !d.hidden
}

func (d dummyViewer) Init() tea.Cmd {
	return nil
}

func (d dummyViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d, nil
}

func (d dummyViewer) View() string {
	return d.state
}

func TestModel_View(t *testing.T) {

	tests := []struct {
		name       string
		taskGen    func(testing.TB) Model
		iterations int
	}{
		{
			name: "gocase",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()

				// └─ a
				//    └─ a-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"}))

				return subject
			},
		},
		{
			name: "sibling branches (one extra level)",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()

				// ├─ a
				// │  ├─ a-a
				// │  └─ a-b
				// └─ b
				//    └─ b-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b"}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))

				return subject
			},
		},
		{
			name: "sibling branches (lots of extra levels)",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()

				// ├─ a
				// │  ├─ a-a
				// │  ├─ a-b
				// │  │  ├─ a-b-a
				// │  │  ├─ a-b-b
				// │  │  └─ a-b-c
				// │  ├─ a-c
				// │  │  └─ a-c-a
				// │  └─ a-d
				// └─ b
				//    ├─ b-a
				//    │  ├─ b-a-a
				//    │  │  └─ b-a-a-a
				//    │  │     └─ b-a-a-a-a
				//    │  └─ b-a-b
				//    └─ b-b

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("a-b", "a-b-a", dummyViewer{state: "a-b-a"}))
				require.NoError(t, subject.Add("a-b", "a-b-b", dummyViewer{state: "a-b-b"}))
				require.NoError(t, subject.Add("a-b", "a-b-c", dummyViewer{state: "a-b-c"}))
				require.NoError(t, subject.Add("a", "a-c", dummyViewer{state: "a-c"}))
				require.NoError(t, subject.Add("a-c", "a-c-a", dummyViewer{state: "a-c-a"}))
				require.NoError(t, subject.Add("a", "a-d", dummyViewer{state: "a-d"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b"}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))
				require.NoError(t, subject.Add("b-a", "b-a-a", dummyViewer{state: "b-a-a"}))
				require.NoError(t, subject.Add("b-a-a", "b-a-a-a", dummyViewer{state: "b-a-a-a"}))
				require.NoError(t, subject.Add("b-a-a-a", "b-a-a-a-a", dummyViewer{state: "b-a-a-a-a"}))
				require.NoError(t, subject.Add("b-a", "b-a-b", dummyViewer{state: "b-a-b"}))
				require.NoError(t, subject.Add("b", "b-b", dummyViewer{state: "b-b"}))

				return subject
			},
		},
		{
			name: "multiline node",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()

				// ├─ a
				// │  more a...
				// │  ├─ a-a
				// │  │  a-a continued...
				// │  │  more a-a!
				// │  └─ a-b
				// └─ b
				//    more b...
				//    └─ b-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a\nmore a..."}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a\na-a continued...\nmore a-a!"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b\nmore b..."}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))

				return subject
			},
		},
		{
			name: "padded multiline node",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()
				subject.VerticalPadMultilineNodes = true

				// ├─ a
				// │  more a...
				// │  ├─ a-a
				// │  │  a-a continued...
				// │  │  more a-a!
				// │  └─ a-b
				// └─ b
				//    more b...
				//    └─ b-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a\nmore a..."}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a\na-a continued...\nmore a-a!"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("a", "a-c", dummyViewer{state: "a-c"}))
				require.NoError(t, subject.Add("a", "a-d", dummyViewer{state: "a-d"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b\nmore b..."}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))

				return subject
			},
		},
		{
			name: "hidden nodes",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()

				// └─ a
				//    └─ a-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"})) // shown as a leaf instead of a fork
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b", hidden: true}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b", hidden: true}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"})) // gets pruned entirely

				return subject
			},
		},
		{
			name: "margin",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()
				subject.Margin = "   "

				//    ├─ a
				//    │  ├─ a-a
				//    │  └─ a-b
				//    └─ b
				//       └─ b-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b"}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))

				return subject
			},
		},
		{
			name: "roots without prefix",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()
				subject.RootsWithoutPrefix = true

				// a
				// ├──a-a
				// └──a-b
				// b
				// └──b-a

				require.NoError(t, subject.Add("", "a", dummyViewer{state: "a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: "a-a"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: "a-b"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: "b"}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: "b-a"}))

				return subject
			},
		},
		{
			name: "horizontal padding",
			taskGen: func(tb testing.TB) Model {
				subject := NewModel()
				subject.Padding = "   "
				subject.RootsWithoutPrefix = true

				// ├─ a
				// │  ├─ a-a
				// │  ├─ a-b
				// │  │  ├─ a-b-a
				// │  │  ├─ a-b-b
				// │  │  └─ a-b-c
				// │  ├─ a-c
				// │  │  └─ a-c-a
				// │  └─ a-d
				// └─ b
				//    ├─ b-a
				//    │  ├─ b-a-a
				//    │  │  └─ b-a-a-a
				//    │  │     └─ b-a-a-a-a
				//    │  └─ b-a-b
				//    └─ b-b

				require.NoError(t, subject.Add("", "a", dummyViewer{state: " ✔ a"}))
				require.NoError(t, subject.Add("a", "a-a", dummyViewer{state: " ✔ a-a"}))
				require.NoError(t, subject.Add("a", "a-b", dummyViewer{state: " ✔ a-b"}))
				require.NoError(t, subject.Add("a-b", "a-b-a", dummyViewer{state: " ✔ a-b-a"}))
				require.NoError(t, subject.Add("a-b", "a-b-b", dummyViewer{state: " ✔ a-b-b"}))
				require.NoError(t, subject.Add("a-b", "a-b-c", dummyViewer{state: " ✔ a-b-c"}))
				require.NoError(t, subject.Add("a", "a-c", dummyViewer{state: " ⠼ a-c"}))
				require.NoError(t, subject.Add("a-c", "a-c-a", dummyViewer{state: " ⠼ a-c-a"}))
				require.NoError(t, subject.Add("a", "a-d", dummyViewer{state: " ⠼ a-d"}))
				require.NoError(t, subject.Add("", "b", dummyViewer{state: " ⠼ b"}))
				require.NoError(t, subject.Add("b", "b-a", dummyViewer{state: " ⠼ b-a"}))
				require.NoError(t, subject.Add("b-a", "b-a-a", dummyViewer{state: " ⠼ b-a-a"}))
				require.NoError(t, subject.Add("b-a-a", "b-a-a-a", dummyViewer{state: " ⠼ b-a-a-a"}))
				require.NoError(t, subject.Add("b-a-a-a", "b-a-a-a-a", dummyViewer{state: " ⠼ b-a-a-a-a"}))
				require.NoError(t, subject.Add("b-a", "b-a-b", dummyViewer{state: " ⠼ b-a-b"}))
				require.NoError(t, subject.Add("b", "b-b", dummyViewer{state: " ⠼ b-b"}))

				return subject
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m tea.Model = tt.taskGen(t)
			tsk, ok := m.(Model)
			require.True(t, ok)
			got := testutil.RunModel(t, tsk, tt.iterations, nil)
			t.Log(got)
			snaps.MatchSnapshot(t, got)
		})
	}
}
