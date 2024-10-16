/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/philipparndt/kubectx/internal/kube"
	"github.com/spf13/cobra"
	"sort"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available contexts",
	Long:  `List all available contexts`,
	Run: func(cmd *cobra.Command, args []string) {
		config := kube.LoadDefault()

		items := make([]string, 0, len(config.Contexts))

		for name, _ := range config.Contexts {
			items = append(items, name)
		}

		// Sort by label
		sort.Slice(items, func(i, j int) bool {
			return items[i] < items[j]
		})

		for _, name := range items {
			fmt.Println(name)
		}
	},
	Aliases: []string{"l"},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
