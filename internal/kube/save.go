package kube

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Save(config *api.Config) {
	// Save changes to kubeconfig file
	err := clientcmd.WriteToFile(*config, FileName())
	if err != nil {
		panic(err)
	}
}
