package bubbly

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/wagoodman/go-partybus"
)

var _ tea.Model = (*dummyModel)(nil)

type dummyModel struct {
	id string
}

func (d dummyModel) Init() tea.Cmd {
	return nil
}

func (d dummyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(string); ok {
		d.id = msg
	}
	return d, nil
}

func (d dummyModel) View() string {
	return d.id
}

func dummyMsg(s any) tea.Cmd {
	return func() tea.Msg {
		return s
	}
}

func TestEventDispatcher_Handle(t *testing.T) {

	tests := []struct {
		name       string
		subject    *EventDispatcher
		event      partybus.Event
		wantModels []tea.Model
		wantCmd    tea.Cmd
	}{
		{
			name: "simple event",
			subject: func() *EventDispatcher {
				d := NewEventDispatcher()
				d.AddHandler("test", func(e partybus.Event) ([]tea.Model, tea.Cmd) {
					return []tea.Model{dummyModel{id: "model"}}, dummyMsg("updated")
				})
				return d
			}(),
			event: partybus.Event{
				Type: "test",
			},
			wantModels: []tea.Model{dummyModel{id: "model"}},
			wantCmd:    dummyMsg("updated"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModels, gotCmd := tt.subject.Handle(tt.event)
			if !reflect.DeepEqual(gotModels, tt.wantModels) {
				t.Errorf("Handle() got = %v (model), want %v", gotModels, tt.wantModels)
			}

			if gotCmd != nil && tt.wantCmd == nil {
				t.Fatal("got command, but want nil")
			} else if gotCmd == nil && tt.wantCmd != nil {
				t.Fatal("did not get command, but wanted one")
			}

			var (
				gotMsg  tea.Msg
				wantMsg tea.Msg
			)

			if gotCmd != nil {
				gotMsg = gotCmd()
			}

			if tt.wantCmd != nil {
				wantMsg = tt.wantCmd()
			}

			if !assert.Equal(t, wantMsg, gotMsg) {
				t.Errorf("Handle() got = %v (msg), want %v", gotMsg, wantMsg)
			}

		})
	}
}

func TestEventDispatcher_RespondsTo(t *testing.T) {

	d := NewEventDispatcher()
	d.AddHandler("test", func(e partybus.Event) ([]tea.Model, tea.Cmd) {
		return []tea.Model{dummyModel{id: "test-model"}}, dummyMsg("test-msg")
	})

	d.AddHandler("something", func(e partybus.Event) ([]tea.Model, tea.Cmd) {
		return []tea.Model{dummyModel{id: "something-model"}}, dummyMsg("something-msg")
	})

	tests := []struct {
		name    string
		subject *EventDispatcher
		want    []partybus.EventType
	}{
		{
			name:    "responds to registered event",
			subject: d,
			want:    []partybus.EventType{"test", "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.subject.RespondsTo()
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("RespondsTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerCollection_RespondsTo(t *testing.T) {
	d1 := NewEventDispatcher()
	d1.AddHandler("test", func(e partybus.Event) ([]tea.Model, tea.Cmd) {
		return []tea.Model{dummyModel{id: "test-model"}}, dummyMsg("test-msg")
	})

	d2 := NewEventDispatcher()
	d2.AddHandler("something", func(e partybus.Event) ([]tea.Model, tea.Cmd) {
		return []tea.Model{dummyModel{id: "something-model"}}, dummyMsg("something-msg")
	})

	subject := NewHandlerCollection(d1, d2)

	tests := []struct {
		name    string
		subject *HandlerCollection
		want    []partybus.EventType
	}{
		{
			name:    "responds to registered event from all handlers",
			subject: subject,
			want:    []partybus.EventType{"test", "something"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.subject.RespondsTo()
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("RespondsTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
