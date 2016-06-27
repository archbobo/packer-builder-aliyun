// The aliyun package contains a packer.Builder implementation that
// builds images for Aliyun Platform
package aliyun

import (
	"github.com/mitchellh/multistep"
	"log"
	"github.com/mitchellh/packer/packer"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/common"
)

// The unique ID for this builder
const BuilderId = "packer.aliyun"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return warnings, errs
	}
	b.config = *c

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	client := ecs.NewClient(b.config.AccessKeyId, b.config.AccessKeySecret)


	// Setup the state
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps
	steps := []multistep.Step{
		new(stepCreateInstance),
		new(stepStartInstance),
		new(stepInstanceInfo),
		&communicator.StepConnect{
			Config: &b.config.Comm,
			Host: commHost,
			SSHConfig: sshConfig,
		},
		new(common.StepProvision),
		new(stepCreateSnapshot),
		new(stepCreateImage),
	}

	// Run the steps
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("image_id"); !ok {
		log.Println("Failed to find image_id in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		imageName: b.config.ImageName,
		imageId: state.Get("image_id").(string),
		regionId: b.config.RegionId,
		client: client,
	}

	return artifact, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}