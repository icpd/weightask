package priortask

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MockTask struct {
	priority int
	work     func() (any, error)
}

func (m MockTask) Priority() int {
	return m.priority
}

func (m MockTask) PerformTask(_ context.Context) (any, error) {
	return m.work()
}

func TestTasker_ProcessTasks(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []Task
		want      any
		expectErr bool
	}{
		{
			name: "Test with non-err",
			tasks: []Task{
				MockTask{priority: 1, work: func() (any, error) { return "Task1Result", nil }},
				MockTask{priority: 2, work: func() (any, error) { return "Task2Result", nil }},
			},
			want:      "Task2Result",
			expectErr: false,
		},
		{
			name: "Test with err",
			tasks: []Task{
				MockTask{priority: 1, work: func() (any, error) { return "Task1Result", nil }},
				MockTask{priority: 2, work: func() (any, error) { time.Sleep(time.Microsecond * 500); return "Task2Result", nil }},
				MockTask{priority: 3, work: func() (any, error) { return "Task3Result", fmt.Errorf("Task3Error") }},
				MockTask{priority: 4, work: func() (any, error) { time.Sleep(time.Second); return "Task4Result", nil }},
				MockTask{priority: 1, work: func() (any, error) { time.Sleep(time.Microsecond * 300); return "Task5Result", nil }},
			},
			want:      "Task4Result",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTaskController()
			for _, task := range tt.tasks {
				tc.AddTask(task)
			}
			got, err := tc.ProcessTasks(context.Background())
			if (err != nil) != tt.expectErr {
				t.Errorf("TaskController.ProcessTasks() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskController.ProcessTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTasker_ProcessTasksWithTimeoutCtx(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []Task
		want      any
		expectErr bool
	}{
		{
			name: "Test",
			tasks: []Task{
				MockTask{priority: 1, work: func() (any, error) { return "Task1Result", nil }},
				MockTask{priority: 2, work: func() (any, error) { time.Sleep(time.Second * 900); return "Task2Result", nil }},
				MockTask{priority: 3, work: func() (any, error) { return "Task3Result", fmt.Errorf("Task3Error") }},
				MockTask{priority: 4, work: func() (any, error) { time.Sleep(time.Second); return "Task4Result", nil }},
				MockTask{priority: 1, work: func() (any, error) { time.Sleep(time.Microsecond * 300); return "Task5Result", nil }},
			},
			want:      "Task1Result",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTaskController(WithPriorityList(new(PrioritySlice)))
			for _, task := range tt.tasks {
				tc.AddTask(task)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*500)
			defer cancel()
			got, err := tc.ProcessTasks(ctx)
			if (err != nil) != tt.expectErr {
				t.Errorf("TaskController.ProcessTasks() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskController.ProcessTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
