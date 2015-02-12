package marine

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

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
	t.Skip()
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

func TestListingIP(t *testing.T) {
	t.Skip()

	log.Info("Start testing")
	base, err := Import(os.Getenv("GOPATH")+"/files/ubuntu-14.10-server-amd64.ova", 512)
	assert.NoError(t, err)
	boxes, err := base.Clone(1, "box")
	assert.NoError(t, err)

	defer func() {
		boxes[0].Stop()
		boxes[0].Remove()
		base.Remove()
	}()

	boxes[0].StartAndWait()
	if out, err := boxes[0].Sudo(`sed -i "\$aauto eth1\niface eth1 inet dhcp\n" /etc/network/interfaces`); err == nil {
		log.Info(string(out))
	}
	if out, err := boxes[0].Sudo("cat /etc/network/interfaces"); err == nil {
		log.Info(string(out))
	}
	if out, err := boxes[0].Sudo("ifup eth1"); err == nil {
		log.Info(string(out))
	}
	time.Sleep(3 * time.Second)
	if out, err := boxes[0].Run("/sbin/ip addr show dev eth1"); err == nil {
		log.Info("\n" + string(out))
	}

}

func TestGettingIP(t *testing.T) {
	t.Skip()
	log.Info("Start testing")
	base, err := Import(os.Getenv("GOPATH")+"/files/ubuntu-14.10-server-amd64.ova", 512)
	assert.NoError(t, err)
	boxes, err := base.Clone(1, "box")
	assert.NoError(t, err)

	defer func() {
		boxes[0].Stop()
		boxes[0].Remove()
		base.Remove()
	}()

	boxes[0].StartAndWait()
	assert.NoError(t, boxes[0].SetupIPAddr())
	ip, err := boxes[0].GetIPAddr()
	log.Infof("IP Address (eth1): %s", ip)
	assert.NoError(t, err)
	assert.Equal(t, true, strings.HasPrefix(ip, "192.168.99."))
}

func TestInstallDocker(t *testing.T) {
	log.Info("Start testing")
	// base, err := Import(os.Getenv("GOPATH")+"/files/ubuntu-14.10-server-amd64.ova", 512, "docker")
	// assert.NoError(t, err)
	base := &Machine{Name: "base"}
	boxes, err := base.Clone(1, "box")
	assert.NoError(t, err)

	defer func() {
		boxes[0].Stop()
		boxes[0].Remove()
	}()

	boxes[0].StartAndWait()
	out, err := boxes[0].Sudo(`bash -c "docker version"`)
	if err != nil {
		log.Error(err)
	}
	log.Info(out)
}
