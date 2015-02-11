package marine

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Machine struct {
	Name           string
	ID             string
	ForwardingPort string
}

func (m *Machine) Clone(num int, prefix string) ([]*Machine, error) {
	err := Clone(m.Name, prefix, num, "vboxnet0")
	if err != nil {
		return nil, err
	}
	machines := make([]*Machine, num)
	for i := 1; i <= num; i++ {
		name := fmt.Sprintf("%s%03d", prefix, i)
		machines[i-1] = &Machine{Name: name}
	}
	return machines, nil
}

func (m *Machine) StartAndWait() error {
	out, err := exec.Command(VBOX_MANAGE, "showvminfo", m.Name, "--machinereadable").Output()
	if err != nil {
		return fmt.Errorf("Cannot get vminfo %s", m.Name)
	}
	port := ""
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Forwarding") {
			re := regexp.MustCompile(`Forwarding\(\d+\)="ssh,tcp,[0-9\.]*,(\d+),[0-9\.]*,\d+"`)
			result := re.FindStringSubmatch(line)
			if len(result) == 2 {
				m.ForwardingPort = result[1]
				break
			}
		}
	}
	if port == "" {
		return fmt.Errorf("Cannot find port: %s", m.Name)
	}
	log.Infof("Found %s = %s", m.Name, m.ForwardingPort)
	return StartAndWait(m.Name, m.ForwardingPort)
}

func (m *Machine) Run(cmd string) (string, error) {
	_, sess, err := ConnectToHost("ubuntu@127.0.0.1:"+m.ForwardingPort, "reverse")
	defer sess.Close()
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (m *Machine) Sudo(cmd string) (string, error) {
	// "/bin/bash -c 'echo reverse | sudo -S whoami'"
	sudo := fmt.Sprintf("/bin/bash -c 'echo reverse | sudo -S %s'", cmd)
	return m.Run(sudo)
}

func (m *Machine) Remove() error {
	return Remove(m.Name)
}

func (m *Machine) Stop() error {
	return Stop(m.Name)
}
