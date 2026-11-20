package prompt

import (
	"context"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anchore/bubbly/v2"
)

func typeText(m tea.Model, s string) tea.Model {
	for _, r := range s {
		m, _ = m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return m
}

func TestPrompt_EnterResponds(t *testing.T) {
	prompter := bubbly.NewPrompter("name?", false)
	var m tea.Model = New(prompter)
	m.Init()

	m = typeText(m, "bob")
	assert.Contains(t, m.View().Content, "name?")
	assert.Contains(t, m.View().Content, "bob")

	// a key release must not submit the prompt, only a press
	m, _ = m.Update(tea.KeyReleaseMsg{Code: tea.KeyEnter})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := prompter.Response(ctx)
	require.Error(t, err)

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got, err := prompter.Response(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "bob", got)

	// further input is ignored once the prompt is complete
	m = typeText(m, "x")
	assert.NotContains(t, m.View().Content, "bobx")
}

func TestPrompt_EmptyValueShowsValidationError(t *testing.T) {
	prompter := bubbly.NewPrompter("name?", false)
	var m tea.Model = New(prompter)
	m.Init()

	assert.Contains(t, m.View().Content, "value required")

	// note: this matches the bubbletea v1 behavior (bubbly v0.x), where enter on an invalid value still responds
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got, err := prompter.Response(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
	assert.Contains(t, m.View().Content, "value required")
}

func TestPrompt_SensitiveIsMasked(t *testing.T) {
	prompter := bubbly.NewPrompter("secret?", true)
	var m tea.Model = New(prompter)
	m.Init()

	m = typeText(m, "bob")
	t.Logf("before enter: %q", m.View().Content)
	assert.NotContains(t, m.View().Content, "bob")
	assert.Contains(t, m.View().Content, "●●●")

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	t.Logf("after enter: %q", m.View().Content)
	assert.NotContains(t, m.View().Content, "bob")
	assert.Contains(t, m.View().Content, "●●●")
}
