package notify

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ipatrikeev/watchdog/config"
	"os"
	"path/filepath"
)

var failTemplate = "❌ %s fail: %v"
var recoverTemplate = "✅ %s recover"

type Notifier struct {
	Senders []Sender
}

func (n *Notifier) Fail(entity config.MonitoredEntity, details interface{}) {
	if shouldNotify(entity, false) {
		msg := fmt.Sprintf(failTemplate, entity.Name, details)
		n.notifyAll(msg)
	}
}

func (n *Notifier) Success(entity config.MonitoredEntity) {
	if shouldNotify(entity, true) {
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
		return errors.New("no senders specified")
	}
	return nil
}

// shouldNotify checks whether the message is needed and won't be repeated several times
func shouldNotify(entity config.MonitoredEntity, success bool) bool {
	fileName := statusFileName(entity)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if success {
			return false
		} else {
			_, err = os.Create(fileName)
			if err != nil {
				fmt.Printf("Couldn't create status file for %v, file path: %s\n", entity.Name, fileName)
			}
			return true
		}
	} else {
		if success {
			err = os.Remove(fileName)
			if err != nil {
				fmt.Printf("Couldn't remove status file for %v, file path: %s\n", entity.Name, fileName)
			}
			return true
		} else {
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
