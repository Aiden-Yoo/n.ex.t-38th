package gornir

import (
	"context"
	"encoding/json"
	"fmt"
	"gornir_w_goeapi/util"

	"github.com/aristanetworks/goeapi"
	"github.com/nornir-automation/gornir/pkg/gornir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type InterfacesInfo struct {
	Interfaces map[string]interface{} `json:"interfaces"`
}

type TaskShowInterfaces struct {
	SaveToDB bool
}

func (t *TaskShowInterfaces) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	hostname := host.Hostname
	username := host.Username
	password := host.Password

	db, err := gorm.Open(sqlite.Open("results.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	info := InterfacesInfo{}

	node, err := goeapi.Connect("http", hostname, username, password, 80)
	if err != nil {
		return info, fmt.Errorf("failed to connect to %s: %v", hostname, err)
	}

	commands := []string{"show interfaces"}
	result, err := node.RunCommands(commands, "json")
	if err != nil {
		return info, fmt.Errorf("failed to execute command on %s: %v", hostname, err)
	}

	if len(result.Result) == 0 {
		return info, fmt.Errorf("empty result from device")
	}

	interfacesData := result.Result[0]

	rawJSON, err := json.Marshal(interfacesData)
	if err != nil {
		return info, fmt.Errorf("marshal failed: %v", err)
	}

	if err := json.Unmarshal(rawJSON, &info); err != nil {
		return info, fmt.Errorf("unmarshal failed: %v", err)
	}

	if t.SaveToDB && db != nil {
		err = util.SaveCommandResult(db, hostname, "show interfaces", interfacesData, "show_interfaces")
		if err != nil {
			fmt.Println("save failed:", err)
		}
	}

	return info, nil
}

func (t *TaskShowInterfaces) Metadata() *gornir.TaskMetadata {
	return &gornir.TaskMetadata{Identifier: "show_interfaces"}
}

func (i InterfacesInfo) Format() string {
	return fmt.Sprintf("show interfaces: %d개 인터페이스", len(i.Interfaces))
}

func (i InterfacesInfo) String() string {
	b, _ := json.MarshalIndent(i.Interfaces, "", "  ")
	return fmt.Sprintf("show interfaces 결과:\n%s", string(b))
}
