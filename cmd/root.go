package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dispatcher",
	Short: "A delivery route dispatcher using KMeans clustering and Euclidean distance sorting",
	Long: `Dispatcher is a CLI tool that takes order locations and assigns them to drivers
in an optimized delivery route using KMeans clustering and Euclidean distance sorting.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
