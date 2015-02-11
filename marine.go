package marine

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var VBOX_MANAGE = "VBoxManage"

func Import(file string, memory int) error {
	name := "base"
	cmd := exec.Command(VBOX_MANAGE, "import", file,
		"--vsys", "0", "--vmname", name,
		"--memory", fmt.Sprintf("%d", memory),
	)
	out, err := cmd.Output()
	if err == nil {
		log.Infof("Imported \"%s\"", name)
	}
	if err != nil {
		log.Errorf("Error: %s\n%s", err, string(out))
		return err
	}

	_, err = exec.Command(VBOX_MANAGE, "snapshot", "base", "take", "origin").Output()
	log.Info("Snapshot \"origin\" taken")
	return err
}

func Modify(name string, adapter string, i int) error {
	err := exec.Command(VBOX_MANAGE, "modifyvm", name,
		"--natpf1", fmt.Sprintf("ssh,tcp,127.0.0.1,%d,,22", 52200+i),
		"--nic2", "hostonly",
		"--hostonlyadapter2", adapter,
		"--cableconnected2", "on",
		"--nicpromisc2", "allow-vms",
	).Run()
	if err == nil {
		log.Infof("Modified nic2 for \"%s\"", name)
	}
	return err
}

func Clone(baseName string, prefix string, num int, adapter string) error {
	for i := 1; i <= num; i++ {
		name := fmt.Sprintf("%s%03d", prefix, i)
		cmd := exec.Command(VBOX_MANAGE, "clonevm",
			baseName,
			"--snapshot", "origin",
			"--options", "link",
			"--name", name,
			"--register")
		out, err := cmd.Output()
		if err != nil {
			return err
		} else {
			err = Modify(name, adapter, i)
			log.Infof("Clone: %s", strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func Remove(args ...string) error {
	for _, name := range args {
		if name == "base" {
			err := exec.Command(VBOX_MANAGE, "snapshot", "base", "delete", "origin").Run()
			if err != nil {
				log.Info("Removed snapshot \"base/origin\"")
			}
		}
		cmd := exec.Command(VBOX_MANAGE, "unregistervm", name, "--delete")
		_, err := cmd.Output()
		log.Infof("Removed \"%s\"", name)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartAndWait(name string, port string) error {
	err := exec.Command(VBOX_MANAGE, "startvm", name).Run()
	log.Infof("Started \"%s\"", name)
	if err != nil {
		return err
	}
	err = WaitForTCP("127.0.0.1:" + port)
	if err == nil {
		log.Infof("VM \"%s\" ready to connect", name)
	}
	return err
}

func Stop(name string) error {
	err := exec.Command(VBOX_MANAGE, "controlvm", name, "acpipowerbutton").Run()
	if err == nil {
		log.Infof("Stopping \"%s\"", name)
	}

	for {
		st := GetState(name)
		if st == "poweroff" {
			log.Infof("VM \"%s\" is now %s", name, st)
			break
		} else if st == "error" {
			return fmt.Errorf("GetState: %s", st)
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func GetState(name string) string {
	out, err := exec.Command(VBOX_MANAGE, "showvminfo", name, "--machinereadable").Output()
	if err != nil {
		return "exec error"
	}
	str := string(out)
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VMState=") {
			v := strings.Split(line, "=")[1]
			return v[1 : len(v)-1]
		}
	}
	return "unknown error"
}

func WaitForTCP(addr string) error {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		defer conn.Close()
		if _, err = conn.Read(make([]byte, 1)); err != nil {
			continue
		}
		break
	}
	return nil
}
