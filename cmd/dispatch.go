package cmd

import (
	"fmt"

	"github.com/lalsaady/dispatcher/client"
	"github.com/lalsaady/dispatcher/dispatcher"
	"github.com/spf13/cobra"
)

var (
	address []string
	driver  []string
)

var dispatchCmd = &cobra.Command{
	Use:     "dispatch",
	Aliases: []string{"dsp"},
	Short:   "Dispatch orders to drivers",
	Long:    `Dispatch orders to drivers using KMeans clustering. Example: dsp -a "123 Main St,456 Elm St" -d "Alice,Bob"`,
	RunE:    runDispatch,
}

func init() {
	rootCmd.AddCommand(dispatchCmd)
	dispatchCmd.Flags().StringSliceVarP(&address, "address", "a", []string{}, `Address (e.g. -a "123 Main St,456 Random Rd")`)
	dispatchCmd.Flags().StringSliceVarP(&driver, "driver", "d", []string{}, `List of drivers (e.g. -d "Alice,Bob")`)
	dispatchCmd.MarkFlagRequired("address")
	dispatchCmd.MarkFlagRequired("driver")
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
