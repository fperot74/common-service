package configuration

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/cloudtrust/common-service/v2/database/sqltypes"
	"github.com/cloudtrust/common-service/v2/log"
)

const (
	selectBothConfigsStmt  = `SELECT configuration, admin_configuration FROM realm_configuration WHERE realm_id = ? AND configuration IS NOT NULL AND admin_configuration IS NOT NULL`
	selectConfigStmt       = `SELECT configuration FROM realm_configuration WHERE realm_id = ? AND configuration IS NOT NULL`
	selectAdminConfigStmt  = `SELECT admin_configuration FROM realm_configuration WHERE realm_id = ? AND admin_configuration IS NOT NULL`
	selectContextKeyConfig = `SELECT id, label, identities_realm, customer_realm, configuration FROM context_key_configuration WHERE customer_realm = ?`
	deleteContextKey       = `DELETE from context_key_configuration WHERE id = ?`
	setContextKey          = `
		INSERT INTO context_key_configuration (id, label, identities_realm, customer_realm, configuration)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE id=?, label=?, identities_realm=?, customer_realm=?, configuration=?
	`
	selectAllAuthzStmt = `SELECT realm_id, group_name, action, target_realm_id, target_group_name FROM authorizations;`
)

// ConfigurationReaderDBModule struct
type ConfigurationReaderDBModule struct {
	db        sqltypes.CloudtrustDB
	authScope map[string]bool
	logger    log.Logger
}

// NewConfigurationReaderDBModule returns a ConfigurationDB module.
func NewConfigurationReaderDBModule(db sqltypes.CloudtrustDB, logger log.Logger, actions ...[]string) *ConfigurationReaderDBModule {
	var authScope map[string]bool
	if len(actions) > 0 {
		authScope = make(map[string]bool)
		for _, actionSet := range actions {
			for _, filter := range actionSet {
				authScope[filter] = true
			}
		}
	}
	return &ConfigurationReaderDBModule{
		db:        db,
		authScope: authScope,
		logger:    logger,
	}
}

// GetRealmConfigurations returns both configuration and admin configuration of a realm
func (c *ConfigurationReaderDBModule) GetRealmConfigurations(ctx context.Context, realmID string) (RealmConfiguration, RealmAdminConfiguration, error) {
	var configJSON, adminConfigJSON string
	row := c.db.QueryRow(selectBothConfigsStmt, realmID)

	switch err := row.Scan(&configJSON, &adminConfigJSON); err {
	case sql.ErrNoRows:
		c.logger.Warn(ctx, "msg", "Realm Configuration not found in DB", "err", err.Error())
		return RealmConfiguration{}, RealmAdminConfiguration{}, err

	default:
		if err != nil {
			return RealmConfiguration{}, RealmAdminConfiguration{}, err
		}

		realmConf, err := NewRealmConfiguration(configJSON)
		if err != nil {
			return RealmConfiguration{}, RealmAdminConfiguration{}, err
		}

		realmAdminConf, err := NewRealmAdminConfiguration(adminConfigJSON)
		return realmConf, realmAdminConf, err
	}
}

// GetConfiguration returns a realm configuration
func (c *ConfigurationReaderDBModule) GetConfiguration(ctx context.Context, realmID string) (RealmConfiguration, error) {
	var configJSON string
	row := c.db.QueryRow(selectConfigStmt, realmID)

	switch err := row.Scan(&configJSON); err {
	case sql.ErrNoRows:
		c.logger.Warn(ctx, "msg", "Realm Configuration not found in DB", "err", err.Error())
		return RealmConfiguration{}, err

	default:
		if err != nil {
			return RealmConfiguration{}, err
		}

		return NewRealmConfiguration(configJSON)
	}
}

// GetAdminConfiguration returns a realm admin configuration
func (c *ConfigurationReaderDBModule) GetAdminConfiguration(ctx context.Context, realmID string) (RealmAdminConfiguration, error) {
	var configJSON string
	row := c.db.QueryRow(selectAdminConfigStmt, realmID)

	var err = row.Scan(&configJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			c.logger.Warn(ctx, "msg", "Realm Admin Configuration not found in DB", "err", err.Error())
		}
		return RealmAdminConfiguration{}, err
	}
	return NewRealmAdminConfiguration(configJSON)
}

// GetContextKeyConfiguration returns context key configuration for a given customer realm
func (c *ConfigurationReaderDBModule) GetContextKeyConfiguration(ctx context.Context, customerRealm string) ([]RealmContextKey, error) {
	rows, err := c.db.Query(selectContextKeyConfig, customerRealm)
	if err != nil {
		c.logger.Warn(ctx, "msg", "Can't get context key configuration", "realm", customerRealm, "err", err.Error())
		return nil, err
	}
	defer rows.Close()

	var res []RealmContextKey
	for rows.Next() {
		ctxKeyConf, err := c.scanContextKeyConfiguration(rows)
		if err != nil {
			c.logger.Warn(ctx, "msg", "Can't get context key configuration. Scan failed", "realm", customerRealm, "err", err.Error())
			return nil, err
		}
		res = append(res, ctxKeyConf)
	}
	if err = rows.Err(); err != nil {
		c.logger.Warn(ctx, "msg", "Can't get context key configuration. Failed to iterate on every items", "realm", customerRealm, "err", err.Error())
		return nil, err
	}

	return res, nil
}

