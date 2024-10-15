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
					fmt.Println("Context", name, "already exists, skipping")
					continue
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
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
