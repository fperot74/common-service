package configuration

import (
	"context"
	"database/sql"

	"github.com/cloudtrust/common-service/v2/database/sqltypes"
	"github.com/cloudtrust/common-service/v2/log"
)

const (
	selectBothConfigsStmt = `SELECT configuration, admin_configuration FROM realm_configuration WHERE realm_id = ? AND configuration IS NOT NULL AND admin_configuration IS NOT NULL`
	selectConfigStmt      = `SELECT configuration FROM realm_configuration WHERE realm_id = ? AND configuration IS NOT NULL`
	selectAdminConfigStmt = `SELECT admin_configuration FROM realm_configuration WHERE realm_id = ? AND admin_configuration IS NOT NULL`
	selectAllAuthzStmt    = `SELECT realm_id, group_name, action, target_realm_id, target_group_name FROM authorizations;`
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
