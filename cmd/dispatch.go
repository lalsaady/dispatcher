package cmd

import (
	"fmt"
	"os"

	"github.com/lalsaady/dispatcher/client"
	"github.com/lalsaady/dispatcher/dispatcher"
	"github.com/spf13/cobra"
)

var (
	address []string
	driver  []string
)

var dspCmd = &cobra.Command{
	Use:     "dispatch",
	Aliases: []string{"dsp"},
	Short:   "A delivery route dispatcher using KMeans clustering and Euclidean distance sorting.",
	Long:    `Dispatch orders to drivers using KMeans clustering. Example: dispatch -a "123 Main St,456 Elm St" -d "Alice,Bob"`,
	RunE:    runDispatch,
}

func init() {
	dspCmd.Flags().StringSliceVarP(&address, "address", "a", []string{}, `Address (e.g. -a "123 Main St,456 Random Rd")`)
	dspCmd.Flags().StringSliceVarP(&driver, "driver", "d", []string{}, `List of drivers (e.g. -d "Alice,Bob")`)
	dspCmd.MarkFlagRequired("address")
	dspCmd.MarkFlagRequired("driver")
}

func Execute() {
	if err := dspCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runDispatch(cmd *cobra.Command, args []string) error {
	d, err := dispatcher.NewDispatcher(client.NewKMeansClient(), client.NewGeocoderClient())
	if err != nil {
		return fmt.Errorf("error creating dispatcher: %w", err)
	}
	routes, err := d.AssignRoutes(address, driver)
	if err != nil {
		return fmt.Errorf("error assigning routes: %w", err)
	}
	// Print routes
	fmt.Println("\nAssigned Routes:")
	for driver, route := range routes {
		fmt.Printf("\nDriver %s (%d orders):\n", driver, len(route))
		for i, order := range route {
			fmt.Printf("%d. %s (%.6f, %.6f)\n", i+1, order.Address, order.Lat, order.Lon)
		}
	}
	return nil
}
