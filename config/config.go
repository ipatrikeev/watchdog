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
}

func (entity *MonitoredEntity) String() string {
	return fmt.Sprintf("%s (%s) checking every %v",
		entity.Name, entity.HealthUrl, entity.CheckPeriod)
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
	Senders  []MessageSender
}
