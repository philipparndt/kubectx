package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/cui"
	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
	"log"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Switch to a different context",
	Long:  `Switch to a different context. If no context is provided, a list of available contexts will be shown.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := kube.Load()

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

		kube.Save(config)

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
