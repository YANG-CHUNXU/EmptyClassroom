package config

import (
	"EmptyClassroom/logs"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	ConfigPathKey = "CONFIG_PATH"
)

type CampusConfig struct {
	Name         string `json:"name"`
	Id           int    `json:"id,omitempty"`
	HasRealtime  bool   `json:"has_realtime"`
	ReplaceRegex []struct {
		Regex   string `json:"regex"`
		Replace string `json:"replace"`
	} `json:"replace_regex"`
}

type NotificationConfig struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	Duration         int    `json:"duration"`
	Type             string `json:"type"`
	ShowNotification bool   `json:"showNotification"`
	Start            string `json:"start"`
	End              string `json:"end"`
}

type Config struct {
	ClassTable   ClassTableConfig   `json:"class_table"`
	Campus       []CampusConfig     `json:"campus"`
	Notification NotificationConfig `json:"notification"`
}

type ClassTableConfig struct {
	StartWeek     string                `json:"start_week"`
	EndWeek       string                `json:"end_week"`
	UnableReason  string                `json:"unable_reason"`
	IsAvailable   bool                  `json:"is_available"`
	ClassTableMap map[string]ClassTable `json:"class_table_map"`
}

type ClassTable struct {
	Class   []ClassTableClassroomInfo `json:"class"`
	TypeMap map[string]string         `json:"typeMap"`
}

type ClassTableClassroomInfo struct {
	Campus  string    `json:"campus"`
	Seat    string    `json:"seat"`
	Name    string    `json:"name"`
	Classes [][][]int `json:"classes"`
}

var GlobalConfig *Config

//go:embed data/*.json
var embeddedConfigFS embed.FS

func loadConfigFile(configPath string, filename string) ([]byte, error) {
	if configPath != "" {
		fullPath := filepath.Join(configPath, filename)
		configContent, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}
		return configContent, nil
	}

	configContent, err := fs.ReadFile(embeddedConfigFS, filepath.ToSlash(filepath.Join("data", filename)))
	if err != nil {
		return nil, err
	}
	return configContent, nil
}

func mustLoadConfigFile(configPath string, filename string) []byte {
	configContent, err := loadConfigFile(configPath, filename)
	if err != nil {
		logs.CtxError(context.Background(), "load config file failed: %v", err)
		panic(err)
	}
	return configContent
}

func InitConfig() {
	configPath := os.Getenv(ConfigPathKey)
	if configPath != "" {
		info, err := os.Stat(configPath)
		if err != nil {
			logs.CtxError(context.Background(), "stat config directory failed: %v", err)
			panic(err)
		}
		if !info.IsDir() {
			err = errors.New("CONFIG_PATH is not a directory")
			logs.CtxError(context.Background(), "%v", err)
			panic(err)
		}
	}

	configContent := mustLoadConfigFile(configPath, "config.json")
	GlobalConfig = new(Config)
	err := json.Unmarshal(configContent, GlobalConfig)
	if err != nil {
		logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
		panic(err)
	}
	for _, building := range GlobalConfig.Campus {
		configContent = mustLoadConfigFile(configPath, building.Name+".json")
		buildingConfig := new(ClassTable)
		err = json.Unmarshal(configContent, buildingConfig)
		if err != nil {
			logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
			panic(err)
		}
		if GlobalConfig.ClassTable.ClassTableMap == nil {
			GlobalConfig.ClassTable.ClassTableMap = make(map[string]ClassTable)
		}
		GlobalConfig.ClassTable.ClassTableMap[building.Name] = *buildingConfig
	}

	configContent = mustLoadConfigFile(configPath, "notification.json")
	notificationConfig := new(NotificationConfig)
	err = json.Unmarshal(configContent, notificationConfig)
	if err != nil {
		logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
		panic(err)
	}
	GlobalConfig.Notification = *notificationConfig
}

func GetConfig() Config {
	if GlobalConfig == nil {
		InitConfig()
	}
	return *GlobalConfig
}
