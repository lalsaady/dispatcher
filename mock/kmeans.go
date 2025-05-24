package mock

import (
	"github.com/lalsaady/dispatcher/model"
	"github.com/muesli/clusters"
	"github.com/stretchr/testify/mock"
)

type MockKMeans struct {
	mock.Mock
}

func (m *MockKMeans) Partition(orders []model.Location, k int, hub model.Location) (clusters.Clusters, error) {
	args := m.Called(orders, k, hub)
	return args.Get(0).(clusters.Clusters), args.Error(1)
}

func NewMockKMeans() *MockKMeans {
	return &MockKMeans{}
}
