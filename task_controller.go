package weightask

import (
	"context"
	"errors"
	"sync"
)

var ErrNoResult = errors.New("no result")

type Task interface {
	// Weight returns the weight value of the task
	Weight() int
	// PerformTask Execute the task
	PerformTask(ctx context.Context) (any, error)
}

type WeightList interface {
	Sort()
	Add(val int)
	Remove(val int)
	GetTopWeight() int
}

type TaskReport struct {
	weight int
	result any
	err    error
}

type TaskController struct {
	tasks           []Task
	reportCh        chan *TaskReport
	effectiveReport *TaskReport
	weightList      WeightList
}

func (t *TaskController) AddTask(task Task) {
	t.tasks = append(t.tasks, task)
}

func (t *TaskController) ProcessTasks(ctx context.Context) (any, error) {
	var wg sync.WaitGroup
	for _, task := range t.tasks {
		t.weightList.Add(task.Weight())

		wg.Add(1)
		go func(ctx context.Context, tsk Task) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case result := <-t.do(ctx, tsk):
				t.reportCh <- result
			}
		}(ctx, task)
	}

	go func() {
		wg.Wait()
		close(t.reportCh)
	}()

	t.weightList.Sort()
	for rst := range t.reportCh {
		if rst == nil {
			continue
		}

		if rst.err != nil {
			t.weightList.Remove(rst.weight)
			continue
		}

		if rst.weight == t.weightList.GetTopWeight() {
			return rst.result, nil
		}

		if t.effectiveReport == nil {
			t.effectiveReport = rst
			continue
		}

		if rst.weight > t.effectiveReport.weight {
			t.weightList.Remove(t.effectiveReport.weight)
			t.effectiveReport = rst
		} else {
			t.weightList.Remove(rst.weight)
		}
	}

	if t.effectiveReport != nil && t.effectiveReport.result != nil {
		return t.effectiveReport.result, nil
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return nil, ErrNoResult
}

func (t *TaskController) do(ctx context.Context, tsk Task) <-chan *TaskReport {
	trCh := make(chan *TaskReport, 1)
	go func() {
		rst, err := tsk.PerformTask(ctx)
		trCh <- &TaskReport{result: rst, err: err, weight: tsk.Weight()}
	}()

	return trCh
}

type Option func(*TaskController)

func WithWeightList(pl WeightList) Option {
	return func(t *TaskController) {
		t.weightList = pl
	}
}

func NewTaskController(opts ...Option) *TaskController {
	t := &TaskController{
		reportCh: make(chan *TaskReport),
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.weightList == nil {
		t.weightList = new(WeightSlice)
	}

	return t
}
