package cmd

import (
	"fmt"
	"time"

	"github.com/adowair/kvdb/kv"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the value for a key",
	Long: `Set the value of a key in the database.
If the key does not exist, it is created.
If the key already exists, its old value is overwritten.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, val := args[0], args[1]
		if err := kv.Set(key, val, time.Now()); err != nil {
			return err
		}

		fmt.Printf("Set %s <= %s\n", key, val)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
