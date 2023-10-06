package cmd

import (
	"fmt"

	"github.com/adowair/kvdb/kv"
	"github.com/spf13/cobra"
)

const timeFormat = "2006-01-02 15:04:05"

// tsCmd represents the ts command
var tsCmd = &cobra.Command{
	Use:   "ts",
	Short: "Get the created and last-modified timestamps for a key",
	Long: `Get the times this key was first and last set
These timestamps are durable--moving a database folder will not affect them.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		first, last, err := kv.Timestamps(key)
		if err != nil {
			return err
		}

		fmt.Printf("Key %s:\n", key)
		fmt.Printf("  First set on %s\n", first.Format(timeFormat))
		fmt.Printf("  Last set on  %s\n", last.Format(timeFormat))
		return nil
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
