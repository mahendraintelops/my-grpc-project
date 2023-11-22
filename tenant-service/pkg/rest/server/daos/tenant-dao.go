package daos

import (
	"errors"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TenantDao struct {
	db *gorm.DB
}

func NewTenantDao() (*TenantDao, error) {
	sqlClient, err := sqls.InitGORMSQLiteDB()
	if err != nil {
		return nil, err
	}
	err = sqlClient.DB.AutoMigrate(models.Tenant{})
	if err != nil {
		return nil, err
	}
	return &TenantDao{
		db: sqlClient.DB,
	}, nil
}

func (tenantDao *TenantDao) CreateTenant(m *models.Tenant) (*models.Tenant, error) {
	if err := tenantDao.db.Create(&m).Error; err != nil {
		log.Debugf("failed to create tenant: %v", err)
		return nil, err
	}

	log.Debugf("tenant created")
	return m, nil
}

func (tenantDao *TenantDao) ListTenants() ([]*models.Tenant, error) {
	var tenants []*models.Tenant

	// TODO populate associations here with your own logic - https://gorm.io/docs/belongs_to.html
	if err := tenantDao.db.Find(&tenants).Error; err != nil {
		log.Debugf("failed to list tenants: %v", err)
		return nil, err
	}

	log.Debugf("tenant listed")
	return tenants, nil
}

func (tenantDao *TenantDao) GetTenant(id int64) (*models.Tenant, error) {
	var m *models.Tenant
	if err := tenantDao.db.Where("id = ?", id).First(&m).Error; err != nil {
		log.Debugf("failed to get tenant: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("tenant retrieved")
	return m, nil
}

func (tenantDao *TenantDao) UpdateTenant(id int64, m *models.Tenant) (*models.Tenant, error) {
	if id == 0 {
		return nil, errors.New("invalid tenant ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	var tenant *models.Tenant
	if err := tenantDao.db.Where("id = ?", id).First(&tenant).Error; err != nil {
		log.Debugf("failed to find tenant for update: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	if err := tenantDao.db.Save(&m).Error; err != nil {
		log.Debugf("failed to update tenant: %v", err)
		return nil, err
	}
	log.Debugf("tenant updated")
	return m, nil
}

func (tenantDao *TenantDao) DeleteTenant(id int64) error {
	var m *models.Tenant
	if err := tenantDao.db.Where("id = ?", id).Delete(&m).Error; err != nil {
		log.Debugf("failed to delete tenant: %v", err)
		return err
	}

	log.Debugf("tenant deleted")
	return nil
}
