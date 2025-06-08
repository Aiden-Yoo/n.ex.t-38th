package main

import (
	"context"
	"os"

	eapi "gornir_w_goeapi/tasks"
	"gornir_w_goeapi/util"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
)

func main() {
	log := logger.NewLogrus(false)

	file := "./config/inventory.yaml"
	filter := "arista"
	plugin := inventory.FromYAML{HostsFile: file}
	inv, err := plugin.Create()
	if err != nil {
		log.Fatal(err)
	}

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(runner.Parallel())

	// Applying filter (optional)
	gr = gr.Filter(func(h *gornir.Host) bool {
		groups, ok := h.Data["groups"].([]interface{})
		if !ok {
			return false
		}
		for _, g := range groups {
			if s, ok := g.(string); ok && s == filter {
				return true
			}
		}
		return false
	})

	// Open an SSH connection towards the devices
	results, err := gr.RunSync(
		context.Background(),
		// &eapi.TaskAll{
		// 	Tasks: []gornir.Task{
		// 		&eapi.TaskShowVersion{},
		// 		&eapi.TaskShowVlan{},
		// 	},
		// },
		// &eapi.TaskShowVersion{},
		// &eapi.TaskShowTest{},
		&eapi.TaskShowVlan{SaveToDB: true},
		// eapi.NewTaskShowAll(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 결과를 저장할 슬라이스
	var storedResults []*gornir.JobResult

	// 새 채널 생성 (RenderResults용)
	resultsCopy := make(chan *gornir.JobResult, 10) // 버퍼 크기는 적절히 조정

	// 채널에서 결과를 하나씩 읽어옴
	for result := range results {
		storedResults = append(storedResults, result)
	}

	// 저장된 결과를 새 채널에 복사
	go func() {
		for _, r := range storedResults {
			resultsCopy <- r
		}
		close(resultsCopy)
	}()

	// 커스텀 결과 출력
	util.CustomRenderResults(os.Stdout, storedResults, "장비 정보")

	// 결과 출력
	output.RenderResults(os.Stdout, resultsCopy, "장비 정보", true)
}
