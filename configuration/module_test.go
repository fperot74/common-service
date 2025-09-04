package configuration

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/v2/configuration/mock"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type dbMocks struct {
	mockCtrl *gomock.Controller
	db       *mock.CloudtrustDB
	sqlRow   *mock.SQLRow
	sqlRows  *mock.SQLRows
	tx       *mock.Transaction
	logger   log.Logger
}

func newDbMocks(t *testing.T) *dbMocks {
	var mockCtrl = gomock.NewController(t)
	return &dbMocks{
		mockCtrl: mockCtrl,
		db:       mock.NewCloudtrustDB(mockCtrl),
		sqlRow:   mock.NewSQLRow(mockCtrl),
		sqlRows:  mock.NewSQLRows(mockCtrl),
		tx:       mock.NewTransaction(mockCtrl),
		logger:   log.NewNopLogger(),
	}
}

func (m *dbMocks) NewConfigurationReaderDBModule(actions ...[]string) *ConfigurationReaderDBModule {
	return NewConfigurationReaderDBModule(m.db, m.logger, actions...)
}

func (m *dbMocks) Finish() {
	m.mockCtrl.Finish()
}

func TestGetRealmConfigurations(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var realmID = "myrealm"
	var ctx = context.TODO()
	var module = mocks.NewConfigurationReaderDBModule()

	t.Run("SQL error", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(errors.New("SQL error"))
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows)
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Invalid JSON", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0]).(*string) = `{`
			*(dest[1]).(*string) = `{}`
			return nil
		})
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0]).(*string) = `{}`
			*(dest[1]).(*string) = `{}`
			return nil
		})
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.Nil(t, err)
	})
}

func TestGetConfiguration(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var realmID = "myrealm"
	var ctx = context.TODO()
	var module = mocks.NewConfigurationReaderDBModule()

	t.Run("SQL error", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).Return(errors.New("SQL error"))
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(ptrConfig *string) error {
			*ptrConfig = `{}`
			return nil
		})
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.Nil(t, err)
	})
}

func TestGetAdminConfiguration(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var realmID = "myrealm"
	var ctx = context.TODO()
	var module = mocks.NewConfigurationReaderDBModule()

	t.Run("SQL error", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).Return(errors.New("SQL error"))
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mocks.db.EXPECT().QueryRow(gomock.Any(), realmID).Return(mocks.sqlRow)
		mocks.sqlRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(ptrConfig *string) error {
			*ptrConfig = `{}`
			return nil
		})
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.Nil(t, err)
	})
}

func TestGetContextKeyConfiguration(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var customerRealm = "customer-realm"
	var ctx = context.TODO()

	var module = mocks.NewConfigurationReaderDBModule()

	t.Run("Query fails", func(t *testing.T) {
		var sqlError = errors.New("SQL error")
		mocks.db.EXPECT().Query(gomock.Any(), customerRealm).Return(nil, sqlError)

		var _, err = module.GetContextKeyConfiguration(ctx, customerRealm)
		assert.Equal(t, sqlError, err)
	})

	// Now, all scenarii will be executed such as query always returns a sqlRows value
	mocks.db.EXPECT().Query(gomock.Any(), customerRealm).Return(mocks.sqlRows, nil).AnyTimes()
	mocks.sqlRows.EXPECT().Close().AnyTimes()

	t.Run("Scan fails", func(t *testing.T) {
		var scanError = errors.New("scan error")
		mocks.sqlRows.EXPECT().Next().Return(true)
		mocks.sqlRows.EXPECT().Scan(gomock.Any()).Return(scanError)

		var _, err = module.GetContextKeyConfiguration(ctx, customerRealm)
		assert.Equal(t, scanError, err)
	})

	t.Run("Error during iteration", func(t *testing.T) {
		iterationErr := errors.New("iteration error")
		mocks.sqlRows.EXPECT().Next().Return(false)
		mocks.sqlRows.EXPECT().Err().Return(iterationErr)

		var _, err = module.GetContextKeyConfiguration(ctx, customerRealm)
		assert.Equal(t, iterationErr, err)
	})

	t.Run("Invalid JSON found in DB", func(t *testing.T) {
		mocks.sqlRows.EXPECT().Next().Return(true)
		mocks.sqlRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0]).(*string) = "uuid1"
			*(dest[1]).(*string) = "label1"
			*(dest[2]).(*string) = "identities-realm"
			*(dest[3]).(*string) = customerRealm
			*(dest[4]).(*string) = "{"
			return nil
		})

		var _, err = module.GetContextKeyConfiguration(ctx, customerRealm)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "JSON input")
	})

	t.Run("Query returns 2 rows", func(t *testing.T) {
		mockSelectContextKeyResult(mocks, []RealmContextKey{
			{
				ID:              "uuid1",
				Label:           "label1",
				IdentitiesRealm: "identities-realm1",
				CustomerRealm:   customerRealm,
				Config:          ContextKeyConfiguration{},
			},
			{
				ID:              "uuid2",
				Label:           "label2",
				IdentitiesRealm: "identities-realm2",
				CustomerRealm:   customerRealm,
				Config:          ContextKeyConfiguration{},
			},
		})

		var res, err = module.GetContextKeyConfiguration(ctx, customerRealm)
		assert.Nil(t, err)
		assert.Len(t, res, 2)
	})
}

