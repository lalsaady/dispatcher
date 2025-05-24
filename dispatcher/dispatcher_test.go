package dispatcher

import (
	"errors"
	"testing"

	"github.com/lalsaady/dispatcher/client"
	mocks "github.com/lalsaady/dispatcher/mock"
	"github.com/lalsaady/dispatcher/model"
	"github.com/muesli/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testObservation is a test implementation of client.OrderObserver
type testObservation struct {
	id          int
	coordinates clusters.Coordinates
}

func (o testObservation) GetID() int {
	return o.id
}

func (o testObservation) Coordinates() clusters.Coordinates {
	return o.coordinates
}

func (o testObservation) Distance(c clusters.Coordinates) float64 {
	return o.coordinates.Distance(c)
}

func TestDispatcherInvalid(t *testing.T) {
	_, err := NewDispatcher(client.NewKMeansClient(), client.NewGeocoderClient(), model.Location{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "hub address and coords are required")
	dispatcher, err := NewDispatcher(client.NewKMeansClient(), client.NewGeocoderClient(), model.Location{Address: "123 Main St, Cleveland, OH", Lat: 41.503322, Lon: -81.698311})
	assert.NoError(t, err)
	_, err = dispatcher.AssignRoutes([]string{}, []string{"Alice", "Bob"})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no addresses provided")
	routes, err := dispatcher.AssignRoutes([]string{
		"123 Main St, Cleveland, OH",
	}, []string{})
	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, err.Error(), "no drivers provided")
}

func TestDispatcherValid(t *testing.T) {
	addresses := []string{
		"123 Main St, Cleveland, OH",
		"345 Dummy Ave, Cleveland, OH",
	}
	drivers := []string{"Alice", "Bob"}

	// Setup mocks
	mockKM := mocks.NewMockKMeans()
	mockGeo := mocks.NewMockGeocoder()
	mockClusters := clusters.Clusters{
		{Observations: []clusters.Observation{
			testObservation{
				id:          1,
				coordinates: clusters.Coordinates{40.988612, -80.698871},
			},
		}},
		{Observations: []clusters.Observation{
			testObservation{
				id:          2,
				coordinates: clusters.Coordinates{41.073612, -80.609771},
			},
		}},
	}
	mockGeo.On("GetCoords", mock.Anything).Return(model.Points{Lat: 40.988612, Lon: -80.698871}, nil)
	mockKM.On("Partition", mock.Anything, mock.Anything, mock.Anything).Return(mockClusters, nil)

	dispatcher, err := NewDispatcher(mockKM, mockGeo, model.Location{Address: "123 Main St, Cleveland, OH", Lat: 41.503322, Lon: -81.698311})
	assert.NoError(t, err)
	routes, err := dispatcher.AssignRoutes(addresses, drivers)
	assert.NoError(t, err)
	assert.NotNil(t, routes)
	assert.Equal(t, 2, len(routes))
	assert.Equal(t, 1, len(routes["Alice"]))
	assert.Equal(t, 1, len(routes["Bob"]))
	mockKM.AssertExpectations(t)
	mockGeo.AssertExpectations(t)
}

func TestDispatcherError(t *testing.T) {
	addresses := []string{
		"123 Main St, Cleveland, OH",
		"345 Dummy Ave, Cleveland, OH",
	}
	drivers := []string{"Alice", "Bob"}

	mockKM := mocks.NewMockKMeans()
	mockGeo := mocks.NewMockGeocoder()
	// Setup geocoder to succeed
	mockGeo.On("GetCoords", mock.Anything).Return(model.Points{Lat: 40.988612, Lon: -80.698871}, nil)
	// Setup kmeans to fail
	mockKM.On("Partition", mock.Anything, mock.Anything, mock.Anything).Return(clusters.Clusters{}, errors.New("mock error"))

	dispatcher, err := NewDispatcher(mockKM, mockGeo, model.Location{Address: "123 Main St, Cleveland, OH", Lat: 41.503322, Lon: -81.698311})
	assert.NoError(t, err)
	routes, err := dispatcher.AssignRoutes(addresses, drivers)
	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, "error executing k-means algorithm: mock error", err.Error())
	mockKM.AssertExpectations(t)
	mockGeo.AssertExpectations(t)
}
