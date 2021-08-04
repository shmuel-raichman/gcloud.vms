// B"H
package vms

import (
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

func StartVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Start(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Starting vm: ", instanceName)
	fmt.Println("Start status: ", res.Status)
}

func StopVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Stop(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Starting vm: ", instanceName)
	fmt.Println("Stop status: ", res.Status)
}

func GetVMs(computeService *compute.Service, projectID string, zone string) ([]*compute.Instance, error) {

	list, err := computeService.Instances.List(projectID, zone).Do()
	if err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		fmt.Printf("No VMs in project: %s\n", projectID)
		return nil, nil
	}

	// status := string[]

	for _, vm := range list.Items {
		fmt.Print("VM Name is: ", vm.Name)
		// fmt.Printf("%+v\n", vm)
		fmt.Println(" - VM state is: ", vm.Status)
	}

	return list.Items, nil
}

func GetVMStatus(computeService *compute.Service, projectID string, zone string, instanceName string) (string, error) {
	res, err := computeService.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return "", err
	}
	fmt.Println(instanceName, "IS: ", res.Status)
	return res.Status, nil
}
