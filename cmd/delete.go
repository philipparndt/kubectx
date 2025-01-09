package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a context",
	Long:  `Delete a context. If no context is provided, a list of available contexts will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := kube.LoadDefault()

		contexts := kube.SelectContext(config, args)

		var deleted []string
		for name := range config.Contexts {
			for _, n := range contexts {
				if name == n {
					deleted = append(deleted, name)
					ctx := config.Contexts[name]

					// Check if the cluster and user are still used by other contexts
					clusterUsed := false
					userUsed := false
					for _, otherCtx := range config.Contexts {
						if otherCtx.Cluster == ctx.Cluster && otherCtx != ctx {
							clusterUsed = true
						}
						if otherCtx.AuthInfo == ctx.AuthInfo && otherCtx != ctx {
							userUsed = true
						}
					}

					// Delete context
					delete(config.Contexts, name)

					// Delete cluster if not used
					if !clusterUsed {
						delete(config.Clusters, ctx.Cluster)
					}

					// Delete user if not used
					if !userUsed {
						delete(config.AuthInfos, ctx.AuthInfo)
					}
				}
			}
		}

		if len(deleted) == 0 {
			fmt.Println("No contexts found")
			return
		}

		kube.Backup()
		kube.Save(config)
		fmt.Println("Deleted contexts:", deleted)
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
