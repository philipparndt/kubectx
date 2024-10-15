/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/cui"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find kubeconfig file (default is $HOME/.kube/config)
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		// Load the kubeconfig file
		config, err := clientcmd.LoadFromFile(kubeconfig)
		if err != nil {
			panic(err)
		}

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

			// Save changes to kubeconfig file
			err = clientcmd.WriteToFile(*config, "./tmp/config")
			if err != nil {
				log.Panic(err)
			}

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
