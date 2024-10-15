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

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
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

		toBeSelected := ""
		if len(args) > 1 {
			log.Panic("Too many arguments")
		} else if len(args) == 1 {
			toBeSelected = args[0]
		} else {
			ctx := cui.SelectContext(config)
			if ctx == nil {
				log.Panic("No context selected")
			}
			toBeSelected = ctx.Name
		}

		fmt.Println("Selected context:", toBeSelected)
		config.CurrentContext = toBeSelected

		// Save changes to kubeconfig file
		err = clientcmd.WriteToFile(*config, "./tmp/config")
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Kubeconfig updated with new context!")

	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
