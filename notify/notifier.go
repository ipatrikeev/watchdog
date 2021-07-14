package notify

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ipatrikeev/watchdog/config"
)

var failTemplate = "❌ %s fail: %v"
var recoverTemplate = "✅ %s recover"

type Notifier struct {
	Senders []Sender
}

func (n *Notifier) Fail(entity config.MonitoredEntity, details interface{}) {
	if n.shouldNotify(entity, false) {
		msg := fmt.Sprintf(failTemplate, entity.Name, details)
		n.notifyAll(msg)
	}
}

func (n *Notifier) Success(entity config.MonitoredEntity) {
	if n.shouldNotify(entity, true) {
		msg := fmt.Sprintf(recoverTemplate, entity.Name)
		n.notifyAll(msg)
	}
}

func (n *Notifier) notifyAll(msg string) {
	for _, s := range n.Senders {
		s.Send(msg)
	}
}

func (n *Notifier) Validate() error {
	if len(n.Senders) == 0 {
		return errors.New("no notifiers specified")
	}
	return nil
}

// shouldNotify checks whether the message is needed and won't be repeated several times
func (n *Notifier) shouldNotify(entity config.MonitoredEntity, success bool) (res bool) {
	fileName := statusFileName(entity)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if success {
			return false
		} else {
			if _, err := n.incrementFails(fileName, entity); err != nil {
				return true
			}
			return entity.FailsAllowed < 1
		}
	} else {
		if success {
			fails, err := n.removeStatusFile(fileName, entity)
			if err != nil {
				return true
			}
			return fails > entity.FailsAllowed
		} else {
			fails, err := n.incrementFails(fileName, entity)
			if err != nil {
				return true
			}
			// notify on first fails overflow
			if fails == entity.FailsAllowed+1 {
				return true
			}
			return false
		}
	}
}

func statusFileName(entity config.MonitoredEntity) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	hasher := sha1.New()
	hasher.Write([]byte(entity.Name))
	hashedName := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return exPath + "/" + hashedName
}

func (n *Notifier) readCurrentFails(fileName string, entity config.MonitoredEntity) (int, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Couldn't read status file for %s, file path %s, error: %v",
			entity.Name, fileName, err)
		return 1, err
	}
	prevFails, err := strconv.Atoi(string(content))
	if err != nil {
		fmt.Printf("Couldn't parse fail count: %v, error: %v", content, err)
		return 1, err
	}
	return prevFails, nil
}

func (n *Notifier) closeFile(f *os.File, fileName string) {
	err := f.Close()
	if err != nil {
		n.handleFailure(fmt.Sprintf("Couldn't close status file: %s, error: %v", fileName, err))
	}
}

func (n *Notifier) removeStatusFile(fileName string, entity config.MonitoredEntity) (int, error) {
	fails, _ := n.readCurrentFails(fileName, entity)

	err := os.Remove(fileName)
	if err != nil {
		n.handleFailure(fmt.Sprintf("Couldn't remove status file for %v, file path: %s", entity.Name, fileName))
	}

	return fails, err
}

func (n *Notifier) incrementFails(fileName string, entity config.MonitoredEntity) (int, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := n.writeFailCount(fileName, 1)
		if err != nil {
			n.handleFailure(fmt.Sprintf("Couldn't write to status file for %v, file path: %s, err: %v",
				entity.Name, fileName, err))
			return -1, err
		}
		return 1, nil
	} else {
		prevFails, err := n.readCurrentFails(fileName, entity)
		if err != nil {
			n.handleFailure(fmt.Sprintf("Couldn't read current fail count for %s, err: %v", entity.Name, err))
			return -1, err
		}
		fails := prevFails + 1

		err = n.writeFailCount(fileName, fails)
		if err != nil {
			n.handleFailure(fmt.Sprintf("Couldn't write fail count for %s, err: %v", entity.Name, err))
			return -1, err
		}
		return fails, err
	}
}

func (n *Notifier) writeFailCount(fileName string, fails int) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer n.closeFile(f, fileName)
	if err != nil {
		return err
	}

	_, err = f.WriteString(strconv.Itoa(fails))
	return err
}

func (n *Notifier) handleFailure(err string) {
	n.notifyAll(err)
}
