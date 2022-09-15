package io

import "github.com/stretchr/testify/mock"


type MockReader struct {
	mock.Mock
}


func (m *MockReader) Read(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}
