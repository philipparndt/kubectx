package kube

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
)

func FileName() string {
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

func Load() *api.Config {
	// Load the kubeconfig file
	config, err := clientcmd.LoadFromFile(FileName())
	if err != nil {
		panic(err)
	}

	return config
}
