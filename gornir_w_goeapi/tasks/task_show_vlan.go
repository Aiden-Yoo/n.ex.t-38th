package gornir

import (
	"context"
	"encoding/json"
	"fmt"
	"gornir_w_goeapi/util"
	"sort"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/nornir-automation/gornir/pkg/gornir"
)

type VlanInterfaces struct {
	Blocked         *bool `json:"blocked"`
	PrivatePromoted bool  `json:"privatePromoted"`
}

type Vlans struct {
	Dynamic    bool                      `json:"dynamic"`
	Interfaces map[string]VlanInterfaces `json:"interfaces"`
	Name       string                    `json:"name"`
	Status     string                    `json:"status"`
}

type VlanInfo struct {
	SourceDetail string           `json:"sourceDetail"`
	Vlans        map[string]Vlans `json:"vlans"`
}

type TaskShowVlan struct {
	SaveToDB bool
}

func (t *TaskShowVlan) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	// get host info
	hostname := host.Hostname
	username := host.Username
	password := host.Password

	// DB file path
	db, err := util.InitDB("results.db")
	if err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// initialize result struct
	info := VlanInfo{}

	// connect to device using goeapi
	node, err := goeapi.Connect("http", hostname, username, password, 80) // https, 443
	if err != nil {
		return info, fmt.Errorf("failed to connect to %s: %v", hostname, err)
	}

	// get version info
	commands := []string{"show vlan"}
	result, err := node.RunCommands(commands, "json")
	if err != nil {
		return info, fmt.Errorf("failed to execute command on %s: %v", hostname, err)
	}

	// check result
	if len(result.Result) == 0 {
		return info, fmt.Errorf("empty result from device")
	}

	versionData := result.Result[0]

	rawJSON, err := json.Marshal(versionData)
	if err != nil {
		return info, fmt.Errorf("marshal failed: %v", err)
	}

	// fmt.Println(string(rawJSON))

	if err := json.Unmarshal(rawJSON, &info); err != nil {
		return info, fmt.Errorf("unmarshal failed: %v", err)
	}

	// save to db
	if t.SaveToDB && db != nil {
		err = util.SaveCommandResult(db, hostname, "show vlan", versionData)
		if err != nil {
			fmt.Println("save failed:", err)
		}
	}

	return info, nil
}

func (t *TaskShowVlan) Metadata() *gornir.TaskMetadata {
	return &gornir.TaskMetadata{Identifier: "show_vlan"}
}

func (v VlanInfo) Format() string {
	result := fmt.Sprintf("show vlan: %s\n", v.SourceDetail)
	for vlanID, vlan := range v.Vlans {
		var intf string
		if len(vlan.Interfaces) > 0 {
			// Interface list -> slice
			interfaces := make([]string, 0, len(vlan.Interfaces))
			for ifName := range vlan.Interfaces {
				interfaces = append(interfaces, ifName)
			}
			// Interface name sorting
			sort.Strings(interfaces)

			intf = "(" + strings.Join(interfaces, ", ") + ")"
		}
		result += fmt.Sprintf("  - VLAN %s (%s): interfaces=%d %s\n", vlanID, vlan.Name, len(vlan.Interfaces), intf)
	}
	return result
}

func (v VlanInfo) String() string {
	result := fmt.Sprintf("\nSourceDetail: %s\n", v.SourceDetail)
	for vlanID, vlan := range v.Vlans {
		result += fmt.Sprintf("VLAN %s:\n", vlanID)
		result += fmt.Sprintf("  Name: %s\n", vlan.Name)
		result += fmt.Sprintf("  Status: %s\n", vlan.Status)
		result += fmt.Sprintf("  Dynamic: %t\n", vlan.Dynamic)
		result += fmt.Sprintln("  Interfaces:")
		for ifName, iface := range vlan.Interfaces {
			blockedStr := "null"
			if iface.Blocked != nil {
				blockedStr = fmt.Sprintf("%v", *iface.Blocked)
			}
			result += fmt.Sprintf("    - %s: blocked=%s, promoted=%t\n", ifName, blockedStr, iface.PrivatePromoted)
		}
	}
	return result
}
