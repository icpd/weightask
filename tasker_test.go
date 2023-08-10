package priortask

import (
	"context"
	"fmt"
	"runtime"
	"testing"
)

type MockTask struct {
	priority int
	work     func() (any, error)
}

func (m MockTask) Priority() int {
	return m.priority
}

func (m MockTask) Execute() (any, error) {
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
				MockTask{priority: 2, work: func() (any, error) { return "Task2Result", nil }},
				MockTask{priority: 3, work: func() (any, error) { return "Task3Result", fmt.Errorf("Task2Error") }},
			},
			want:      "Task2Result",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasker := NewTasker()
			for _, task := range tt.tasks {
				tasker.AddTask(task)
			}
			got, err := tasker.Process(context.Background())
			if (err != nil) != tt.expectErr {
				t.Errorf("Tasker.Process() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tasker.Process() = %v, want %v", got, tt.want)
			}
			t.Logf("GOROUTINES: %d", runtime.NumGoroutine())
		})
	}
}
