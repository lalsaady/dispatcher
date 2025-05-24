package client

import (
	"sort"

	"github.com/lalsaady/dispatcher/model"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// KMeansClient defines the interface for kmeans clustering
type KMeansClient interface {
	Partition(orders []model.Location, numClusters int, hub model.Location) (clusters.Clusters, error)
}

type KMeans struct{}

func NewKMeansClient() KMeansClient {
	return &KMeans{}
}

// Observer defines the interface for an order
type Observer interface {
	clusters.Observation
	GetID() int
}

// Observation implements Observer with coordinates and ID
type Observation struct {
	ID     int
	coords clusters.Coordinates
}

func (o Observation) GetID() int {
	return o.ID
}

func (o Observation) Coordinates() clusters.Coordinates {
	return o.coords
}

func (o Observation) Distance(c clusters.Coordinates) float64 {
	return o.coords.Distance(c)
}

// NewObservation creates a new Observation
func NewObservation(id int, lat, lon float64) Observation {
	return Observation{
		ID:     id,
		coords: clusters.Coordinates{lat, lon},
	}
}

// Partition groups orders into clusters using k-means clustering
func (k *KMeans) Partition(orders []model.Location, numClusters int, hub model.Location) (clusters.Clusters, error) {
	var observations clusters.Observations

	// Prepare observations with order IDs
	for _, order := range orders {
		ob := NewObservation(order.ID, order.Lat, order.Lon)
		observations = append(observations, ob)
	}

	// Perform k-means clustering
	km := kmeans.New()
	clusters, err := km.Partition(observations, numClusters)
	if err != nil {
		return nil, err
	}

	// Sort each cluster by distance from hub
	for _, cluster := range clusters {
		sortClusterByDistance(cluster, hub)
	}

	return clusters, nil
}

// Helper functions
func sortClusterByDistance(cluster clusters.Cluster, hub model.Location) {
	sort.Slice(cluster.Observations, func(i, j int) bool {
		obsI := cluster.Observations[i].(Observation)
		obsJ := cluster.Observations[j].(Observation)
		coordsI := obsI.Coordinates()
		coordsJ := obsJ.Coordinates()
		distI := euclideanDistance(coordsI[0], coordsI[1], hub)
		distJ := euclideanDistance(coordsJ[0], coordsJ[1], hub)
		return distI < distJ
	})
}

func euclideanDistance(lat2, lon2 float64, hub model.Location) float64 {
	return (hub.Lat-lat2)*(hub.Lat-lat2) + (hub.Lon-lon2)*(hub.Lon-lon2)
}