func mockSelectContextKeyResult(mocks *dbMocks, items []RealmContextKey) {
	for _, item := range items {
		var bytes, _ = json.Marshal(item.Config)
		mocks.sqlRows.EXPECT().Next().Return(true)
		mocks.sqlRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0]).(*string) = item.ID
			*(dest[1]).(*string) = item.Label
			*(dest[2]).(*string) = item.IdentitiesRealm
			*(dest[3]).(*string) = item.CustomerRealm
			*(dest[4]).(*string) = string(bytes)
			return nil
		})
	}
	mocks.sqlRows.EXPECT().Next().Return(false)
	mocks.sqlRows.EXPECT().Err()
}

func TestSetContextKeyConfiguration(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var customerRealm = "customer-realm"
	var realmContext1 = RealmContextKey{
		ID:              "uuid1",
		Label:           "label1",
		IdentitiesRealm: "identities-realm",
	}
	var emptyContextKeys = []RealmContextKey{}
	var ctx = context.TODO()

	var module = mocks.NewConfigurationReaderDBModule()

	t.Run("Select query fails", func(t *testing.T) {
		var sqlError = errors.New("SQL error")
		mocks.db.EXPECT().Query(gomock.Any(), customerRealm).Return(nil, sqlError)

		var err = module.SetContextKeyConfiguration(ctx, customerRealm, emptyContextKeys)
		assert.Equal(t, sqlError, err)
	})

	// All select queries will return a valid result
	mocks.db.EXPECT().Query(gomock.Any(), customerRealm).Return(mocks.sqlRows, nil).AnyTimes()
	mocks.sqlRows.EXPECT().Close().AnyTimes()

	t.Run("Begin transaction fails", func(t *testing.T) {
		var sqlError = errors.New("tx error")
		mockSelectContextKeyResult(mocks, []RealmContextKey{})
		mocks.db.EXPECT().BeginTx(ctx, nil).Return(nil, sqlError)

		var err = module.SetContextKeyConfiguration(ctx, customerRealm, emptyContextKeys)
		assert.Equal(t, sqlError, err)
	})

	// Now all transactions will be created without error
	mocks.db.EXPECT().BeginTx(ctx, nil).Return(mocks.tx, nil).AnyTimes()
	mocks.tx.EXPECT().Close().AnyTimes()

	t.Run("Deletion", func(t *testing.T) {
		t.Run("Delete failure", func(t *testing.T) {
			var deleteError = errors.New("deeelt eorrr")
			mockSelectContextKeyResult(mocks, []RealmContextKey{realmContext1})
			mocks.tx.EXPECT().Exec(deleteContextKey, realmContext1.ID).Return(nil, deleteError)

			var err = module.SetContextKeyConfiguration(ctx, customerRealm, emptyContextKeys)
			assert.Equal(t, deleteError, err)
		})

		t.Run("Delete success", func(t *testing.T) {
			mockSelectContextKeyResult(mocks, []RealmContextKey{realmContext1})
			mocks.tx.EXPECT().Exec(deleteContextKey, realmContext1.ID).Return(nil, nil)
			mocks.tx.EXPECT().Commit()

			var err = module.SetContextKeyConfiguration(ctx, customerRealm, emptyContextKeys)
			assert.Nil(t, err)
		})
	})
	t.Run("Creation/Update", func(t *testing.T) {
		var createContextKeys = []RealmContextKey{realmContext1}

		t.Run("Creation/update failure", func(t *testing.T) {
			var creationError = errors.New("craetoin eorrr")
			mockSelectContextKeyResult(mocks, []RealmContextKey{})
			mocks.tx.EXPECT().Exec(setContextKey, gomock.Any()).Return(nil, creationError)

			var err = module.SetContextKeyConfiguration(ctx, customerRealm, createContextKeys)
			assert.Equal(t, creationError, err)
		})
		t.Run("Creation/update success", func(t *testing.T) {
			mockSelectContextKeyResult(mocks, []RealmContextKey{})
			mocks.tx.EXPECT().Exec(setContextKey, gomock.Any()).Return(nil, nil)
			mocks.tx.EXPECT().Commit().Return(nil)

			var err = module.SetContextKeyConfiguration(ctx, customerRealm, createContextKeys)
			assert.Nil(t, err)
		})
		t.Run("Commit fails", func(t *testing.T) {
			var commitError = errors.New("mmcoit eorrr")
			mockSelectContextKeyResult(mocks, []RealmContextKey{})
			mocks.tx.EXPECT().Exec(setContextKey, gomock.Any()).Return(nil, nil)
			mocks.tx.EXPECT().Commit().Return(commitError)

			var err = module.SetContextKeyConfiguration(ctx, customerRealm, createContextKeys)
			assert.Equal(t, commitError, err)
		})
	})
}

