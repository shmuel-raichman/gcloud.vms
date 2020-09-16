// B"H
package main

import (
	"context"
	"fmt"
	"log"
	// "io/ioutil"

	compute "google.golang.org/api/compute/v1"
	// "google.golang.org/api/googleapi"
	// "google.golang.org/api/internal/gensupport"
	"google.golang.org/api/option"
	// "google.golang.org/api/option/internaloption"
	// "google.golang.org/api/transport/http"
)

func main(){
	projectID := "shmulik-general-dev"
	instanceName := "ester-wp"
	zone := "us-central1-a"
	// prefix := "https://www.googleapis.com/compute/v1/projects/" + projectID

	fmt.Println("Starting my app...")
	jsonPath := "./resources/shmulik-general-dev.json"
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}

	startVM(computeService, projectID, zone, instanceName)
	// getVMs(computeService, projectID, zone)
	// stopVM(computeService, projectID, zone, instanceName)
	// getVMs(computeService, projectID, zone)
	// getVMStatus(computeService, projectID, zone, instanceName)
}

func startVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Start(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(res)
	fmt.Println("Starting vm: ", instanceName)
	fmt.Println("Start status: ", res.Status)
	// etag := res.Header.Get("etag")

	// fmt.Println("")
	// inst, err := computeService.Instances.Get(instanceName, projectID, zone).Do()
	// inst, err := computeService.Instances.Get(instanceName, projectID, zone).IfNoneMatch(etag).Do()
	// if err != nil {
	// 	log.Println(err)
	// }
	getVMStatus(computeService, projectID, zone, instanceName)
	// inst, err := computeService.Instances.Get(instanceName, projectID, zone).IfNoneMatch(etag).Do()

	// log.Printf("Got compute.Instance, err: %#v, %v", inst.Name, inst.Status)
}


func stopVM(computeService *compute.Service, instanceName string, projectID string, zone string) {
	res, err := computeService.Instances.Stop(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Stop vm res", res.Status)
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
	res, _ := computeService.Instances.Get(projectID, zone, instanceName).Do()
	fmt.Println(res.Status)

}