func (c *ConfigurationReaderDBModule) scanContextKeyConfiguration(scanner sqltypes.SQLRow) (RealmContextKey, error) {
	var (
		id              string
		label           string
		identitiesRealm string
		customerRealm   string
		configJSON      string
	)

	err := scanner.Scan(&id, &label, &identitiesRealm, &customerRealm, &configJSON)
	if err != nil {
		return RealmContextKey{}, err
	}

	config, err := NewContextKeyConfiguration(configJSON)
	if err != nil {
		return RealmContextKey{}, err
	}

	return RealmContextKey{
		ID:              id,
		Label:           label,
		IdentitiesRealm: identitiesRealm,
		CustomerRealm:   customerRealm,
		Config:          config,
	}, nil
}

func (c *ConfigurationReaderDBModule) getContextKeyUUIDSet(contextKeys []RealmContextKey) map[string]struct{} {
	existingKeys := map[string]struct{}{}
	for _, ctxKey := range contextKeys {
		existingKeys[ctxKey.ID] = struct{}{}
	}
	return existingKeys
}

func (c *ConfigurationReaderDBModule) getExistingContextKeyUUID(ctx context.Context, customerRealm string) (map[string]struct{}, error) {
	existing, err := c.GetContextKeyConfiguration(ctx, customerRealm)
	if err != nil {
		c.logger.Warn(ctx, "msg", "Can't get existing context key configuration", "realm", customerRealm, "err", err.Error())
		return nil, err
	}
	return c.getContextKeyUUIDSet(existing), nil
}

// SetContextKeyConfiguration sets the context key configuration for a given customer realm
func (c *ConfigurationReaderDBModule) SetContextKeyConfiguration(ctx context.Context, customerRealm string, contextKeys []RealmContextKey) error {
	providedKeys := c.getContextKeyUUIDSet(contextKeys)
	existingKeys, err := c.getExistingContextKeyUUID(ctx, customerRealm)
	if err != nil {
		return err
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		c.logger.Warn(ctx, "msg", "Can't start transaction to set context key", "realm", customerRealm, "err", err.Error())
		return err
	}
	defer tx.Close()

	for existingKey := range existingKeys {
		if _, ok := providedKeys[existingKey]; !ok {
			if _, err := tx.Exec(deleteContextKey, existingKey); err != nil {
				c.logger.Warn(ctx, "msg", "Can't delete a context key configuration", "id", existingKey, "err", err.Error())
				return err
			}
		}
	}

	for _, ck := range contextKeys {
		configJSON, err := json.Marshal(ck.Config)
		if err != nil {
			c.logger.Warn(ctx, "msg", "Can't convert to JSON", "realm", customerRealm, "id", ck.ID, "err", err.Error())
			return err
		}
		if _, err := tx.Exec(setContextKey, ck.ID, ck.Label, ck.IdentitiesRealm, ck.CustomerRealm, configJSON); err != nil {
			c.logger.Warn(ctx, "msg", "Can't set context key in db", "realm", customerRealm, "id", ck.ID, "err", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		c.logger.Warn(ctx, "msg", "Can't set context key in db. Commit transaction failed", "realm", customerRealm, "err", err.Error())
		return err
	}

	return nil
}

// GetAuthorizations returns authorizations
func (c *ConfigurationReaderDBModule) GetAuthorizations(ctx context.Context) ([]Authorization, error) {
	// Get Authorizations from DB
	rows, err := c.db.Query(selectAllAuthzStmt)
	if err != nil {
		c.logger.Warn(ctx, "msg", "Can't get authorizations", "err", err.Error())
		return nil, err
	}
	defer rows.Close()

	var authz Authorization
	var res = make([]Authorization, 0)
	for rows.Next() {
		authz, err = c.scanAuthorization(rows)
		if err != nil {
			c.logger.Warn(ctx, "msg", "Can't get authorizations. Scan failed", "err", err.Error())
			return nil, err
		}
		if c.isInAuthorizationScope(*authz.Action) {
			res = append(res, authz)
		}
	}
	if err = rows.Err(); err != nil {
		c.logger.Warn(ctx, "msg", "Can't get authorizations. Failed to iterate on every items", "err", err.Error())
		return nil, err
	}

	return res, nil
}

func (c *ConfigurationReaderDBModule) scanAuthorization(scanner sqltypes.SQLRow) (Authorization, error) {
	var (
		realmID         string
		groupName       string
		action          string
		targetGroupName sql.NullString
		targetRealmID   sql.NullString
	)

	err := scanner.Scan(&realmID, &groupName, &action, &targetRealmID, &targetGroupName)
	if err != nil {
		return Authorization{}, err
	}

	var authz = Authorization{
		RealmID:   &realmID,
		GroupName: &groupName,
		Action:    &action,
	}

	if targetRealmID.Valid {
		authz.TargetRealmID = &targetRealmID.String
	}

	if targetGroupName.Valid {
		authz.TargetGroupName = &targetGroupName.String
	}

	return authz, nil
}

func (c *ConfigurationReaderDBModule) isInAuthorizationScope(action string) bool {
	if c.authScope != nil {
		if _, ok := c.authScope[action]; !ok {
			return false
		}
	}
	return true
}
