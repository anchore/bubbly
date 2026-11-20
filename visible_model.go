package bubbly

import tea "charm.land/bubbletea/v2"

type VisibleModel interface {
	IsVisible() bool
	tea.Model
}
