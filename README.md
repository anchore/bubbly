# Bubbly

A small collection of [bubbletea](https://github.com/charmbracelet/bubbletea) TUI components and helpers.

## Versions

| bubbly module                    | bubbletea                    | install                                  |
|----------------------------------|------------------------------|------------------------------------------|
| `github.com/anchore/bubbly/v2`   | v2 (`charm.land/bubbletea/v2`) | `go get github.com/anchore/bubbly/v2@latest` |
| `github.com/anchore/bubbly` (v0.2.x) | v1 (`github.com/charmbracelet/bubbletea`) | `go get github.com/anchore/bubbly@v0.2`  |

The two modules have different import paths, so an application can stay on v0.2.x until it moves its own
`tea.Program` to bubbletea v2. Migrating means changing the import path and adapting to the bubbletea v2 API:
models return `tea.View` from `View()` (use `tea.NewView(s)`, and `.Content` to read a child model's rendering),
key presses arrive as `tea.KeyPressMsg`, and lipgloss colors are `color.Color` values (`lipgloss.Color("#777777")`).

## Status

***Consider this project to be in alpha. The API is not stable and may change at any time.***
