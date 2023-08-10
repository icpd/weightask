# Priority Task
## Introduction
This Go-based project provides a system for managing and executing tasks based on their priority. Specifically, it concurrently executes a collection of tasks and returns the result of the one with the highest priority.
## How It Works
In the context of the Priority Task Processor, a "Task" is defined as any operation with a defined Priority and an Execute function. Here's what these mean:
- **Priority**: An integer indicating relative importance. Higher numbers indicate higher priority.  
- **Execute**: The function that performs the given task.    

Task processing involves two major steps:  
1. **Adding Tasks**: Tasks to be executed are added into the Priority Task Processor using the AddTask method. Each task is an instance of a struct that implements the Task interface.
2. **Processing Tasks**: The Process function executes all tasks concurrently and retrieves the result from the task with the highest priority. If multiple tasks share the same highest priority, it will return the first result.  
## How to Use
Here's a basic example of how to use the Priority Task Processor:  
1. Implement the Task interface in your task struct.
2. Initiate a new instance of Tasker via NewTasker().
3. Add tasks into the Tasker via AddTask.
4. Call the Process method to begin executing tasks. This will return the result from the highest-priority task.
```go
package main

import (
    "github.com/icpd/priortask"
)


type MyTask struct {
    ...
}

func (t *MyTask) Execute() (any, error) {
     ...
}

func (t *MyTask) Priority() int {
     ...
}

func main() {
    tasker := priortask.NewTasker()
    
    task := &MyTask{...}
    tasker.AddTask(task)
    ...
    
    result, err := tasker.ProcessTasks(context.Background())
    ...
}
```