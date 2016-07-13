
package testify

import (
	mock "github.com/stretchr/testify/mock"

	url "net/url"
)

type MockFuelCan struct {
	mock.Mock
}


func (m *MockFuelCan) Capacity() float32 {
	ret := m.Called()
	
	var r0 float32

	if r0f, ok := ret.Get(0).(func() float32); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).(float32)
	}

	return r0
}

func (m *MockFuelCan) Fill(v Vehicle) Vehicle {
	ret := m.Called(v)
	
	var r0 Vehicle

	if r0f, ok := ret.Get(0).(func(Vehicle) Vehicle); ok {
			r0 = r0f(v)
	} else {
			r0 = ret.Get(0).(Vehicle)
	}

	return r0
}

func (m *MockFuelCan) Qty() float32 {
	ret := m.Called()
	
	var r0 float32

	if r0f, ok := ret.Get(0).(func() float32); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).(float32)
	}

	return r0
}


type MockVehicle struct {
	mock.Mock
}


func (m *MockVehicle) Attach(wheels ...Wheel) {
	m.Called(wheels)
	

	return 
}

func (m *MockVehicle) Drive(dir string,dist float32) url.URL {
	ret := m.Called(dir,dist)
	
	var r0 url.URL

	if r0f, ok := ret.Get(0).(func(string,float32) url.URL); ok {
			r0 = r0f(dir,dist)
	} else {
			r0 = ret.Get(0).(url.URL)
	}

	return r0
}

func (m *MockVehicle) EnterLeave(m0 map[int]string,b1 bool) {
	m.Called(m0,b1)
	

	return 
}

func (m *MockVehicle) Occupants() map[int]string {
	ret := m.Called()
	
	var r0 map[int]string

	if r0f, ok := ret.Get(0).(func() map[int]string); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).(map[int]string)
	}

	return r0
}

func (m *MockVehicle) Range() int {
	ret := m.Called()
	
	var r0 int

	if r0f, ok := ret.Get(0).(func() int); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).(int)
	}

	return r0
}

func (m *MockVehicle) Refuel(fc FuelCan) *FuelCan {
	ret := m.Called(fc)
	
	var r0 *FuelCan

	if r0f, ok := ret.Get(0).(func(FuelCan) *FuelCan); ok {
			r0 = r0f(fc)
	} else {
			r0 = ret.Get(0).(*FuelCan)
	}

	return r0
}

func (m *MockVehicle) Wheels() []Wheel {
	ret := m.Called()
	
	var r0 []Wheel

	if r0f, ok := ret.Get(0).(func() []Wheel); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).([]Wheel)
	}

	return r0
}


type MockWheel struct {
	mock.Mock
}


func (m *MockWheel) Diameter() float32 {
	ret := m.Called()
	
	var r0 float32

	if r0f, ok := ret.Get(0).(func() float32); ok {
			r0 = r0f()
	} else {
			r0 = ret.Get(0).(float32)
	}

	return r0
}

