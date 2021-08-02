// B"H
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	vms "smuel1414/gcloud.vms/vms"
)

var (
	script string = `

	export GITHUB_PASSWORD=` + os.Getenv("GITHUB_PASSWORD") + `
	export GITHUB_USERNAME=` + os.Getenv("GITHUB_USERNAME") + `
	export GITHUB_USERNAME=` + os.Getenv("GITHUB_INITAL_REPO") + `
	export VM_GCLOUD_USER= ` + os.Getenv("VM_GCLOUD_USER") + `
	export VM_SSH_USER= ` + os.Getenv("VM_SSH_USER") + `
	export DOCKER_COMPOSE_VERSION= ` + os.Getenv("DOCKER_COMPOSE_VERSION") + `

	sudo apt-get install git tree vim -y
	cd /opt
	mkdir init
	cd init
	git clone https://$GITHUB_USERNAME:$GITHUB_PASSWORD@github.com/$GITHUB_USERNAME/$GITHUB_INITAL_REPO

	cd $GITHUB_INITAL_REPO
	chmod +x installdocker.sh
	./installdocker.sh

	sudo usermod -aG docker $VM_GCLOUD_USER
	sudo usermod -aG docker $VM_SSH_USER

	echo "DONE INITIALIZING STARTUP SCRIPT"
	`
)

func main() {

	var usage string = "Usage: action must be supplied from - start, stop, status, create, delete"
	argv := os.Args
	if len(argv) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	var action string
	var instanceName string = "my-vm"

	switch len(argv) {
	case 1:
		log.Fatal(usage)
	case 2:
		action = argv[1]
	case 3:
		action = argv[1]
		instanceName = argv[2]
	default:
		println(usage)
	}

	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	jsonPath := os.Getenv("GCLOUD_SERVICE_ACCOUNT_JSON_PATH")
	zone := "us-central1-a"

	fmt.Println(instanceName)

	fmt.Println("Starting my app...")
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}

	scopesForInst := []string{
		"https://www.googleapis.com/auth/devstorage.read_only",
		"https://www.googleapis.com/auth/logging.write",
		"https://www.googleapis.com/auth/monitoring.write",
		"https://www.googleapis.com/auth/servicecontrol",
		"https://www.googleapis.com/auth/service.management.readonly",
		"https://www.googleapis.com/auth/trace.append",
	}

	instanceConfig := vms.InstanceConfig{
		ProjectID:     projectID,
		Zone:          zone,
		Name:          instanceName,
		StartupScript: script,
		MachineType:   "g1-small",
		ImageProject:  "debian-cloud",
		ImageFamily:   "debian-10",
		Scopes:        scopesForInst,
	}

	switch action {
	case "start":
		log.Println("action: ", action)
		startVM(computeService, projectID, zone, instanceName)
		if err := getVMStatus(computeService, projectID, zone, instanceConfig.Name); err != nil {
			log.Fatal(err)
		}
	case "stop":
		log.Println("action: ", action)
		stopVM(computeService, projectID, zone, instanceName)
		if err := getVMStatus(computeService, projectID, zone, instanceConfig.Name); err != nil {
			log.Println(err)
		}
	case "status-all":
		log.Println("action: ", action)
		if err := getVMs(computeService, projectID, zone); err != nil {
			log.Println(err)
		}
	case "status":
		log.Println("action: ", action)
		if err := getVMStatus(computeService, projectID, zone, instanceConfig.Name); err != nil {
			log.Println(err)
		}
	case "create":
		log.Println("action: ", action)
		vms.CreateInstance(computeService, ctx, &instanceConfig)
		err := vms.PollForSerialOutput(computeService, ctx, &instanceConfig, "DONE INITIALIZING STARTUP SCRIPT", "error is now")
		if err != nil {
			log.Println(err)
		}
		if err := getVMStatus(computeService, projectID, zone, instanceConfig.Name); err != nil {
			log.Fatal(err)
		}
	case "delete":
		log.Println("action: ", action)
		vms.DeleteInstance(computeService, ctx, &instanceConfig)
		if err := getVMStatus(computeService, projectID, zone, instanceConfig.Name); err != nil {
			log.Println(err)
		}

	default:
		log.Println("action: ", action)
		log.Println(usage)
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
}

func getVMs(computeService *compute.Service, projectID string, zone string) error {

	list, err := computeService.Instances.List(projectID, zone).Do()
	if err != nil {
		return err
	}

	// listpretty, err := json.MarshalIndent(list, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	//Marshal

	if len(list.Items) == 0 {
		fmt.Printf("No VMs in project: %s\n", projectID)
		return nil
	}

	for _, vm := range list.Items {
		fmt.Print("VM Name is: ", vm.Name)
		// fmt.Printf("%+v\n", vm)
		fmt.Println(" - VM state is: ", vm.Status)
	}
	return nil
}

func getVMStatus(computeService *compute.Service, projectID string, zone string, instanceName string) error {
	res, err := computeService.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return err
	}
	fmt.Println(instanceName, "IS: ", res.Status)
	return nil
}
