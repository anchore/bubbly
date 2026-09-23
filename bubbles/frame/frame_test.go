package frame

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockModel struct {
	updateCalled bool
	view         string
}

func (m mockModel) Init() tea.Cmd {
	return nil
}

func (m mockModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	m.updateCalled = true
	return m, nil
}

func (m mockModel) View() string {
	return m.view
}

type mockTerminalElement struct {
	mockModel
	isAlive bool
}

func (m mockTerminalElement) IsAlive() bool {
	return m.isAlive
}

type mockVisibleElement struct {
	mockModel
	isHidden bool
}

func (m mockVisibleElement) IsHidden() bool {
	return m.isHidden
}

type mockImprintableElement struct {
	mockModel
	shouldImprint bool
}

func (m mockImprintableElement) ShouldImprint() bool {
	return m.shouldImprint
}

func TestFrame_Update_PruneTerminalElement(t *testing.T) {
	frame := New()
	model := &mockTerminalElement{isAlive: false}

	frame.AppendModel(model)

	m, _ := frame.Update(nil)
	actual := m.(*Frame)

	assert.Empty(t, actual.models)
}

func TestFrame_Update_HideVisibleElement(t *testing.T) {
	frame := New()
	model := mockVisibleElement{isHidden: true}

	frame.AppendModel(model)

	m, _ := frame.Update(nil)
	actual := m.(*Frame)

	require.NotEmpty(t, actual.models)
	assert.True(t, actual.models[0].hidden)
}

func TestFrame_Update_ImprintImprintableElement(t *testing.T) {
	frame := New()
	model := &mockImprintableElement{shouldImprint: true, mockModel: mockModel{view: "imprinted!"}}

	frame.AppendModel(model)

	m, cmds := frame.Update(nil)
	actual := m.(*Frame)

	assert.True(t, actual.models[0].expired)
	assert.NotNil(t, cmds)
}

func TestFrame_Update_UpdateElement(t *testing.T) {
	frame := New()
	model := &mockModel{}

	frame.AppendModel(model)

	m, cmds := frame.Update(nil)
	actual := m.(*Frame)

	assert.Nil(t, cmds)
	assert.True(t, actual.models[0].model.(mockModel).updateCalled)
}

// Mock model for VisibleElement and TerminalElement
type mockVisibleTerminalModel struct {
	view          string
	isHidden      bool
	isAlive       bool
	shouldImprint bool
}

func (m mockVisibleTerminalModel) Init() tea.Cmd {
	return nil
}

func (m mockVisibleTerminalModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mockVisibleTerminalModel) View() string {
	return m.view
}

func (m mockVisibleTerminalModel) IsHidden() bool {
	return m.isHidden
}

func (m mockVisibleTerminalModel) IsAlive() bool {
	return m.isAlive
}

func (m mockVisibleTerminalModel) ShouldImprint() bool {
	return m.shouldImprint
}

func TestFrame_View_HiddenElement(t *testing.T) {
	frame := New()
	model := mockVisibleTerminalModel{view: "visible", isHidden: true, isAlive: true}

	frame.AppendModel(model)

	m, _ := frame.Update(nil)

	assert.Empty(t, m.View())
}

func TestFrame_View_VisibleElement(t *testing.T) {
	frame := New()
	model := mockVisibleTerminalModel{view: "visible", isHidden: false, isAlive: true}

	frame.AppendModel(model)

	m, _ := frame.Update(nil)

	assert.Contains(t, m.View(), "visible")
}

func TestFrame_View_DeadElement(t *testing.T) {
	frame := New()
	model := mockVisibleTerminalModel{view: "should not be seen", isHidden: false, isAlive: false}

	frame.AppendModel(model)

	m, _ := frame.Update(nil)

	assert.NotContains(t, m.View(), "should not be seen")
}

func TestFrame_View_WithFooter(t *testing.T) {
	frame := New()
	frame.ShowFooter(true)
	frame.TruncateFooter(false)
	model := mockVisibleTerminalModel{view: "visible", isHidden: false, isAlive: true}

	frame.AppendModel(model)
	frame.Footer().Write([]byte("log line 1\nlog line 2"))

	m, _ := frame.Update(nil)

	viewOutput := m.View()

	assert.Equal(t, "visible\nlog line 1\nlog line 2\n", viewOutput)
}

func TestFrame_View_WithTruncatedFooter(t *testing.T) {
	frame := New()
	frame.ShowFooter(true)
	frame.TruncateFooter(true)
	frame.windowSize = tea.WindowSizeMsg{Height: 3} // room for the model, one log line, and the trailing blank row
	model := mockVisibleTerminalModel{view: "visible", isHidden: false, isAlive: true}

	frame.AppendModel(model)
	frame.Footer().Write([]byte("log line 1\nlog line 2\nlog line 3"))

	m, _ := frame.Update(nil)

	viewOutput := m.View()

	assert.Equal(t, "visible\nlog line 3\n", viewOutput)
}

func TestFrame_View_FooterLayout(t *testing.T) {
	tests := []struct {
		name     string
		views    []string
		footer   string
		truncate bool
		height   int
		want     string
	}{
		{
			name:     "every row is newline terminated",
			views:    []string{"visible"},
			truncate: true,
			height:   5,
			want:     "visible\n",
		},
		{
			name:     "no models means no leading newline",
			footer:   "log line 1\nlog line 2",
			truncate: true,
			height:   5,
			want:     "log line 1\nlog line 2\n",
		},
		{
			name:     "trailing newline in footer does not cost a row",
			views:    []string{"visible"},
			footer:   "log line 1\nlog line 2\nlog line 3\n",
			truncate: true,
			height:   4,
			want:     "visible\nlog line 2\nlog line 3\n",
		},
		{
			name:     "blank footer lines are dropped",
			views:    []string{"visible"},
			footer:   "\nlog line 1\n\nlog line 2\n",
			truncate: false,
			want:     "visible\nlog line 1\nlog line 2\n",
		},
		{
			name:     "multi-line and multiple models count every row",
			views:    []string{"a\nb", "c"},
			footer:   "log line 1\nlog line 2",
			truncate: true,
			height:   5,
			want:     "a\nb\nc\nlog line 2\n",
		},
		{
			name:     "models taller than the window hide the footer",
			views:    []string{"a\nb\nc"},
			footer:   "log line 1\nlog line 2",
			truncate: true,
			height:   2,
			want:     "a\nb\nc\n",
		},
		{
			name:     "no window size yet hides the footer",
			views:    []string{"visible"},
			footer:   "log line 1\nlog line 2",
			truncate: true,
			want:     "visible\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := New()
			frame.TruncateFooter(tt.truncate)
			frame.windowSize = tea.WindowSizeMsg{Height: tt.height}
			for _, v := range tt.views {
				frame.AppendModel(mockVisibleTerminalModel{view: v, isHidden: false, isAlive: true})
			}
			frame.Footer().Write([]byte(tt.footer))

			m, _ := frame.Update(nil)

			assert.Equal(t, tt.want, m.View())
		})
	}
}

func TestFrame_View_NoFooter(t *testing.T) {
	frame := New()
	frame.ShowFooter(false)
	model := mockVisibleTerminalModel{view: "visible", isHidden: false, isAlive: true}

	frame.AppendModel(model)
	frame.Footer().Write([]byte("log line 1\nlog line 2"))

	m, _ := frame.Update(nil)

	viewOutput := m.View()

	assert.Contains(t, viewOutput, "visible")
	assert.NotContains(t, viewOutput, "log line 1")
	assert.NotContains(t, viewOutput, "log line 2")
}
