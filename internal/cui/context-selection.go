package cui

import (
	"k8s.io/client-go/tools/clientcmd/api"
)

type Ctx struct {
	Name    string
	Context *api.Context
}

type CtxLabelProvider struct {
}

func (c CtxLabelProvider) Label(ctx Ctx) string {
	return ctx.Name
}

func SelectContext(config *api.Config) *Ctx {
	ctxs := make([]Ctx, 0, len(config.Contexts))
	for name, ctx := range config.Contexts {
		ctxs = append(ctxs, Ctx{Name: name, Context: ctx})
	}

	return Select(ctxs, CtxLabelProvider{})
}
