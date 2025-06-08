package util

import (
	"encoding/json"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// VLAN 정보 저장용 모델
type VlanResult struct {
	ID         uint `gorm:"primaryKey"`
	Host       string
	VlanID     string
	VlanName   string
	Status     string
	Interfaces string // 인터페이스 목록을 콤마로 저장
	CreatedAt  time.Time
}

// Command 결과 저장용 모델
type CommandResult struct {
	ID        uint `gorm:"primaryKey"`
	Host      string
	Command   string // ex: "show vlan", "show version"
	Result    string // JSON 등 직렬화된 결과
	CreatedAt time.Time
}

// DB 초기화 및 마이그레이션
func InitDB(filepath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&CommandResult{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// VLAN 결과 저장 함수
func SaveVlanResult(db *gorm.DB, host, vlanID, vlanName, status string, interfaces []string) error {
	result := VlanResult{
		Host:       host,
		VlanID:     vlanID,
		VlanName:   vlanName,
		Status:     status,
		Interfaces: joinInterfaces(interfaces),
		CreatedAt:  time.Now(),
	}
	return db.Create(&result).Error
}

// 인터페이스 슬라이스를 콤마로 합침
func joinInterfaces(interfaces []string) string {
	return strings.Join(interfaces, ",")
}

func SaveCommandResult(db *gorm.DB, host, command string, result interface{}) error {
	// result를 JSON 문자열로 변환
	resultStr, err := toJSONString(result)
	if err != nil {
		return err
	}
	rec := CommandResult{
		Host:      host,
		Command:   command,
		Result:    resultStr,
		CreatedAt: time.Now(),
	}
	return db.Create(&rec).Error
}

func toJSONString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
