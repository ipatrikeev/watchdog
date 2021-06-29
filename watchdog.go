package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ipatrikeev/watchdog/config"
	"github.com/ipatrikeev/watchdog/notify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var cfgPath = flag.String("config-path", "./config.yml", "Path to YAML application config")
var apiTimeout = 3 * time.Second

func main() {
	flag.Parse()

	cfg := config.AppConfig{}
	configContent, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("Can't read app config content at %s, error: %v", *cfgPath, err)
	}

	err = yaml.Unmarshal(configContent, &cfg)
	if err != nil {
		log.Fatalf("Couldn't construct app config: %v", err)
	}

	senders, err := parseSenders(cfg)
	if err != nil {
		log.Fatalf("Couldn't parse notifiers: %v", err)
	}
	notifier := notify.Notifier{Senders: senders}
	if err = notifier.Validate(); err != nil {
		log.Fatalf("Invalid notifiiers: %v", err)
	}

	for _, e := range cfg.Entities {
		go monitorSingle(e, notifier)
	}

	select {}
}

func monitorSingle(entity config.MonitoredEntity, notifier notify.Notifier) {
	fmt.Printf("Monitoring %v\n", &entity)
	for {
		processCheck(entity, notifier)
	}
}

func processCheck(entity config.MonitoredEntity, notifier notify.Notifier) {
	time.Sleep(entity.CheckPeriod)
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, entity.HealthUrl, nil)
	if err != nil {
		notifier.Fail(entity, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || !entity.CheckStatus(resp.StatusCode) {
		var info interface{}
		if err == nil {
			info = resp.Status
		} else {
			info = err
		}
		notifier.Fail(entity, info)
	} else {
		notifier.Success(entity)
	}
}

func parseSenders(cfg config.AppConfig) ([]notify.Sender, error) {
	var senders []notify.Sender

	for _, s := range cfg.Senders {
		switch strings.ToLower(s.Name) {
		case "telegram":
			sender, err := notify.NewTelegramSender(s.Params)
			if err != nil {
				return nil, err
			}
			senders = append(senders, sender)
		default:
			return nil, fmt.Errorf("unsupported sender type: %s", s.Name)
		}
	}

	return senders, nil
}
