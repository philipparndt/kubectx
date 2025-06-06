package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/kube"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectx",
	Short: "Mange your kubernetes contexts",
	Long: `kubectx is a tool to manage your kubernetes contexts.
It allows you to switch between contexts, delete contexts and import new contexts to you local
kubernetes configuration.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result := cmd.Help()

		cfg := kube.LoadDefault()
		fmt.Printf("\n-------\nCurrent context: <%s>\n", cfg.CurrentContext)

		return result
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectx.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
