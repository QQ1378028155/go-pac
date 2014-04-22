package cmd

import "errors"
import "os/exec"

func SyncCmd(name string, arg []string) error {
	if arg == nil {
		return errors.New("invalid args")
	}
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
