package aliyun

import (
	"github.com/mitchellh/multistep"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/mitchellh/packer/packer"
	"fmt"
)

type stepCreateInstance struct {
	instanceId string
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	// Create the instance based on configuration
	ui.Say("Creating instance...")
	instanceId, err := client.CreateInstance(&ecs.CreateInstanceArgs{
		RegionId: c.RegionId,
		ImageId: c.BaseImageId,
		InstanceType: c.InstanceType,
		SecurityGroupId: c.SecurityGroupId,
		InstanceName: c.InstanceName,
		Password: c.Comm.SSHPassword,
	})

	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// We use this in cleanup
	s.instanceId = instanceId

	ui.Message("Instance has been created!")
	// Store the instance id for later
	state.Put("instance_id", instanceId)

	/*
	ui.Say("Allocating public ip address...")
	publicIpAddress, err := client.AllocatePublicIpAddress(instanceId)

	if err != nil {
		err := fmt.Errorf("Error allocating public ip address: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Public IP allocated: %s", publicIpAddress))
	*/

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	// if the instanceid isn't there, we probably never created it
	if s.instanceId == "" {
		return
	}

	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)

	// Delete the instance we just created
	ui.Say("Delete instance...")
	err := client.DeleteInstance(s.instanceId)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error deleting instance. Please delete it manually: %s", err))
	} else {
		ui.Message("Instance has been deleted!")
	}

}