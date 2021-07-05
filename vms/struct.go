package vms

// InstanceConfig is configuration for starting single GCE instance for
// profiling agent test case.
type InstanceConfig struct {
	ProjectID     string
	Zone          string
	Name          string
	StartupScript string
	MachineType   string
	ImageProject  string
	ImageFamily   string
	Scopes        []string
}
