package dispatcher

import (
	"errors"
	"testing"

	KMeansMock "github.com/lalsaady/dispatcher/mock"
	"github.com/muesli/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDispatcherInvalid(t *testing.T) {
	dispatcher := NewDispatcher(nil)
	routes, err := dispatcher.AssignRoutes([]Location{}, []string{"Alice", "Bob"})
	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, err.Error(), "no orders provided")
	routes, err = dispatcher.AssignRoutes([]Location{
		{ID: 1, Address: "123 Main St, Cleveland, OH", Lat: 40.988612, Lon: -80.698871},
	}, []string{})
	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, err.Error(), "no drivers provided")
}

func TestDispatcherValid(t *testing.T) {
	orders := []Location{
		{ID: 1, Address: "123 Main St, Cleveland, OH", Lat: 40.988612, Lon: -80.698871},
		{ID: 2, Address: "345 Dummy Ave, Cleveland, OH", Lat: 41.073612, Lon: -80.609771},
	}
	drivers := []string{"Alice", "Bob"}
	dispatcher := NewDispatcher(nil)
	routes, err := dispatcher.AssignRoutes(orders, drivers)
	assert.NoError(t, err)
	assert.NotNil(t, routes)
	assert.Equal(t, len(routes), 2)
	assert.Equal(t, len(routes["Alice"]), 1)
	assert.Equal(t, len(routes["Bob"]), 1)
}

func TestDispatcherError(t *testing.T) {
	orders := []Location{
		{ID: 1, Address: "123 Main St, Cleveland, OH", Lat: 40.988612, Lon: -80.698871},
		{ID: 2, Address: "345 Dummy Ave, Cleveland, OH", Lat: 41.073612, Lon: -80.609771},
	}
	drivers := []string{"Alice", "Bob"}

	mockKM := KMeansMock.NewMockKMeans()
	mockKM.On("Partition", mock.Anything, mock.Anything).Return(clusters.Clusters{}, errors.New("mock error"))

	dispatcher := NewDispatcher(mockKM)
	routes, err := dispatcher.AssignRoutes(orders, drivers)
	assert.Error(t, err)
	assert.Nil(t, routes)
	assert.Equal(t, err.Error(), "error executing k-means algorithm")
	mockKM.AssertExpectations(t)
}
