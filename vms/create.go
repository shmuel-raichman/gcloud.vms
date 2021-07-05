package vms

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	monitorWriteScope = "https://www.googleapis.com/auth/monitoring.write"
)

// TestRunner has common elements used for testing profiling agents on a range
// of environments.
// type TestRunner struct {
// 	Client *http.Client
// }

// // // GCETestRunner supports testing a profiling agent on GCE.
// type GCETestRunner struct {
// 	TestRunner
// 	ComputeService *compute.Service
// }

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
				Name: "External NAT",
			}},
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

	// Poll status of the operation to create the instance.
	// getOpCall := computeService.ZoneOperations.Get(inst.ProjectID, inst.Zone, op.Name)
	// for {
	// 	if err := checkOpErrors(op); err != nil {
	// 		return fmt.Errorf("failed to create instance: %v", err)
	// 	}
	// 	if op.Status == "DONE" {
	// 		return nil
	// 	}

	// 	if err := gax.Sleep(ctx, 5*time.Second); err != nil {
	// 		return err
	// 	}

	// 	op, err = getOpCall.Do()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get operation: %v", err)
	// 	}
	// }
	return waitForOperation(computeService, ctx, inst.ProjectID, inst.Zone, op)
}

// checkOpErrors returns nil if the operation does not have any errors and an
// error summarizing all errors encountered if the operation has errored.
func checkOpErrors(op *compute.Operation) error {
	if op.Error == nil || len(op.Error.Errors) == 0 {
		return nil
	}

	var errs []string
	for _, e := range op.Error.Errors {
		if e.Message != "" {
			errs = append(errs, e.Message)
		} else {
			errs = append(errs, e.Code)
		}
	}
	return errors.New(strings.Join(errs, ","))
}

// func waitForOperation(computeService *compute.Service, ctx context.Context, project, zone string, op *compute.Operation) error {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return fmt.Errorf("timeout waiting for operation to complete")
// 		case <-ticker.C:
// 			result, err := computeService.ZoneOperations.Get(project, zone, op.Name).Do()
// 			if err != nil {
// 				return fmt.Errorf("ZoneOperations.Get: %s", err)
// 			}

// 			if result.Status == "DONE" {
// 				if result.Error != nil {
// 					var errors []string
// 					for _, e := range result.Error.Errors {
// 						errors = append(errors, e.Message)
// 					}
// 					return fmt.Errorf("operation %q failed with error(s): %s", op.Name, strings.Join(errors, ", "))
// 				}
// 				return nil
// 			}
// 		}
// 	}
// }
