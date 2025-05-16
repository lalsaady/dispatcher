package main

import (
	"errors"
	"fmt"
	"sort"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

type KMeans interface {
	Partition(observations clusters.Observations, k int) (clusters.Clusters, error)
}

type Dispatcher struct {
	kmeans KMeans
}

func NewDispatcher(km KMeans) *Dispatcher {
	if km == nil {
		km = kmeans.New()
	}
	return &Dispatcher{kmeans: km}
}

// Location represents a delivery location
type Location struct {
	ID      int
	Address string
	Lat     float64
	Lon     float64
}

// OrderObservation links an Order ID with coordinates
type OrderObservation struct {
	ID     int
	coords clusters.Coordinates
}

func (o OrderObservation) Coordinates() clusters.Coordinates {
	return o.coords
}

func (o OrderObservation) Distance(c clusters.Coordinates) float64 {
	return o.coords.Distance(c)
}

var hub = Location{
	Address: "2800 Euclid Ave, Cleveland, OH",
	Lat:     41.502069,
	Lon:     -81.669011,
}

func distance(lat2, lon2 float64) float64 {
	return (hub.Lat-lat2)*(hub.Lat-lat2) + (hub.Lon-lon2)*(hub.Lon-lon2)
}

func (d *Dispatcher) AssignRoutes(orders []Location, drivers []string) (map[string][]Location, error) {
	if len(orders) == 0 {
		return nil, errors.New("no orders provided")
	}
	if len(drivers) == 0 {
		return nil, errors.New("no drivers provided")
	}

	var observations clusters.Observations
	idToIndex := make(map[int]int)

	// Prepare observations with order IDs
	for i, order := range orders {
		ob := OrderObservation{
			ID:     order.ID,
			coords: clusters.Coordinates{order.Lat, order.Lon},
		}
		observations = append(observations, ob)
		idToIndex[order.ID] = i
	}

	// If there are too many orders, increase the ordersPerDriver
	ordersPerDriver := 2
	totalOrders := len(orders)
	minOrdersPerDriver := totalOrders / len(drivers)
	if minOrdersPerDriver > ordersPerDriver {
		ordersPerDriver = minOrdersPerDriver
	}

	// Perform k-means clustering
	clustered, err := d.kmeans.Partition(observations, ordersPerDriver)
	if err != nil {
		return nil, errors.New("error executing k-means algorithm")
	}

	// Assign drivers and order deliveries by distance from hub
	driverRoutes := make(map[string][]Location)
	for i, cluster := range clustered {
		driver := drivers[i]
		locations := make([]Location, len(cluster.Observations))

		// Get locations for this cluster
		for j, obs := range cluster.Observations {
			orderObs := obs.(OrderObservation)
			idx := idToIndex[orderObs.ID]
			locations[j] = orders[idx]
		}

		// Sort locations by distance from hub
		sort.Slice(locations, func(i, j int) bool {
			distI := distance(locations[i].Lat, locations[i].Lon)
			distJ := distance(locations[j].Lat, locations[j].Lon)
			return distI < distJ
		})

		// Limit to ordersPerDriver unless we have too many orders
		if len(locations) > ordersPerDriver {
			locations = locations[:ordersPerDriver]
		}

		driverRoutes[driver] = locations
	}

	fmt.Printf("Driver routes: %v\n", driverRoutes)
	return driverRoutes, nil
}

func main() {
	dispatcher := NewDispatcher(nil)
	orders := []Location{
		{ID: 1, Address: "123 Main St, Cleveland, OH", Lat: 41.498612, Lon: -81.694471},
		{ID: 2, Address: "456 Dummy Ave, Cleveland, OH", Lat: 41.478727, Lon: -81.738038},
		{ID: 3, Address: "789 Random Rd, Cleveland, OH", Lat: 41.506967, Lon: -81.599513},
		{ID: 4, Address: "1010 Idk St, Cleveland, OH", Lat: 41.477112, Lon: -81.649591},
	}
	drivers := []string{"Alice", "Bob"}
	dispatcher.AssignRoutes(orders, drivers)
}