func TestGetAuthorizations(t *testing.T) {
	var mocks = newDbMocks(t)
	defer mocks.Finish()

	var notInScopeAction = "ACTION#1"
	var allowedAction = "ACTION#2"
	var actions = []string{allowedAction}
	var ctx = context.TODO()

	var module = mocks.NewConfigurationReaderDBModule(actions)

	t.Run("Query fails", func(t *testing.T) {
		var sqlError = errors.New("SQL error")
		mocks.db.EXPECT().Query(gomock.Any()).Return(nil, sqlError)

		var _, err = module.GetAuthorizations(ctx)
		assert.Equal(t, sqlError, err)
	})

	// Now, query will always be successful
	mocks.db.EXPECT().Query(gomock.Any()).Return(mocks.sqlRows, nil).AnyTimes()
	mocks.sqlRows.EXPECT().Close().AnyTimes()

	t.Run("scan fails", func(t *testing.T) {
		var scanError = errors.New("scan error")
		mocks.sqlRows.EXPECT().Next().Return(true)
		mocks.sqlRows.EXPECT().Scan(gomock.Any()).Return(scanError)

		var _, err = module.GetAuthorizations(ctx)
		assert.Equal(t, scanError, err)
	})

	t.Run("error during iteration", func(t *testing.T) {
		iterationErr := errors.New("iteration error")
		mocks.sqlRows.EXPECT().Next().Return(false)
		mocks.sqlRows.EXPECT().Err().Return(iterationErr)

		_, err := module.GetAuthorizations(ctx)
		assert.Equal(t, iterationErr, err)
	})

	t.Run("Query returns 2 rows", func(t *testing.T) {
		gomock.InOrder(
			mocks.sqlRows.EXPECT().Next().Return(true),
			mocks.sqlRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
				*(dest[0]).(*string) = "realm#1"
				*(dest[1]).(*string) = "group#1"
				*(dest[2]).(*string) = notInScopeAction
				*(dest[3]).(*sql.NullString) = sql.NullString{Valid: false}
				*(dest[4]).(*sql.NullString) = sql.NullString{Valid: false}
				return nil
			}),
			mocks.sqlRows.EXPECT().Next().Return(true),
			mocks.sqlRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
				*(dest[0]).(*string) = "realm#2"
				*(dest[1]).(*string) = "group#2"
				*(dest[2]).(*string) = allowedAction
				*(dest[3]).(*sql.NullString) = sql.NullString{Valid: true, String: "targetRealm"}
				*(dest[4]).(*sql.NullString) = sql.NullString{Valid: true, String: "targetGroup"}
				return nil
			}),
			mocks.sqlRows.EXPECT().Next().Return(false),
			mocks.sqlRows.EXPECT().Err(),
		)

		var res, err = module.GetAuthorizations(ctx)
		assert.Nil(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, allowedAction, *res[0].Action)
	})
}

func TestIsInScope(t *testing.T) {
	var module = NewConfigurationReaderDBModule(nil, nil)
	assert.True(t, module.isInAuthorizationScope("any-auth-will-be-considered-in-scope"))

	module = NewConfigurationReaderDBModule(nil, nil, []string{"autun"}, []string{"auth2", "auth3"})
	assert.False(t, module.isInAuthorizationScope("any-auth-will-be-considered-in-scope"))
	assert.True(t, module.isInAuthorizationScope("autun"))
	assert.True(t, module.isInAuthorizationScope("auth2"))
	assert.True(t, module.isInAuthorizationScope("auth3"))
}
