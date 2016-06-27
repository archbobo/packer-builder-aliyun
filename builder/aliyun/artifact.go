package aliyun

import (
	"fmt"
	"log"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
)

type Artifact struct {
	// The name of the image
	imageName string

	// The ID of the image
	imageId string

	// The ID of the region
	regionId common.Region

	// The client for making API calls
	client *ecs.Client
}

// BuilderId returns the builder Id.
func (*Artifact) BuilderId() string {
	return BuilderId
}

// Destroy destroys the Aliyun image represented by the artifact.
func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: '%s' (ID: %s)", a.imageName, a.imageId)
	err := a.client.DeleteImage(a.regionId, a.imageId)
	return  err
}

// Files returns the files represented by the artifact.
func (*Artifact) Files() []string {
	// No files with Aliyun
	return nil
}

func (a *Artifact) Id() string {
	return a.imageId
}

// String returns the string representation of the artifact.
func (a *Artifact) String() string {
	return fmt.Sprintf("A disk image was created: '%v' (ID: %v) in region '%v'", a.imageName, a.imageId, a.regionId)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}