package gornir

import (
	"context"
	"fmt"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

type TaskAll struct {
	Tasks []gornir.Task
}

func (t *TaskAll) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	results := make([]interface{}, 0)
	for _, task := range t.Tasks {
		fmt.Printf("Running task: %s on host: %s\n", task.Metadata().Identifier, host.Hostname)

		result, err := task.Run(ctx, logger, host)
		if err != nil {
			return results, fmt.Errorf("error in task %s: %v", task.Metadata().Identifier, err)
		}
		results = append(results, result)
	}
	return results, nil
}

func (t *TaskAll) Metadata() *gornir.TaskMetadata {
	return &gornir.TaskMetadata{Identifier: "task_all"}
}

// temp
func NewTaskShowAll() (gornir.TaskInstanceResult, error) {
	return &TaskAll{
		Tasks: []gornir.Task{
			&TaskShowVersion{},
			// &TaskShowInterfaces{},
			&TaskShowVlan{},
		},
	}, nil
}
