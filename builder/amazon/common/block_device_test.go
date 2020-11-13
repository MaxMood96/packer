package common

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/packer/packer-plugin-sdk/config"
)

func TestBlockDevice(t *testing.T) {
	cases := []struct {
		Config *BlockDevice
		Result *ec2.BlockDeviceMapping
	}{
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				SnapshotId:          "snap-1234",
				VolumeType:          "standard",
				VolumeSize:          8,
				DeleteOnTermination: true,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId:          aws.String("snap-1234"),
					VolumeType:          aws.String("standard"),
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName: "/dev/sdb",
				VolumeSize: 8,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(false),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				VolumeType:          "io1",
				VolumeSize:          8,
				DeleteOnTermination: true,
				IOPS:                1000,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("io1"),
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(true),
					Iops:                aws.Int64(1000),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				VolumeType:          "io2",
				VolumeSize:          8,
				DeleteOnTermination: true,
				IOPS:                1000,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("io2"),
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(true),
					Iops:                aws.Int64(1000),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				VolumeType:          "gp2",
				VolumeSize:          8,
				DeleteOnTermination: true,
				Encrypted:           config.TriTrue,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(true),
					Encrypted:           aws.Bool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				VolumeType:          "gp2",
				VolumeSize:          8,
				DeleteOnTermination: true,
				Encrypted:           config.TriTrue,
				KmsKeyId:            "2Fa48a521f-3aff-4b34-a159-376ac5d37812",
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(8),
					DeleteOnTermination: aws.Bool(true),
					Encrypted:           aws.Bool(true),
					KmsKeyId:            aws.String("2Fa48a521f-3aff-4b34-a159-376ac5d37812"),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:          "/dev/sdb",
				VolumeType:          "standard",
				DeleteOnTermination: true,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("standard"),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:  "/dev/sdb",
				VirtualName: "ephemeral0",
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName:  aws.String("/dev/sdb"),
				VirtualName: aws.String("ephemeral0"),
			},
		},
		{
			Config: &BlockDevice{
				DeviceName: "/dev/sdb",
				NoDevice:   true,
			},

			Result: &ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/sdb"),
				NoDevice:   aws.String(""),
			},
		},
	}

	for _, tc := range cases {
		var amiBlockDevices BlockDevices = []BlockDevice{*tc.Config}

		var launchBlockDevices BlockDevices = []BlockDevice{*tc.Config}

		expected := []*ec2.BlockDeviceMapping{tc.Result}

		amiResults := amiBlockDevices.BuildEC2BlockDeviceMappings()
		if diff := cmp.Diff(expected, amiResults); diff != "" {
			t.Fatalf("Bad block device: %s", diff)
		}

		launchResults := launchBlockDevices.BuildEC2BlockDeviceMappings()
		if diff := cmp.Diff(expected, launchResults); diff != "" {
			t.Fatalf("Bad block device: %s", diff)
		}
	}
}
