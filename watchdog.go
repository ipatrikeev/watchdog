package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ipatrikeev/watchdog/config"
	"github.com/ipatrikeev/watchdog/notify"
	"gopkg.in/yaml.v2"
)

var cfgPath = flag.String("config-path", "./config.yml", "Path to YAML application config")
var apiTimeout = 10 * time.Second

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
		info := getFailInfo(resp, err)
		notifier.Fail(entity, info)
	} else {
		notifier.Success(entity)
	}
}

func getFailInfo(resp *http.Response, err error) interface{} {
	var info string

	if err == nil {
		// try to read response body as it may contain useful info
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			info = err.Error()
		} else {
			info = resp.Status + ": " + string(bodyBytes)
		}
	} else {
		info = err.Error()
	}
	return info
}

func parseSenders(cfg config.AppConfig) ([]notify.Sender, error) {
	var senders []notify.Sender

	for _, s := range cfg.Senders {
		var sender notify.Sender
		switch strings.ToLower(s.Name) {
		case "telegram":
			var err error
			sender, err = notify.NewTelegramSender(s.Params)
			if err != nil {
				return nil, err
			}
		case "console":
			sender = &notify.ConsoleSender{}
		default:
			return nil, fmt.Errorf("unsupported sender type: %s", s.Name)
		}
		senders = append(senders, sender)
	}

	return senders, nil
}
