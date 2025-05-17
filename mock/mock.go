package mock

import (
	"github.com/muesli/clusters"
	"github.com/stretchr/testify/mock"
)

type MockKMeans struct {
	mock.Mock
}

func (m *MockKMeans) Partition(observations clusters.Observations, k int) (clusters.Clusters, error) {
	args := m.Called(observations, k)
	return args.Get(0).(clusters.Clusters), args.Error(1)
}

func NewMockKMeans() *MockKMeans {
	return &MockKMeans{}
}
