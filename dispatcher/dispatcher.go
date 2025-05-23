package dispatcher

import (
	"errors"
	"fmt"

	"github.com/lalsaady/dispatcher/client"
	"github.com/lalsaady/dispatcher/model"
)

type Dispatcher struct {
	kmeans   client.KMeansClient
	geocoder client.GeocoderClient
}

func NewDispatcher(km client.KMeansClient, g client.GeocoderClient) (*Dispatcher, error) {
	if km == nil {
		return nil, errors.New("kmeans client is required")
	}
	if g == nil {
		return nil, errors.New("geocoder client is required")
	}
	return &Dispatcher{
		kmeans:   km,
		geocoder: g,
	}, nil
}

func (d *Dispatcher) AssignRoutes(addresses []string, drivers []string) (map[string][]model.Location, error) {
	if len(addresses) == 0 {
		return nil, errors.New("no addresses provided")
	}
	if len(drivers) == 0 {
		return nil, errors.New("no drivers provided")
	}

	// Get coordinates for all orders first
	ordersWithCoords := make([]model.Location, len(addresses))
	for i, address := range addresses {
		coords, err := d.geocoder.GetCoords(address)
		if err != nil {
			return nil, fmt.Errorf("error getting coords from google maps api: %v", err)
		}
		ordersWithCoords[i] = model.Location{
			ID:      i + 1,
			Address: address,
			Lat:     coords.Lat,
			Lon:     coords.Lon,
		}
	}

	// If there are too many orders, increase the ordersPerDriver
	ordersPerDriver := 2
	minOrdersPerDriver := len(addresses) / len(drivers)
	if minOrdersPerDriver > ordersPerDriver {
		ordersPerDriver = minOrdersPerDriver
	}

	// Partition orders into clusters
	clusters, err := d.kmeans.Partition(ordersWithCoords, ordersPerDriver)
	if err != nil {
		return nil, fmt.Errorf("error executing k-means algorithm: %v", err)
	}

	// Assign drivers to clusters
	driverRoutes := make(map[string][]model.Location)
	for i, cluster := range clusters {
		driver := drivers[i]
		locations := make([]model.Location, len(cluster.Observations))

		// Get locations for this cluster
		for j, obs := range cluster.Observations {
			orderObs, ok := obs.(client.Observer)
			if !ok {
				return nil, fmt.Errorf("error getting locations from k-means algorithm")
			}
			for _, order := range ordersWithCoords {
				if order.ID == orderObs.GetID() {
					locations[j] = order
					break
				}
			}
		}

		driverRoutes[driver] = locations
	}

	return driverRoutes, nil
}
