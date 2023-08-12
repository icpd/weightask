package priortask

import (
	"context"
	"errors"
	"sync"
)

var ErrNoResult = errors.New("no result")

type Task interface {
	// Priority returns the priority value of the task
	Priority() int
	// PerformTask Execute the task
	PerformTask() (any, error)
}

type PriorityList interface {
	Sort()
	Add(val int)
	Remove(val int)
	GetTopPriority() int
}

type TaskReport struct {
	priority int
	result   any
	err      error
}

type TaskController struct {
	tasks           []Task
	reportCh        chan *TaskReport
	effectiveReport *TaskReport
	priorityList    PriorityList
	err             error
}

func (t *TaskController) AddTask(task Task) {
	t.tasks = append(t.tasks, task)
	t.priorityList.Add(task.Priority())
}

func (t *TaskController) ProcessTasks(c context.Context) (any, error) {
	t.priorityList.Sort()

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	var wg sync.WaitGroup
	for _, task := range t.tasks {
		wg.Add(1)
		go func(ctx context.Context, tsk Task) {
			defer wg.Done()

			rst, err := tsk.PerformTask()

			select {
			case <-ctx.Done():
				return
			case t.reportCh <- &TaskReport{result: rst, err: err, priority: tsk.Priority()}:
			}
		}(ctx, task)
	}

	go func() {
		wg.Wait()
		close(t.reportCh)
	}()

	for rst := range t.reportCh {
		if rst == nil {
			continue
		}

		if rst.err != nil {
			t.priorityList.Remove(rst.priority)
			continue
		}

		if rst.priority == t.priorityList.GetTopPriority() {
			return rst.result, nil
		}

		if t.effectiveReport == nil {
			t.effectiveReport = rst
			continue
		}

		if rst.priority > t.effectiveReport.priority {
			t.priorityList.Remove(t.effectiveReport.priority)
			t.effectiveReport = rst
		} else {
			t.priorityList.Remove(rst.priority)
		}
	}

	if t.effectiveReport.result != nil {
		return t.effectiveReport.result, nil
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return nil, ErrNoResult
}

type Option func(*TaskController)

func WithPriorityList(pl PriorityList) Option {
	return func(t *TaskController) {
		t.priorityList = pl
	}
}

func NewTaskController(opts ...Option) *TaskController {
	t := &TaskController{
		reportCh: make(chan *TaskReport),
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.priorityList == nil {
		t.priorityList = new(PrioritySlice)
	}

	return t
}
