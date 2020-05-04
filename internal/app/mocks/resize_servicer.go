// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import domain "GoImageZip/internal/domain"
import mock "github.com/stretchr/testify/mock"

// ResizeServicer is an autogenerated mock type for the ResizeServicer type
type ResizeServicer struct {
	mock.Mock
}

// GetById provides a mock function with given fields: _a0
func (_m *ResizeServicer) GetById(_a0 string) (domain.ImageInfo, error) {
	ret := _m.Called(_a0)

	var r0 domain.ImageInfo
	if rf, ok := ret.Get(0).(func(string) domain.ImageInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(domain.ImageInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields:
func (_m *ResizeServicer) GetHistory() ([]domain.ImageInfo, error) {
	ret := _m.Called()

	var r0 []domain.ImageInfo
	if rf, ok := ret.Get(0).(func() []domain.ImageInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.ImageInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Resize provides a mock function with given fields: _a0
func (_m *ResizeServicer) Resize(_a0 domain.Image) (domain.ImageInfo, error) {
	ret := _m.Called(_a0)

	var r0 domain.ImageInfo
	if rf, ok := ret.Get(0).(func(domain.Image) domain.ImageInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(domain.ImageInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Image) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *ResizeServicer) Update(_a0 domain.ImageInfo) (domain.ImageInfo, error) {
	ret := _m.Called(_a0)

	var r0 domain.ImageInfo
	if rf, ok := ret.Get(0).(func(domain.ImageInfo) domain.ImageInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(domain.ImageInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.ImageInfo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}