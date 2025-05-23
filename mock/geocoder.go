package mock

import (
	"github.com/lalsaady/dispatcher/model"
	"github.com/stretchr/testify/mock"
)

type MockGeocoder struct {
	mock.Mock
}

func (m *MockGeocoder) GetCoords(address string) (model.Points, error) {
	args := m.Called(address)
	return args.Get(0).(model.Points), args.Error(1)
}

func NewMockGeocoder() *MockGeocoder {
	return &MockGeocoder{}
}
