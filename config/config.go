package config

import (
	"fmt"
	"time"
)

type MonitoredEntity struct {
	Name          string
	HealthUrl     string        `yaml:"health-url"`
	CheckPeriod   time.Duration `yaml:"check-period"`
	ValidStatuses []int         `yaml:"valid-statuses"`
	FailsAllowed  int           `yaml:"fails-allowed"`
}

func (entity *MonitoredEntity) String() string {
	failAllowedInfo := ""
	if entity.FailsAllowed > 1 {
		failAllowedInfo = fmt.Sprintf(". Won't notify unless %d fails happen in a row", entity.FailsAllowed)
	}
	return fmt.Sprintf("%s (%s) checking every %v%s",
		entity.Name, entity.HealthUrl, entity.CheckPeriod, failAllowedInfo)
}

func (entity *MonitoredEntity) CheckStatus(status int) bool {
	for _, s := range entity.ValidStatuses {
		if s == status {
			return true
		}
	}
	return false
}

type MessageSender struct {
	Name   string
	Params map[string]string
}

type AppConfig struct {
	Entities []MonitoredEntity
	Senders  []MessageSender `yaml:"notifiers"`
}
