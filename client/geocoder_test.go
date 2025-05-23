package client

import (
	"errors"
	"testing"

	gmock "github.com/lalsaady/dispatcher/mock"
	"github.com/lalsaady/dispatcher/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCoords(t *testing.T) {
	mockGeo := gmock.NewMockGeocoder()
	mockResponse := model.Points{
		Lat: 42.123456,
		Lon: -80.234567,
	}
	mockGeo.On("GetCoords", mock.Anything).Return(mockResponse, nil)
	result, err := mockGeo.GetCoords("Fake Address")
	assert.NoError(t, err)
	assert.Equal(t, mockResponse, result)
}

func TestGetCoords_Error(t *testing.T) {
	mockGeo := gmock.NewMockGeocoder()
	mockGeo.On("GetCoords", mock.Anything).Return(model.Points{}, errors.New("error"))
	result, err := mockGeo.GetCoords("Fake Address")
	assert.Error(t, err)
	assert.Equal(t, model.Points{}, result)
}
