// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudtrust/common-service/v2/database/sqltypes (interfaces: SQLRow,CloudtrustDB,CloudtrustDBFactory)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/cloudtrustdb.go -package=mock -mock_names=SQLRow=SQLRow,CloudtrustDB=CloudtrustDB,CloudtrustDBFactory=CloudtrustDBFactory github.com/cloudtrust/common-service/v2/database/sqltypes SQLRow,CloudtrustDB,CloudtrustDBFactory
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	sqltypes "github.com/cloudtrust/common-service/v2/database/sqltypes"
	gomock "go.uber.org/mock/gomock"
)

// SQLRow is a mock of SQLRow interface.
type SQLRow struct {
	ctrl     *gomock.Controller
	recorder *SQLRowMockRecorder
	isgomock struct{}
}

// SQLRowMockRecorder is the mock recorder for SQLRow.
type SQLRowMockRecorder struct {
	mock *SQLRow
}

// NewSQLRow creates a new mock instance.
func NewSQLRow(ctrl *gomock.Controller) *SQLRow {
	mock := &SQLRow{ctrl: ctrl}
	mock.recorder = &SQLRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SQLRow) EXPECT() *SQLRowMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *SQLRow) Scan(dest ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range dest {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *SQLRowMockRecorder) Scan(dest ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*SQLRow)(nil).Scan), dest...)
}

// CloudtrustDB is a mock of CloudtrustDB interface.
type CloudtrustDB struct {
	ctrl     *gomock.Controller
	recorder *CloudtrustDBMockRecorder
	isgomock struct{}
}

// CloudtrustDBMockRecorder is the mock recorder for CloudtrustDB.
type CloudtrustDBMockRecorder struct {
	mock *CloudtrustDB
}

// NewCloudtrustDB creates a new mock instance.
func NewCloudtrustDB(ctrl *gomock.Controller) *CloudtrustDB {
	mock := &CloudtrustDB{ctrl: ctrl}
	mock.recorder = &CloudtrustDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *CloudtrustDB) EXPECT() *CloudtrustDBMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *CloudtrustDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltypes.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(sqltypes.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *CloudtrustDBMockRecorder) BeginTx(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*CloudtrustDB)(nil).BeginTx), ctx, opts)
}

// Close mocks base method.
func (m *CloudtrustDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *CloudtrustDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*CloudtrustDB)(nil).Close))
}

// Exec mocks base method.
func (m *CloudtrustDB) Exec(query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *CloudtrustDBMockRecorder) Exec(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*CloudtrustDB)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *CloudtrustDB) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *CloudtrustDBMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*CloudtrustDB)(nil).Ping))
}

// Query mocks base method.
func (m *CloudtrustDB) Query(query string, args ...any) (sqltypes.SQLRows, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(sqltypes.SQLRows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *CloudtrustDBMockRecorder) Query(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*CloudtrustDB)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *CloudtrustDB) QueryRow(query string, args ...any) sqltypes.SQLRow {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(sqltypes.SQLRow)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *CloudtrustDBMockRecorder) QueryRow(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*CloudtrustDB)(nil).QueryRow), varargs...)
}

// Stats mocks base method.
func (m *CloudtrustDB) Stats() sql.DBStats {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats")
	ret0, _ := ret[0].(sql.DBStats)
	return ret0
}

// Stats indicates an expected call of Stats.
func (mr *CloudtrustDBMockRecorder) Stats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*CloudtrustDB)(nil).Stats))
}

// CloudtrustDBFactory is a mock of CloudtrustDBFactory interface.
type CloudtrustDBFactory struct {
	ctrl     *gomock.Controller
	recorder *CloudtrustDBFactoryMockRecorder
	isgomock struct{}
}

// CloudtrustDBFactoryMockRecorder is the mock recorder for CloudtrustDBFactory.
type CloudtrustDBFactoryMockRecorder struct {
	mock *CloudtrustDBFactory
}

// NewCloudtrustDBFactory creates a new mock instance.
func NewCloudtrustDBFactory(ctrl *gomock.Controller) *CloudtrustDBFactory {
	mock := &CloudtrustDBFactory{ctrl: ctrl}
	mock.recorder = &CloudtrustDBFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *CloudtrustDBFactory) EXPECT() *CloudtrustDBFactoryMockRecorder {
	return m.recorder
}

// OpenDatabase mocks base method.
func (m *CloudtrustDBFactory) OpenDatabase() (sqltypes.CloudtrustDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenDatabase")
	ret0, _ := ret[0].(sqltypes.CloudtrustDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenDatabase indicates an expected call of OpenDatabase.
func (mr *CloudtrustDBFactoryMockRecorder) OpenDatabase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenDatabase", reflect.TypeOf((*CloudtrustDBFactory)(nil).OpenDatabase))
}
