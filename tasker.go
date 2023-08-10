package priortask

import (
	"context"
	"errors"
	"math"
	"sync"
)

var ErrNoResult = errors.New("no result")

type Task interface {
	// Priority return the priority value of the task
	Priority() int
	// Execute process the task
	Execute() (any, error)
}

type taskRst struct {
	priority int
	result   any
	err      error
}

type Tasker struct {
	tasks        []Task
	rstCh        chan *taskRst
	effectiveRst *taskRst
	maxPriority  int
	taskCount    int
	err          error
}

func (t *Tasker) AddTask(task Task) {
	t.tasks = append(t.tasks, task)
	t.taskCount++
	if task.Priority() > t.maxPriority {
		t.maxPriority = task.Priority()
	}
}

func (t *Tasker) Process(c context.Context) (any, error) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	var wg sync.WaitGroup
	for _, task := range t.tasks {
		wg.Add(1)
		go func(ctx context.Context, tsk Task) {
			defer wg.Done()

			rst, err := tsk.Execute()

			select {
			case <-ctx.Done():
				return
			case t.rstCh <- &taskRst{result: rst, err: err, priority: tsk.Priority()}:
			}
		}(ctx, task)
	}

	go func() {
		wg.Wait()
		close(t.rstCh)
	}()

	for rst := range t.rstCh {
		t.taskCount--
		if rst.err != nil {
			continue
		}

		if rst.priority == t.maxPriority {
			return rst.result, nil
		}

		if t.effectiveRst == nil {
			t.effectiveRst = rst
			continue
		}

		if rst.priority > t.effectiveRst.priority {
			t.effectiveRst = rst
		}
	}

	if t.taskCount == 0 && t.effectiveRst.result != nil { // Check the count of tasks and check whether the result is nil or not
		return t.effectiveRst.result, nil
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return nil, ErrNoResult
}

func NewTasker() *Tasker {
	return &Tasker{
		rstCh:       make(chan *taskRst),
		maxPriority: math.MinInt64, // Initialize maxPriority to the smallest possible int value
	}
}
