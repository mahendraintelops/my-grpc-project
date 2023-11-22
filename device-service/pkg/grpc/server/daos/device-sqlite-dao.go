package daos

import (
	"database/sql"
	"errors"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/grpc/server/daos/clients/sqls"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/grpc/server/models"
	log "github.com/sirupsen/logrus"
)

type DeviceDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateDevices(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS devices(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Name TEXT NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewDeviceDao() (*DeviceDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateDevices(sqlClient)
	if err != nil {
		return nil, err
	}
	return &DeviceDao{
		sqlClient,
	}, nil
}

func (deviceDao *DeviceDao) CreateDevice(m *models.Device) (*models.Device, error) {
	insertQuery := "INSERT INTO devices(Name)values(?)"
	res, err := deviceDao.sqlClient.DB.Exec(insertQuery, m.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("device created")
	return m, nil
}

func (deviceDao *DeviceDao) ListDevices() ([]*models.Device, error) {
	selectQuery := "SELECT * FROM devices"
	rows, err := deviceDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var devices []*models.Device
	for rows.Next() {
		m := models.Device{}
		if err = rows.Scan(&m.Id, &m.Name); err != nil {
			return nil, err
		}
		devices = append(devices, &m)
	}
	if devices == nil {
		devices = []*models.Device{}
	}

	log.Debugf("device listed")
	return devices, nil
}

func (deviceDao *DeviceDao) GetDevice(id int64) (*models.Device, error) {
	selectQuery := "SELECT * FROM devices WHERE Id = ?"
	row := deviceDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Device{}
	if err := row.Scan(&m.Id, &m.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("device retrieved")
	return &m, nil
}

func (deviceDao *DeviceDao) UpdateDevice(id int64, m *models.Device) (*models.Device, error) {
	if id == 0 {
		return nil, errors.New("invalid device ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	device, err := deviceDao.GetDevice(id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE devices SET Name = ? WHERE Id = ?"
	res, err := deviceDao.sqlClient.DB.Exec(updateQuery, m.Name, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("device updated")
	return m, nil
}
