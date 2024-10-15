package kube

import (
	"github.com/philipparndt/kubectx/internal/cui"
	"k8s.io/client-go/tools/clientcmd/api"
)

func SelectContext(config *api.Config, args []string) []string {
	names := args
	if len(names) == 0 {
		ctx := cui.SelectContext(config)
		if ctx != nil {
			names = append(names, ctx.Name)
		}
	}

	if len(names) == 0 {
		return []string{}
	}

	return names
}
