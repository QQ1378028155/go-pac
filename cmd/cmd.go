package cmd

import "errors"
import "os/exec"
import "github.com/QQ1378028155/go-pac/logger"

//SyncCmd can run a command synchronically
func SyncCmd(name string, arg []string) error {
	if arg == nil {
		return errors.New("invalid args")
	}
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	logger.Debug(string(b))
	if err != nil {
		return err
	}
	return nil
}
