package taskprogress

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
	"github.com/wagoodman/go-progress"

	"github.com/anchore/bubbly/bubbles/internal/testutil"
)

func subject(t testing.TB) (*progress.Manual, *progress.Stage, Model) {
	return subjectWaitGroup(t, &sync.WaitGroup{})
}

func subjectWaitGroup(_ testing.TB, wg *sync.WaitGroup) (*progress.Manual, *progress.Stage, Model) {
	prog := &progress.Manual{
		N:     40,
		Total: -1,
		Err:   nil,
	}
	stage := &progress.Stage{
		Current: "working",
	}

	tsk := New(
		wg,
		WithStagedProgressable(progress.StagedProgressable(&struct {
			progress.Stager
			progress.Progressable
		}{
			Stager:       stage,
			Progressable: prog,
		})),
		WithNoStyle(),
	)
	tsk.HideProgressOnSuccess = true
	tsk.TitleOptions = Title{
		Default: "Do work",
		Running: "Doing work",
		Success: "Did work",
		Failed:  "Failed at work :(",
	}
	tsk.Context = []string{
		"at home",
	}
	tsk.WindowSize = tea.WindowSizeMsg{
		Width:  100,
		Height: 60,
	}

	return prog, stage, tsk
}

func TestModel_View(t *testing.T) {

	tests := []struct {
		name       string
		taskGen    func(testing.TB) Model
		iterations int
	}{
		{
			name: "in progress without progress bar",
			taskGen: func(tb testing.TB) Model {
				prog, _, tsk := subject(t)
				prog.N, prog.Total = 40, -1
				return tsk
			},
		},
		{
			name: "in progress with progress bar",
			taskGen: func(tb testing.TB) Model {
				prog, _, tsk := subject(t)
				prog.N, prog.Total = 40, 100
				return tsk
			},
		},
		{
			name: "respond to title width",
			taskGen: func(tb testing.TB) Model {
				prog, stage, tsk := subject(t)
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				tsk.TitleWidth = 20
				return tsk
			},
		},
		{
			name: "hide stage on success",
			taskGen: func(tb testing.TB) Model {
				prog, stage, tsk := subject(t)
				tsk.HideStageOnSuccess = true
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk
			},
		},
		{
			name: "successfully finished hides progress bar",
			taskGen: func(tb testing.TB) Model {
				prog, stage, tsk := subject(t)
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk
			},
		},
		{
			name: "successfully finished keeps progress bar shown",
			taskGen: func(tb testing.TB) Model {
				prog, stage, tsk := subject(t)
				tsk.HideProgressOnSuccess = false
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk
			},
		},
		{
			name: "successfully can hide the entire line",
			taskGen: func(tb testing.TB) Model {
				prog, stage, tsk := subject(t)
				tsk.HideOnSuccess = true
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk
			},
		},
		{
			name: "no context",
			taskGen: func(tb testing.TB) Model {
				_, _, tsk := subject(t)
				tsk.Context = nil
				return tsk
			},
		},
		{
			name: "multiple hints",
			taskGen: func(tb testing.TB) Model {
				_, _, tsk := subject(t)
				tsk.Hints = []string{"info++", "info!++"}
				return tsk
			},
		},
		{
			name: "error",
			taskGen: func(tb testing.TB) Model {
				prog, _, tsk := subject(t)
				prog.SetCompleted()
				prog.Err = errors.New("woops")
				tsk.HideStageOnSuccess = false
				return tsk
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m tea.Model = tt.taskGen(t)
			tsk, ok := m.(Model)
			require.True(t, ok)
			got := testutil.RunModel(t, tsk, tt.iterations, TickMsg{
				Time:     time.Now(),
				Sequence: tsk.sequence,
				ID:       tsk.id,
			})
			t.Log(got)
			snaps.MatchSnapshot(t, got)
		})
	}
}

func Test_WaitGroupDone(t *testing.T) {
	waitGroupDone := func(_ Model, wg *sync.WaitGroup) {
		require.Equal(t, int32(0), waitCount(wg))
	}

	tests := []struct {
		name     string
		taskGen  func(testing.TB) (Model, *sync.WaitGroup)
		validate func(Model, *sync.WaitGroup)
	}{
		{
			name: "wg done when HideOnSuccess not set",
			taskGen: func(tb testing.TB) (Model, *sync.WaitGroup) {
				wg := &sync.WaitGroup{}
				prog, stage, tsk := subjectWaitGroup(t, wg)
				tsk.HideOnSuccess = false
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk, wg
			},
			validate: waitGroupDone,
		},
		{
			name: "wg done when HideOnSuccess set",
			taskGen: func(tb testing.TB) (Model, *sync.WaitGroup) {
				wg := &sync.WaitGroup{}
				prog, stage, tsk := subjectWaitGroup(t, wg)
				tsk.HideOnSuccess = true
				// note: we set progress to have a total size to ensure it is hidden
				prog.N, prog.Total = 100, 100
				stage.Current = "done!"
				return tsk, wg
			},
			validate: waitGroupDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, wg := tt.taskGen(t)
			_ = testutil.RunModel(t, model, 0, TickMsg{
				Time:     time.Now(),
				Sequence: model.sequence,
				ID:       model.id,
			})
			tt.validate(model, wg)
		})
	}
}

func waitCount(wg *sync.WaitGroup) int32 {
	v := reflect.ValueOf(wg).Elem()
	v = v.FieldByName("state1")
	state1 := v.Uint()

	// this is from waitgroup.go state() function:
	if unsafe.Alignof(state1) != 8 && uintptr(unsafe.Pointer(&state1))%8 != 0 {
		state := (*[3]uint32)(unsafe.Pointer(&state1))
		state1 = *(*uint64)(unsafe.Pointer(&state[1]))
	}

	return int32(state1 >> 32)
}
