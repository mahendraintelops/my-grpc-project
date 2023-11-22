package services

import (
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/daos"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/models"
)

type TenantService struct {
	tenantDao *daos.TenantDao
}

func NewTenantService() (*TenantService, error) {
	tenantDao, err := daos.NewTenantDao()
	if err != nil {
		return nil, err
	}
	return &TenantService{
		tenantDao: tenantDao,
	}, nil
}

func (tenantService *TenantService) CreateTenant(tenant *models.Tenant) (*models.Tenant, error) {
	return tenantService.tenantDao.CreateTenant(tenant)
}

func (tenantService *TenantService) ListTenants() ([]*models.Tenant, error) {
	return tenantService.tenantDao.ListTenants()
}

func (tenantService *TenantService) GetTenant(id int64) (*models.Tenant, error) {
	return tenantService.tenantDao.GetTenant(id)
}

func (tenantService *TenantService) UpdateTenant(id int64, tenant *models.Tenant) (*models.Tenant, error) {
	return tenantService.tenantDao.UpdateTenant(id, tenant)
}

func (tenantService *TenantService) DeleteTenant(id int64) error {
	return tenantService.tenantDao.DeleteTenant(id)
}
