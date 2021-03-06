package main

import (
	"kwatch/cli"
)

func main() {
	kubewatch := cli.NewClient()
	kubewatch.Pods = kubewatch.GetPods()
	kubewatch.Spec = kubewatch.GetDeploymentSpec()
	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go kubewatch.Run(1, stop)
	go kubewatch.Parse()

	select {}
}
