// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\fahrizal.aziz_idstar\workshops\bookstore\bookstore-framework\pkg\generateToken.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJWTGenerator is a mock of JWTGenerator interface.
type MockJWTGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockJWTGeneratorMockRecorder
}

// MockJWTGeneratorMockRecorder is the mock recorder for MockJWTGenerator.
type MockJWTGeneratorMockRecorder struct {
	mock *MockJWTGenerator
}

// NewMockJWTGenerator creates a new mock instance.
func NewMockJWTGenerator(ctrl *gomock.Controller) *MockJWTGenerator {
	mock := &MockJWTGenerator{ctrl: ctrl}
	mock.recorder = &MockJWTGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTGenerator) EXPECT() *MockJWTGeneratorMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJWTGenerator) GenerateToken(userId uint, username, email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId, username, email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJWTGeneratorMockRecorder) GenerateToken(userId, username, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJWTGenerator)(nil).GenerateToken), userId, username, email)
}
