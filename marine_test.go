package marine

import (
	"bytes"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSudo(t *testing.T) {
	t.Skip()
	log.Info("Start testing")
	_, err := Import(os.Getenv("GOPATH")+"/files/ubuntu-14.10-server-amd64.ova", 512)
	assert.NoError(t, err)
	assert.NoError(t, Clone("base", "box", 4, "vboxnet0"))
	StartAndWait("box001", "52201")
	defer func() {
		Stop("box001")
		Remove("box001", "box002", "box003", "box004", "base")
	}()

	_, sess, err := ConnectToHost("ubuntu@127.0.0.1:52201", "reverse")
	defer sess.Close()
	assert.NoError(t, err)

	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run("/bin/bash -c 'echo reverse | sudo -S whoami'"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	whoami := strings.TrimSpace(b.String())
	assert.Equal(t, "root", whoami)
	log.Infof("whoami: %s", whoami)
}

func TestNewAPIs(t *testing.T) {
	log.Info("Start testing")
	base, err := Import(os.Getenv("GOPATH")+"/files/ubuntu-14.10-server-amd64.ova", 512)
	assert.NoError(t, err)
	boxes, err := base.Clone(4, "box")
	assert.NoError(t, err)

	defer func() {
		boxes[0].Stop()
		for _, box := range boxes {
			assert.NoError(t, box.Remove())
		}
		assert.NoError(t, base.Remove())
	}()

	boxes[0].StartAndWait()
	if out, err := boxes[0].Run("/usr/bin/whoami"); err == nil {
		assert.Equal(t, "ubuntu", strings.TrimSpace(string(out)))
	}
	if out, err := boxes[0].Sudo("whoami"); err == nil {
		assert.Equal(t, "root", strings.TrimSpace(out))
	}
}
