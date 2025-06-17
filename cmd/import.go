package cmd

import (
	"fmt"
	"os"

	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
)

func validateFiles(files []string) ([]string, error) {
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			return nil, fmt.Errorf("cannot import directory %s", file)
		}
	}

	return files, nil
}

var upgrade = false

// importCmd represents the add command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import contexts from a file to the default kubeconfig",
	Long:  `Import contexts from a file to the default kubeconfig.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No file to import specified")
			return
		}

		files, err := validateFiles(args)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		config := kube.LoadDefault()
		for _, file := range files {
			fmt.Println("Importing", file)

			toBeImported := kube.Load(file)

			// Upgrade keeps the existing names
			if upgrade {
				// Build a map from server to existing cluster name
				serverToName := make(map[string]string)
				for name, cluster := range config.Clusters {
					serverToName[cluster.Server] = name
				}

				// For each imported cluster, if a cluster with the same server exists, rename it in the imported config
				for importedName, importedCluster := range toBeImported.Clusters {
					if existingName, ok := serverToName[importedCluster.Server]; ok && importedName != existingName {
						fmt.Println("Importing cluster", importedName, "as", existingName, "(matched by server)")
						// Update all contexts to use the existing cluster name
						for ctxName, ctx := range toBeImported.Contexts {
							if ctx.Cluster == importedName {
								ctx.Cluster = existingName
								toBeImported.Contexts[ctxName] = ctx
							}
						}
						// If the existing cluster name is not already in toBeImported, add/overwrite it
						toBeImported.Clusters[existingName] = importedCluster
						// Remove the old cluster name from the import set
						delete(toBeImported.Clusters, importedName)
					}
				}
			}

			for name, ctx := range toBeImported.Contexts {
				if _, ok := config.Contexts[name]; ok {
					if upgrade {
						fmt.Println("Upgrading context", name)
					} else {
						fmt.Println("Skipping context", name, "as it already exists (use --upgrade to overwrite)")
						continue
					}
				}

				config.Contexts[name] = ctx
			}
			for name, authInfo := range toBeImported.AuthInfos {
				config.AuthInfos[name] = authInfo
			}
		}

		kube.Backup()
		kube.Save(config)
	},
	Aliases: []string{"add"},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Upgrade existing contexts")
}
