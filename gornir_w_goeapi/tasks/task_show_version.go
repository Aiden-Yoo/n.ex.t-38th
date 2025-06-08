package gornir

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aristanetworks/goeapi"
	"github.com/nornir-automation/gornir/pkg/gornir"
)

type VersionInfo struct {
	Version            string  `json:"version"`
	Model              string  `json:"modelName"`
	SerialNum          string  `json:"serialNumber"`
	SystemMac          string  `json:"systemMacAddress"`
	Architecture       string  `json:"architecture"`
	BootupTimestamp    float64 `json:"bootupTimestamp"`
	ConfigMacAddress   string  `json:"configMacAddress"`
	HardwareRevision   string  `json:"hardwareRevision"`
	HwMacAddress       string  `json:"hwMacAddress"`
	ImageFormatVersion string  `json:"imageFormatVersion"`
	ImageOptimization  string  `json:"imageOptimization"`
	InternalBuildId    string  `json:"internalBuildId"`
	InternalVersion    string  `json:"internalVersion"`
	IsIntlVersion      bool    `json:"isIntlVersion"`
	MemFree            uint64  `json:"memFree"`
	MemTotal           uint64  `json:"memTotal"`
	MfgName            string  `json:"mfgName"`
	Uptime             float64 `json:"uptime"`
}

type TaskShowVersion struct{}

func (t *TaskShowVersion) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	// get host info
	hostname := host.Hostname
	username := host.Username
	password := host.Password

	// initialize result struct
	info := VersionInfo{}

	// connect to device using goeapi
	node, err := goeapi.Connect("http", hostname, username, password, 80) // https, 443
	if err != nil {
		return info, fmt.Errorf("failed to connect to %s: %v", hostname, err)
	}

	// get version info
	commands := []string{"show version"}
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

	return info, nil
}

func (t *TaskShowVersion) Metadata() *gornir.TaskMetadata {
	return &gornir.TaskMetadata{Identifier: "show_version"}
}

func (i VersionInfo) Format() string {
	return fmt.Sprintf("show version:\n  - 버전: %s\n  - 모델: %s\n  - 시리얼 번호: %s\n  - 시스템 MAC: %s\n  - 아키텍처: %s\n",
		i.Version, i.Model, i.SerialNum, i.SystemMac, i.Architecture)
}

func (v VersionInfo) String() string {
	return fmt.Sprintf(
		`{Version: %s, Model: %s, Serial: %s, System MAC: %s, Architecture: %s, Bootup: %.0f, Config MAC: %s, HW Revision: %s, HW MAC: %s, Image Format: %s, Optimization: %s, Build ID: %s, Internal Ver: %s, Intl Version: %t, MemFree: %d, MemTotal: %d, Mfg: %s, Uptime: %.1f seconds (%.1f days)}`,
		v.Version, v.Model, v.SerialNum, v.SystemMac, v.Architecture, v.BootupTimestamp, v.ConfigMacAddress, v.HardwareRevision, v.HwMacAddress, v.ImageFormatVersion, v.ImageOptimization, v.InternalBuildId, v.InternalVersion, v.IsIntlVersion, v.MemFree, v.MemTotal, v.MfgName, v.Uptime, v.Uptime/86400,
	)
}
