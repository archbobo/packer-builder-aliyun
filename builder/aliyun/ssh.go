package aliyun

import (
	"github.com/mitchellh/multistep"
	gossh "golang.org/x/crypto/ssh"
	"github.com/mitchellh/packer/communicator/ssh"
)

func commHost(state multistep.StateBag) (string, error) {
	ipAddress := state.Get("inner_ip").(string)
	return ipAddress, nil
}

func sshConfig(state multistep.StateBag) (*gossh.ClientConfig, error) {
	config := state.Get("config").(Config)

	auth := []gossh.AuthMethod{
		gossh.Password(config.Comm.SSHPassword),
		gossh.KeyboardInteractive(
			ssh.PasswordKeyboardInteractive(config.Comm.SSHPassword)),
	}

	return &gossh.ClientConfig{
		User: config.Comm.SSHUsername,
		Auth: auth,
	}, nil
}


