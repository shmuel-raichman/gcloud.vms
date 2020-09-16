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

	// getVMs(computeService, projectID, zone)
	stopVM(computeService, projectID, zone, instanceName)
	// getVMs(computeService, projectID, zone)
	getVM(computeService, projectID, zone, instanceName)
}

func startVM(computeService *compute.Service, instanceName string, projectID string, zone string) {





	// ########################################################################################################
	// ########################################################################################################
	// ########################################################################################################
	resStart, err := computeService.Instances.Start(projectID, zone, instanceName).Do()
	// fmt.Println(res)
	fmt.Println("")
	etagStart := resStart.Header.Get("Etag")
	// log.Printf("Etag=%v", etag)

	fmt.Println("")
	instStart, err := computeService.Instances.Get(projectID, zone, instanceName).IfNoneMatch(etagStart).Do()
	// log.Printf("Got compute.Instance, err: %#v, %v", instStart.Name, err)
	log.Printf("Got compute.Instance, err: %#v, %v", instStart.Name, err, instStart.Status)
	// ########################################################################################################
	// ########################################################################################################
	// ########################################################################################################
}


func stopVM(computeService *compute.Service, instanceName string, projectID string, zone string) {


	res, err := computeService.Instances.Stop(instanceName, projectID, zone).Do()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Stop vm res", res.Status)
}

func getVMs(computeService *compute.Service, projectID string, zone string){
		// ########################################################################################################
	// ########################################################################################################
	// ########################################################################################################
	res, _ := computeService.Instances.List(projectID, zone).Do()
	// fmt.Println(res.Items)
	// fmt.Println("")
	// etag := res.Header.Get("Etag")
	// log.Printf("Etag=%v", etag)

	// fmt.Println("")
	// inst, err := computeService.Instances.Get(projectID, zone, instanceName).IfNoneMatch(etag).Do()
	// log.Printf("Got compute.Instance, err: %#v, %v", inst.Name, err)
	// ########################################################################################################
	// ########################################################################################################
	// ########################################################################################################
	// res, err := computeService.Instances.Start(projectID, zone, instanceName).Do()
	// fmt.Println(res)
	// fmt.Println("")
	// etag := res.Header.Get("Etag")
	// log.Printf("Etag=%v", etag)

	// fmt.Println("")
	// inst, err := computeService.Instances.Get(projectID, zone, instanceName).IfNoneMatch(etag).Do()
	// log.Printf("Got compute.Instance, err: %#v, %v", inst.Name, err)
	// ########################################################################################################
	// ########################################################################################################
	// ########################################################################################################
	// bodyBytes, err := ioutil.ReadAll(res.Items)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // bodyString := string(bodyBytes)
	// log.Println(bodyString)
	// fmt.Printf("%v", res.Items)
	for _, vm := range res.Items {
		fmt.Print("VM Name is: ", vm.Name)
		// fmt.Printf("%+v\n", vm)
		fmt.Println(" - VM state is: ", vm.Status)
	}
}


func getVM(computeService *compute.Service, projectID string, zone string, instanceName string){
	res, _ := computeService.Instances.Get(projectID, zone, instanceName).Do()
	fmt.Println(res.Status)

}

// // 	fmt.Println(computeService)
// // }

// import (
// 	"context"
// 	"fmt"

// 	"golang.org/x/oauth2/google"
// 	compute "google.golang.org/api/compute/v1"
// )

// func main() {
// 	// Use oauth2.NoContext if there isn't a good context to pass in.
// 	ctx := context.Background()

// 	client, err := google.DefaultClient(ctx, compute.ComputeScope)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	computeService, err := compute.New(client, option.WithCredentialsFile(jsonPath))
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(computeService)
// }
