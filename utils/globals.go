package utils

import "os"

var (
	Script string = `

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

var ScopesForInst = []string{
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/logging.write",
	"https://www.googleapis.com/auth/monitoring.write",
	"https://www.googleapis.com/auth/servicecontrol",
	"https://www.googleapis.com/auth/service.management.readonly",
	"https://www.googleapis.com/auth/trace.append",
}
