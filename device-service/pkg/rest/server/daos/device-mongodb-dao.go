package daos

import (
	"context"
	"errors"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/rest/server/daos/clients/nosqls"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeviceDao struct {
	mongoDBClient *nosqls.MongoDBClient
	collection    *mongo.Collection
}

func NewDeviceDao() (*DeviceDao, error) {
	mongoDBClient, err := nosqls.InitMongoDB()
	if err != nil {
		log.Debugf("init mongoDB failed: %v", err)
		return nil, err
	}
	return &DeviceDao{
		mongoDBClient: mongoDBClient,
		collection:    mongoDBClient.Database.Collection("devices"),
	}, nil
}

func (deviceDao *DeviceDao) CreateDevice(device *models.Device) (*models.Device, error) {
	// create a document for given device
	insertOneResult, err := deviceDao.collection.InsertOne(context.TODO(), device)
	if err != nil {
		log.Debugf("insert failed: %v", err)
		return nil, err
	}
	device.ID = insertOneResult.InsertedID.(primitive.ObjectID).Hex()

	log.Debugf("device created")
	return device, nil
}

func (deviceDao *DeviceDao) ListDevices() ([]*models.Device, error) {
	filters := bson.D{}
	devices, err := deviceDao.collection.Find(context.TODO(), filters)
	if err != nil {
		return nil, err
	}
	var deviceList []*models.Device
	for devices.Next(context.TODO()) {
		var device *models.Device
		if err = devices.Decode(&device); err != nil {
			log.Debugf("decode device failed: %v", err)
			return nil, err
		}
		deviceList = append(deviceList, device)
	}
	if deviceList == nil {
		return []*models.Device{}, nil
	}

	log.Debugf("device listed")
	return deviceList, nil
}

func (deviceDao *DeviceDao) GetDevice(id string) (*models.Device, error) {
	var device *models.Device
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Device{}, nosqls.ErrInvalidObjectID
	}
	filter := bson.D{
		{Key: "_id", Value: objectID},
	}
	if err = deviceDao.collection.FindOne(context.TODO(), filter).Decode(&device); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			return &models.Device{}, nosqls.ErrNotExists
		}
		log.Debugf("decode device failed: %v", err)
		return nil, err
	}

	log.Debugf("device retrieved")
	return device, nil
}
