package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
	"os"
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
			for name, cluster := range toBeImported.Clusters {
				config.Clusters[name] = cluster
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
