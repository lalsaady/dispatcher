package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/lalsaady/dispatcher/client"
	"github.com/lalsaady/dispatcher/dispatcher"
	"github.com/lalsaady/dispatcher/model"
	"github.com/spf13/cobra"
)

var (
	address []string
	driver  []string
	hub     model.Location
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A delivery route dispatcher using KMeans clustering and Euclidean distance sorting.",
	Long:  `Dispatch orders to drivers using KMeans clustering. Example: dispatch -a "123 Main St,456 Elm St" -d "Alice,Bob" or dispatch -a addresses.csv -d drivers.csv`,
	RunE:  runDispatch,
}

func init() {
	runCmd.Flags().StringSliceVarP(&address, "address", "a", []string{}, `Addresses or CSV file (e.g. -a "123 Main St,456 Elm St" or -a addresses.csv)`)
	runCmd.Flags().StringSliceVarP(&driver, "driver", "d", []string{}, `List of drivers or CSV file (e.g. -d "Alice,Bob" or -d drivers.csv)`)
	runCmd.Flags().StringVarP(&hub.Address, "hub", "u", "", `Hub address (e.g. -u "123 Main St")`)
	runCmd.MarkFlagRequired("address")
	runCmd.MarkFlagRequired("driver")
	runCmd.MarkFlagRequired("hub")
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runDispatch(cmd *cobra.Command, args []string) error {
	if len(address) > 0 && strings.Contains(address[0], ".csv") {
		var err error
		address, err = parseCSV(address)
		if err != nil {
			return fmt.Errorf("error parsing addresses: %w", err)
		}
	}
	if len(driver) > 0 && strings.Contains(driver[0], ".csv") {
		var err error
		driver, err = parseCSV(driver)
		if err != nil {
			return fmt.Errorf("error parsing drivers: %w", err)
		}
	}
	coords, err := client.NewGeocoderClient().GetCoords(hub.Address)
	if err != nil {
		return fmt.Errorf("error getting coordinates for hub: %w", err)
	}
	hub.Lat, hub.Lon = coords.Lat, coords.Lon
	d, err := dispatcher.NewDispatcher(client.NewKMeansClient(), client.NewGeocoderClient(), hub)
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
		fmt.Printf("\nDriver %s:\n", driver)
		for i, order := range route {
			fmt.Printf("%d. %s (%.6f, %.6f)\n", i+1, order.Address, order.Lat, order.Lon)
		}
	}
	return nil
}

// Parses csv
func parseCSV(input []string) ([]string, error) {
	var results []string
	for _, item := range input {
		file, err := os.Open(item)
		if err != nil {
			return nil, fmt.Errorf("failed to open CSV file %s: %w", item, err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV file %s: %w", item, err)
		}
		for _, record := range records {
			results = append(results, record...)
		}
	}
	return results, nil
}
