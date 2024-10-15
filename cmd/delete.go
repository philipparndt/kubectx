package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/cui"
	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a context",
	Long:  `Delete a context. If no context is provided, a list of available contexts will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := kube.Load()

		names := args
		if len(names) == 0 {
			ctx := cui.SelectContext(config)
			if ctx != nil {
				names = append(names, ctx.Name)
			}
		}

		if len(names) != 0 {
			deleted := []string{}
			for name := range config.Contexts {
				for _, n := range names {
					if name == n {
						deleted = append(deleted, name)
						ctx := config.Contexts[name]

						// delete context, cluster and user
						delete(config.Contexts, name)
						delete(config.Clusters, ctx.Cluster)
						delete(config.AuthInfos, ctx.AuthInfo)
					}
				}
			}

			if len(deleted) == 0 {
				fmt.Println("No contexts found")
				return
			}

			kube.Save(config)
			fmt.Println("Deleted contexts:", deleted)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
