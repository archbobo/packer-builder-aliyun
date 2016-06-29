package aliyun

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/denverdino/aliyungo/ecs"
	"fmt"
)

type stepStartInstance struct{}

func (s *stepStartInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	instanceId := state.Get("instance_id").(string)

	ui.Say("Starting the instance...")
	err := client.StartInstance(instanceId)
	if err != nil {
		err := fmt.Errorf("Error starting instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Waiting for instance to become running status...")

	err = waitForInstanceStatus(ecs.Running, instanceId, c.RegionId, client, c.StatusTimeout)
	if err != nil {
		err := fmt.Errorf("Error waiting for instance to become running: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Instance has been started!")

	return multistep.ActionContinue

}

func (s *stepStartInstance) Cleanup(state multistep.StateBag) {
	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	instanceId := state.Get("instance_id").(string)

	// Stop the instance we just started
	ui.Say("Stop instance...")
	err := client.StopInstance(instanceId, true)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error stopping instance: %s", err))
	} else {
		err = waitForInstanceStatus(ecs.Stopped, instanceId, c.RegionId, client, c.StatusTimeout)
		if err != nil {
			ui.Error(fmt.Sprintf("Error waiting for instance to become stopped: %s", err))
		} else {
			ui.Message("Instance has been stopped!")
		}
	}

}
