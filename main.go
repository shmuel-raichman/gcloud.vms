// B"H
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main(){
	argv := os.Args
	if len(argv) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: action must be supplied - start, stop, status")
		return
	}

	action := argv[1]
	
	projectID := "shmulik-general-dev"
	instanceName := "ester-wp"
	zone := "us-central1-a"

	if len(argv) == 3 {
		instanceName = argv[2]
	}

	fmt.Println(instanceName)

	fmt.Println("Starting my app...")
	jsonPath := "./resources/shmulik-general-dev.json"
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "start":
		log.Println("action: ", action)
		startVM(computeService, projectID, zone, instanceName)
		getVMStatus(computeService, projectID, zone, instanceName)
	case "stop":
		log.Println("action: ", action)
		stopVM(computeService, projectID, zone, instanceName)
		getVMStatus(computeService, projectID, zone, instanceName)
	case "status":
		log.Println("action: ", action)
		getVMStatus(computeService, projectID, zone, instanceName)
		getVMs(computeService, projectID, zone)
	default:
		log.Println("action: ", action)
		log.Println("Action must be supplied - start, stop, status")
	}

}

func startVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Start(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Starting vm: ", instanceName)
	fmt.Println("Start status: ", res.Status)
}


func stopVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Stop(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Starting vm: ", instanceName)
	fmt.Println("Stop status: ", res.Status)
	// ctx := context.Background()
	// // getVMStatus(computeService, projectID, zone, instanceName)
	// insta, err := computeService.Instances.Get(projectID, zone, instanceName).IfNoneMatch(etag).Context(ctx).Do()
	// // res, err := computeService.Instances.Get(projectID, instanceName, zone).Do()
	// // res, err := computeService.Instances.Get(zone, projectID, instanceName).Do()
	// // res, err := computeService.Instances.Get(zone, instanceName, projectID).Do()
	// // res, err := computeService.Instances.Get(instanceName, zone, projectID).Do()
	// // res, err := computeService.Instances.Get(instanceName, projectID, zone).Do()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(insta.Status)
}

func getVMs(computeService *compute.Service, projectID string, zone string){

	res, _ := computeService.Instances.List(projectID, zone).Do()

	for _, vm := range res.Items {
		fmt.Print("VM Name is: ", vm.Name)
		// fmt.Printf("%+v\n", vm)
		fmt.Println(" - VM state is: ", vm.Status)
	}
}


func getVMStatus(computeService *compute.Service, projectID string, zone string, instanceName string){
	res, err := computeService.Instances.Get(projectID, zone, instanceName).Do()
	// res, err := computeService.Instances.Get(projectID, instanceName, zone).Do()
	// res, err := computeService.Instances.Get(zone, projectID, instanceName).Do()
	// res, err := computeService.Instances.Get(zone, instanceName, projectID).Do()
	// res, err := computeService.Instances.Get(instanceName, zone, projectID).Do()
	// res, err := computeService.Instances.Get(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(res)
	fmt.Println(res.Status)
}
