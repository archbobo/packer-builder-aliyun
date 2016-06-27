package aliyun

import (
	"time"
	"log"
	"errors"
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

// Simply blocks until the snapshot progress is 100%,
// while eventually timing out.
func waitForSnapshotAccomplished(snapshotId string, regionId common.Region,
                             client *ecs.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts += 1

			log.Printf("Checking snapshot progress... (attempt: %d)", attempts)
			snapshots, _, err := client.DescribeSnapshots(&ecs.DescribeSnapshotsArgs{
				RegionId: regionId,
				SnapshotIds: {snapshotId},
			})
			if (err != nil) {
				result <- err
				return
			}
			if snapshots == nil || len(snapshots) != 1 {
				result <- errors.New("fail to get snapshot progress")
				return
			}

			// TODO check snapshot status?
			progress := snapshots[0].Progress
			log.Printf("Current snapshot progress: '%s'", progress)
			if progress == "100%" {
				result <- nil
				return
			}


			// Wait 10 seconds in between
			time.Sleep(10 * time.Second)

			// Verify we shouldn't exist
			select {
			case <- done:
			// We finished, so just exit the goroutine
				return
			default:
			// Keep going
			}

		}

	}()

	log.Printf("Waiting for up to %d seconds for snapshot to become accomplished", timeout/time.Second)
	select {
	case err := <- result:
		return err
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting for snapshot to become accomplished")
		return err
	}
}


// waitForStatus simply blocks until the instance is in
// a state we expect, while eventually timing out.
func waitForInstanceStatus(desiredStatus ecs.InstanceStatus, instanceId string, regionId common.Region,
                          client *ecs.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts += 1

			log.Printf("Checking instance status... (attempt: %d)", attempts)
			instances, _, err := client.DescribeInstances(&ecs.DescribeInstancesArgs{
				RegionId: regionId,
				InstanceIds: {instanceId},
			})
			if (err != nil) {
				result <- err
				return
			}
			if instances == nil || len(instances) != 1 {
				result <- errors.New("fail to get instance status")
				return
			}

			status := instances[0].Status
			log.Printf("Current instance status: '%s'", status)
			if status == desiredStatus {
				result <- nil
				return
			}

			// Wait 3 seconds in between
			time.Sleep(3 * time.Second)

			// Verify we shouldn't exist
			select {
			case <- done:
				// We finished, so just exit the goroutine
			        return
			default:
				// Keep going
			}
		}
	}()

	log.Printf("Waiting for up to %d seconds for instance to become %s", timeout/time.Second, desiredStatus)
	select {
	case err := <- result:
		return err
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting for instance to become '%s'", desiredStatus)
		return err
	}
}