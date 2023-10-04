package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tsCmd represents the ts command
var tsCmd = &cobra.Command{
	Use:   "ts",
	Short: "Get the created and last-modified timestamps for a key",
	Long: `Get the times this key was first and last set
These timestamps are durable--moving a database folder will not affect them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ts called")
	},
}

func init() {
	rootCmd.AddCommand(tsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
