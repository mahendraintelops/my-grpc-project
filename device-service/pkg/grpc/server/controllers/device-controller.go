package controllers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
	"strconv"

	pb "github.com/mahendraintelops/my-grpc-project/device-service/gen/api/v1"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/grpc/server/models"
	"github.com/mahendraintelops/my-grpc-project/device-service/pkg/grpc/server/services"
)

type DeviceServer struct {
	deviceService *services.DeviceService
	pb.UnimplementedDeviceServiceServer
}

func NewDeviceServer() (*DeviceServer, error) {
	deviceService, err := services.NewDeviceService()
	if err != nil {
		return nil, err
	}
	return &DeviceServer{
		deviceService: deviceService,
	}, nil
}

func (*DeviceServer) Ping(_ context.Context, _ *pb.DeviceRequest) (*pb.DeviceResponse, error) {
	return &pb.DeviceResponse{
		Message: "Server is healthy and working!",
	}, nil
}

func (deviceServer *DeviceServer) CreateDevice(_ context.Context, req *pb.CreateDeviceRequest) (*pb.CreateDeviceResponse, error) {
	device := models.Device{

		Name: req.Device.GetName(),
	}

	if _, err := deviceServer.deviceService.CreateDevice(&device); err != nil {
		return nil, err
	}

	return &pb.CreateDeviceResponse{
		Message: "Device Created Successfully!",
	}, nil
}

func (deviceServer *DeviceServer) ListDevices(_ context.Context, _ *pb.ListDevicesRequest) (*pb.ListDevicesResponse, error) {
	devices, err := deviceServer.deviceService.ListDevices()

	if err != nil {
		return nil, err
	}

	var deviceList []*pb.Device
	for _, retrievedDevice := range devices {
		deviceResponse := &pb.Device{
			Id: retrievedDevice.Id,

			Name: retrievedDevice.Name,
		}
		deviceList = append(deviceList, deviceResponse)
	}

	return &pb.ListDevicesResponse{
		Device: deviceList,
	}, nil
}

func (deviceServer *DeviceServer) GetDevice(ctx context.Context, req *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	id := req.GetId()
	retrievedDevice, err := deviceServer.deviceService.GetDevice(id)

	if err != nil {
		return nil, err
	}

	deviceResponse := &pb.Device{
		Id: id,

		Name: retrievedDevice.Name,
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(ctx)
		currentSpan.SetAttributes(attribute.String("device.id", strconv.FormatInt(retrievedDevice.Id, 10)))
	}

	return &pb.GetDeviceResponse{
		Device: deviceResponse,
	}, nil
}

func (deviceServer *DeviceServer) UpdateDevice(_ context.Context, req *pb.UpdateDeviceRequest) (*pb.UpdateDeviceResponse, error) {
	id := req.GetId()

	deviceRequest := models.Device{
		Id: id,

		Name: req.Device.GetName(),
	}
	_, err := deviceServer.deviceService.UpdateDevice(id, &deviceRequest)

	if err != nil {
		return nil, err
	}

	return &pb.UpdateDeviceResponse{
		Message: "Device Updated Successfully!",
	}, nil
}
