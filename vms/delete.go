package vms

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

// DeleteInstance deletes an instance with project id, name, and zone matched
// by inst.
func DeleteInstance(computeService *compute.Service, ctx context.Context, inst *InstanceConfig) error {
	// if op, err := computeService.Instances.Delete(inst.ProjectID, inst.Zone, inst.Name).Context(ctx).Do(); err != nil {
	// 	return fmt.Errorf("Instances.Delete(%s) got error: %v", inst.Name, err)
	// }

	op, err := computeService.Instances.Delete(inst.ProjectID, inst.Zone, inst.Name).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Instances.Delete(%s) got error: %v", inst.Name, err)
	}
	return waitForOperation(computeService, ctx, inst.ProjectID, inst.Zone, op)
}
