package vms

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

// https://gitee.com/arohat/google-cloud-go/blob/v0.34.0/profiler/proftest/proftest.go

func waitForOperation(computeService *compute.Service, ctx context.Context, project, zone string, op *compute.Operation) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for operation to complete")
		case <-ticker.C:
			result, err := computeService.ZoneOperations.Get(project, zone, op.Name).Do()
			if err != nil {
				return fmt.Errorf("ZoneOperations.Get: %s", err)
			}

			if result.Status == "DONE" {
				if result.Error != nil {
					var errors []string
					for _, e := range result.Error.Errors {
						errors = append(errors, e.Message)
					}
					return fmt.Errorf("operation %q failed with error(s): %s", op.Name, strings.Join(errors, ", "))
				}
				return nil
			}
		}
	}
}

// PollForSerialOutput polls serial port 1 of the GCE instance specified by
// inst and returns when the finishString appears in the serial output
// of the instance, or when the context times out.
func PollForSerialOutput(computeService *compute.Service, ctx context.Context, inst *InstanceConfig, finishString, errorString string) error {
	var output string
	defer func() {
		log.Printf("Serial port output for %s:\n%s", inst.Name, output)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(20 * time.Second):
			resp, err := computeService.Instances.GetSerialPortOutput(inst.ProjectID, inst.Zone, inst.Name).Port(1).Context(ctx).Do()
			if err != nil {
				// Transient failure.
				log.Printf("Transient error getting serial port output from instance %s (will retry): %v", inst.Name, err)
				continue
			}
			if resp.Contents == "" {
				log.Printf("Ignoring empty serial port output from instance %s (will retry)", inst.Name)
				continue
			}
			if output = resp.Contents; strings.Contains(output, finishString) {
				return nil
			}
			// if strings.Contains(output, errorString) {
			// 	return fmt.Errorf("failed to execute the prober benchmark script")
			// }
		}
	}
}
