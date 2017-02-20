package aliyun

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/common"
	aliyungo_common "github.com/denverdino/aliyungo/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/template/interpolate"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"log"
	"github.com/mitchellh/packer/common/uuid"
)

const defaultSshUsername = "root"
const defaultSshPassword = "Packer@123"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`

	RegionId aliyungo_common.Region `mapstructure:"region_id"`
	BaseImageId string `mapstructure:"base_image_id"`
	InstanceType string `mapstructure:"instance_type"`
	SecurityGroupId string `mapstructure:"security_group_id"`

	ImageName string `mapstructure:"image_name"`
	ImageDescription string `mapstructure:"image_description"`
	InstanceName string `mapstructure:"instance_name"`
	StatusTimeout      time.Duration `mapstructure:"status_timeout"`
        InternetChargeType aliyungo_common.InternetChargeType `mapstructure:"internet_charge_type"`

	ctx interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)

	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:           &md,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
			        "run_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Defaults
	if c.AccessKeyId == "" {
		// Default to environment variable for access_key_id, if it exists
		c.AccessKeyId = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if c.AccessKeySecret == "" {
		c.AccessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}

	if c.Comm.SSHUsername != "" && c.Comm.SSHUsername != defaultSshUsername {
		log.Printf(
			"on aliyun, for linux like os, only '%s' is accepted as ssh_username, overridding to '%s'",
			defaultSshUsername, defaultSshUsername)
	}
	c.Comm.SSHUsername = defaultSshUsername

	if c.StatusTimeout == 0 {
		// Default to 6 minute timeouts waiting for
		// desired status. i.e waiting for instance to become running.
		c.StatusTimeout = 6 * time.Minute
	}

	var errs *packer.MultiError

	if c.ImageName == "" {
		img, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("Unable to parse image name: %s ", err))
		} else {
			c.ImageName = img
		}
	}

	if c.ImageDescription == "" {
		c.ImageDescription = "Created by Packer for Aliyun"
	}

	if c.InstanceName == "" {
		c.InstanceName = fmt.Sprintf("packer-aliyun-%s", uuid.TimeOrderedUUID())
	}

	if c.Comm.SSHPassword == "" {
		c.Comm.SSHPassword = defaultSshPassword
	}

	if es := c.Comm.Prepare(&c.ctx); len(es) > 0 {
		errs = packer.MultiErrorAppend(errs, es...)
	}
	if c.AccessKeyId == "" {
		// Required configuration that will display errors if not set
		errs = packer.MultiErrorAppend(
			errs, errors.New("access_key_id for auth must be specified"))
	}
	if c.AccessKeySecret == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("access_key_secret for auth must be specified"))
	}


	if c.RegionId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("region_id is required"))
	}

	if c.BaseImageId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("base_image_id is required"))
	}

	if c.InstanceType == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("instance_type is required"))
	}

	if c.SecurityGroupId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("security_group_id is required"))
	}

	if c.InternetChargeType == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("internet_charge_type is required"))
	}

	log.Println(common.ScrubConfig(c, c.AccessKeyId, c.AccessKeySecret, c.Comm.SSHPassword))

	// Check for any errors.
	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}


	return c, nil, nil
}