package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/archcentric/packer-builder-aliyun/builder/aliyun"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(aliyun.Builder))
	server.Serve()
}
