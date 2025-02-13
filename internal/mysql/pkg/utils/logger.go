package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type ModuleLogger struct {
	Logs   LogsStruct  `json:"logs"`
	Result interface{} `json:"result"`
}

type LogsStruct struct {
	Error []string `json:"Error"`
	Info  []string `json:"Info"`
	Debug []string `json:"Debug"`
}

func NewModuleLogger() *ModuleLogger {
	return &ModuleLogger{
		Logs: LogsStruct{
			Error: make([]string, 0),
			Info:  make([]string, 0),
			Debug: make([]string, 0),
		},
	}
}

func (ml *ModuleLogger) AddError(message string) {
	ts := time.Now().Format(time.RFC3339)
	ml.Logs.Error = append(ml.Logs.Error, fmt.Sprintf("%s ERROR: %s", ts, message))
}

func (ml *ModuleLogger) AddInfo(message string) {
	ts := time.Now().Format(time.RFC3339)
	ml.Logs.Info = append(ml.Logs.Info, fmt.Sprintf("%s INFO: %s", ts, message))
}

// AddDebug ajoute un message de debug avec timestamp.
func (ml *ModuleLogger) AddDebug(message string) {
	ts := time.Now().Format(time.RFC3339)
	ml.Logs.Debug = append(ml.Logs.Debug, fmt.Sprintf("%s DEBUG: %s", ts, message))
}

// SetResult définit le résultat final du module.
func (ml *ModuleLogger) SetResult(result interface{}) {
	ml.Result = result
}

// JSON sérialise le logger en JSON.
func (ml *ModuleLogger) JSON() (string, error) {
	data, err := json.Marshal(ml)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var defaultLogger = NewModuleLogger()

func GetLogger() *ModuleLogger {
	return defaultLogger
}
