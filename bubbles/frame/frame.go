package frame

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type VisibleElement interface {
	IsHidden() bool
}

type LimitedElement interface {
	IsAlive() bool
}

type ExpiredElement interface {
	IsExpired() bool
}

//type Printer interface {
//	Printf(template string, args ...interface{})
//	Println(args ...interface{})
//}

type Frame struct {
	footer         *bytes.Buffer
	models         []annotatedModel
	windowSize     tea.WindowSizeMsg
	showFooter     bool
	truncateFooter bool

	//printer Printer
}

type annotatedModel struct {
	model   tea.Model
	expired bool
	hidden  bool
}

func New() *Frame {
	return &Frame{
		footer:         &bytes.Buffer{},
		showFooter:     true,
		truncateFooter: true,
	}
}

func (f Frame) Footer() io.ReadWriter {
	return f.footer
}

func (f *Frame) ShowFooter(set bool) {
	f.showFooter = set
}

func (f *Frame) TruncateFooter(set bool) {
	f.truncateFooter = set
}

func (f *Frame) AppendModel(uiElement tea.Model) {
	f.models = append(f.models, annotatedModel{model: uiElement})
}

//func (f *Frame) Printer(p Printer) {
//	f.printer = p
//}

func (f Frame) Init() tea.Cmd {
	return nil
}

func (f Frame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		f.windowSize = msg
	}

	var cmds []tea.Cmd

	//// prune any models that are no longer alive
	//var alive []tea.Model
	//for _, p := range f.models {
	//	if le, ok := p.(LimitedElement); ok && le.IsAlive() {
	//		alive = append(alive, p)
	//	}
	//}
	//
	//// prune any models that should be persisted
	//var unexpired []tea.Model
	//for _, p := range alive {
	//	if le, ok := p.(ExpiredElement); ok {
	//		if le.IsExpired() {
	//			cmd := tea.Printf("%s", p.View())
	//			cmds = append(cmds, cmd)
	//		} else {
	//			unexpired = append(unexpired, p)
	//		}
	//	}
	//}
	//
	//// update remaining models
	//for i, el := range unexpired {
	//	newEl, cmd := el.Update(msg)
	//	cmds = append(cmds, cmd)
	//	unexpired[i] = newEl
	//}
	//
	//f.models = unexpired

	// prune any models that are no longer alive
	for i := 0; i < len(f.models); i++ {
		if p, ok := f.models[i].model.(LimitedElement); ok && !p.IsAlive() {
			f.models = append(f.models[:i], f.models[i+1:]...)
			i--
			continue
		}

		if f.models[i].expired {
			f.models = append(f.models[:i], f.models[i+1:]...)
			i--
			continue
		}

		if p, ok := f.models[i].model.(VisibleElement); ok && p.IsHidden() {
			f.models[i].hidden = true
		} else {
			f.models[i].hidden = false
		}

		if p, ok := f.models[i].model.(ExpiredElement); ok && p.IsExpired() {
			f.models[i].expired = true
			//if f.printer != nil {
			//	// since this is a blocking send, we need to do this in a goroutine
			//	go f.printer.Printf("%s", f.models[i].model.View())
			//} else {
			// this races with other messages (like exit) since quit (a message) will be faster to process
			// than a command (run in another go routine).
			cmd := tea.Printf("%s", f.models[i].model.View())
			cmds = append(cmds, cmd)
			//}
		}
	}

	for i, el := range f.models {
		if el.expired {
			continue
		}
		newEl, cmd := el.model.Update(msg)
		cmds = append(cmds, cmd)
		f.models[i].model = newEl
	}
	return f, tea.Batch(cmds...)
}

func (f Frame) View() string {
	// all UI elements
	var strs []string
	for _, p := range f.models {
		if p.hidden {
			continue
		}
		rendered := p.model.View()
		if len(rendered) > 0 {
			strs = append(strs, rendered)
		}
	}

	str := strings.Join(strs, "\n")

	// log events
	if f.showFooter {
		contents := f.footer.String()
		if f.truncateFooter {
			logLines := strings.Split(contents, "\n")
			logMax := f.windowSize.Height - strings.Count(str, "\n")
			trimLog := len(logLines) - logMax
			if trimLog > 0 && len(logLines) >= trimLog {
				logLines = logLines[trimLog:]
			}
			for _, line := range logLines {
				if len(line) > 0 {
					str += fmt.Sprintf("%s\n", line)
				}
			}
		} else {
			str += contents
		}
	}
	return str
}
