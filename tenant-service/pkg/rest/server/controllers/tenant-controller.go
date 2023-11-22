package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/models"
	"github.com/mahendraintelops/my-grpc-project/tenant-service/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type TenantController struct {
	tenantService *services.TenantService
}

func NewTenantController() (*TenantController, error) {
	tenantService, err := services.NewTenantService()
	if err != nil {
		return nil, err
	}
	return &TenantController{
		tenantService: tenantService,
	}, nil
}

func (tenantController *TenantController) CreateTenant(context *gin.Context) {
	// validate input
	var input models.Tenant
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger tenant creation
	tenantCreated, err := tenantController.tenantService.CreateTenant(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, tenantCreated)
}

func (tenantController *TenantController) ListTenants(context *gin.Context) {
	// trigger all tenants fetching
	tenants, err := tenantController.tenantService.ListTenants()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, tenants)
}

func (tenantController *TenantController) FetchTenant(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger tenant fetching
	tenant, err := tenantController.tenantService.GetTenant(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("tenant.id", strconv.FormatInt(tenant.Id, 10)))
	}

	context.JSON(http.StatusOK, tenant)
}

func (tenantController *TenantController) UpdateTenant(context *gin.Context) {
	// validate input
	var input models.Tenant
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger tenant update
	if _, err := tenantController.tenantService.UpdateTenant(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (tenantController *TenantController) DeleteTenant(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger tenant deletion
	if err := tenantController.tenantService.DeleteTenant(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
