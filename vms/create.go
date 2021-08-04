// B"H
package vms

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	monitorWriteScope = "https://www.googleapis.com/auth/monitoring.write"
)

// https://gitee.com/arohat/google-cloud-go/blob/v0.34.0/profiler/proftest/proftest.go

// StartInstance starts a GCE Instance with configs specified by inst,
// and which runs the startup script specified in inst. If image project
// is not specified, it defaults to "debian-cloud". If image family is
// not specified, it defaults to "debian-9".
func CreateInstance(computeService *compute.Service, ctx context.Context, inst *InstanceConfig) error {
	imageProject, imageFamily := inst.ImageProject, inst.ImageFamily
	if imageProject == "" {
		imageProject = "debian-cloud"
	}
	if imageFamily == "" {
		imageFamily = "debian-10"
	}
	img, err := computeService.Images.GetFromFamily(imageProject, imageFamily).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get image from family %q in project %q: %v", imageFamily, imageProject, err)
	}

	op, err := computeService.Instances.Insert(inst.ProjectID, inst.Zone, &compute.Instance{
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", inst.Zone, inst.MachineType),
		Name:        inst.Name,
		Disks: []*compute.AttachedDisk{{
			AutoDelete: true, // delete the disk when the VM is deleted.
			Boot:       true,
			Type:       "PERSISTENT",
			Mode:       "READ_WRITE",
			InitializeParams: &compute.AttachedDiskInitializeParams{
				SourceImage: img.SelfLink,
				DiskType:    fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/diskTypes/pd-standard", inst.ProjectID, inst.Zone),
			},
		}},
		NetworkInterfaces: []*compute.NetworkInterface{{
			Network: fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/default", inst.ProjectID),
			AccessConfigs: []*compute.AccessConfig{{
				Name:        "External NAT",
				NetworkTier: "STANDARD",
			}},
			Name: inst.Name + "-interface",
		}},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{{
				Key:   "startup-script",
				Value: googleapi.String(inst.StartupScript),
			}},
		},
		ServiceAccounts: []*compute.ServiceAccount{{
			Email:  "default",
			Scopes: append(inst.Scopes, monitorWriteScope),
		}},
	}).Do()

	if err != nil {
		return fmt.Errorf("failed to create instance: %v", err)
	}

	return waitForOperation(computeService, ctx, inst.ProjectID, inst.Zone, op)
}
