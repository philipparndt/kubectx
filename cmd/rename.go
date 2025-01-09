package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/cui"
	"github.com/philipparndt/kubectx/internal/kube"
	"log"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a context",
	Long:  `Rename a context. If no context is provided, a list of available contexts will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := kube.LoadDefault()
		contexts := kube.SelectContext(config, args)
		if len(contexts) == 0 {
			fmt.Println("No context selected")
			return
		} else if len(contexts) > 1 {
			log.Panic("Too many arguments")
		}

		newName := cui.RenameForm(contexts[0])
		if newName == contexts[0] {
			fmt.Println("No changes made")
			return
		}

		ctx := config.Contexts[contexts[0]]
		delete(config.Contexts, contexts[0])
		config.Contexts[newName] = ctx

		kube.Backup()
		kube.Save(config)

		fmt.Println("Context renamed to:", newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
