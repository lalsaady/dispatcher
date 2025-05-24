package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lalsaady/dispatcher/dispatcher"
	"github.com/spf13/cobra"
)

type Order struct {
	Address string
	Lat     float64
	Lon     float64
}

var (
	orders  []string
	drivers []string
)

var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch orders to drivers",
	Long: `Dispatch orders to drivers using KMeans clustering and Euclidean distance sorting.
Example:
  dispatcher dispatch --order "123 Main St,41.498612,-81.694471" --order "456 Elm St,41.499,-81.695" --drivers Alice,Bob`,
	RunE: runDispatch,
}

func init() {
	rootCmd.AddCommand(dispatchCmd)
	dispatchCmd.Flags().StringArrayVarP(&orders, "order", "o", []string{}, "Order in format 'address,lat,lon' (can be used multiple times)")
	dispatchCmd.Flags().StringSliceVarP(&drivers, "drivers", "d", []string{}, "Comma-separated list of driver names")
	dispatchCmd.MarkFlagRequired("order")
	dispatchCmd.MarkFlagRequired("drivers")
}

func parseOrder(orderStr string) (Order, error) {
	parts := strings.Split(orderStr, ",")
	if len(parts) != 3 {
		return Order{}, fmt.Errorf("invalid order format, expected 'address,lat,lon' but got: %s", orderStr)
	}

	lat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Order{}, fmt.Errorf("invalid latitude: %s", parts[1])
	}

	lon, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Order{}, fmt.Errorf("invalid longitude: %s", parts[2])
	}

	return Order{
		Address: parts[0],
		Lat:     lat,
		Lon:     lon,
	}, nil
}

func runDispatch(cmd *cobra.Command, args []string) error {
	// Parse all orders
	locations := make([]dispatcher.Location, len(orders))
	for i, orderStr := range orders {
		order, err := parseOrder(orderStr)
		if err != nil {
			return fmt.Errorf("error parsing order %d: %w", i+1, err)
		}
		locations[i] = dispatcher.Location{
			ID:      i + 1,
			Address: order.Address,
			Lat:     order.Lat,
			Lon:     order.Lon,
		}
	}

	d := dispatcher.NewDispatcher(nil)
	routes, err := d.AssignRoutes(locations, drivers)
	if err != nil {
		return fmt.Errorf("error assigning routes: %w", err)
	}

	// Print routes in a more readable format
	fmt.Println("\nAssigned Routes:")
	for driver, route := range routes {
		fmt.Printf("\nDriver %s (%d orders):\n", driver, len(route))
		for i, order := range route {
			fmt.Printf("%d. %s (%.6f, %.6f)\n", i+1, order.Address, order.Lat, order.Lon)
		}
	}
	return nil
}
