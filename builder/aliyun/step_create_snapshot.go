package aliyun

import (
	"github.com/mitchellh/multistep"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/mitchellh/packer/packer"
	"fmt"
	"time"
)

type stepCreateSnapshot struct {
	snapshotId string
}

func (s *stepCreateSnapshot) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	diskId := state.Get("disk_id").(string)

	ui.Say("Creating disk snapshot...")
	snapshotId, err := client.CreateSnapshot(&ecs.CreateSnapshotArgs{
		DiskId: diskId,
	})

	if err != nil {
		err := fmt.Errorf("Error creating snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if snapshotId == "" {
		err := fmt.Errorf("Error creating snapshot, snapshotId is empty")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Wait for the snapshot to become accomplished. For snapshots
	// this can end up taking quite a long time, so we hardcode this to
	// 20 minutes.
	err = waitForSnapshotAccomplished(snapshotId, c.RegionId, client, 20*time.Minute)
	if err != nil {
		err := fmt.Errorf("Error waiting for snapshot to become accomplished: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.snapshotId = snapshotId

	// Store the snapshot id for later
	state.Put("snapshot_id", snapshotId)

	return multistep.ActionContinue
}

func (s *stepCreateSnapshot) Cleanup(state multistep.StateBag) {
	// if the snapshot isn't there, we probably never created it
	if s.snapshotId == "" {
		return
	}

	client := state.Get("client").(*ecs.Client)
	ui := state.Get("ui").(packer.Ui)

	// Delete the snapshot we just created
	ui.Say("Delete snapshot...")
	err := client.DeleteSnapshot(s.snapshotId)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error deleting snapshot: %s, please delete it manually", err))
	}
}
