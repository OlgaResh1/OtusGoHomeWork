//go:generate protoc --proto_path=../../../api/ --go_out=../../../internal/pb --go-grpc_out=../../../internal/pb ../../../api/EventService.proto

package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/pb"      //nolint:depguard
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/golang/protobuf/ptypes/empty"                                     //nolint:depguard
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedEventsServiceServer
	logger Logger
	app    Application
	server *grpc.Server
	addr   string
}

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (storage.EventID, error)
	UpdateEvent(ctx context.Context, id storage.EventID, event storage.Event) error
	RemoveEvent(ctx context.Context, id storage.EventID) error
	GetEventsAll(ctx context.Context, ownerID storage.EventOwnerID) ([]storage.Event, error)
	GetEventsForDay(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForWeek(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForMonth(ctx context.Context, ownerID storage.EventOwnerID, date time.Time) ([]storage.Event, error)
	GetEventsForNotification(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error)
	RemoveOldEvents(ctx context.Context, date time.Time) error
}

func fromPBEvent(event *pb.Event) storage.Event {
	return storage.Event{
		OwnerID:       storage.EventOwnerID(event.GetOwnerId()),
		Title:         event.GetTitle(),
		Description:   event.GetDescription(),
		StartDateTime: event.GetStart().AsTime(),
		Duration:      event.GetDuration().AsDuration(),
		TimeToNotify:  event.GetTimeToNotify().AsDuration(),
	}
}

func toPBEvents(events []storage.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:           int64(event.ID),
			OwnerId:      int64(event.OwnerID),
			Title:        event.Title,
			Description:  event.Description,
			Start:        timestamppb.New(event.StartDateTime),
			Duration:     durationpb.New(event.Duration),
			TimeToNotify: durationpb.New(event.TimeToNotify),
		}
	}
	return pbEvents
}

func (s *Server) CreateEvent(ctx context.Context, event *pb.Event) (*pb.CreateEventResponse, error) {
	if event == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	eventID, err := s.app.CreateEvent(ctx, fromPBEvent(event))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateEventResponse{Id: int64(eventID)}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, event *pb.Event) (*pb.UpdateEventResponse, error) {
	if event == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}
	err := s.app.UpdateEvent(ctx, storage.EventID(event.Id), fromPBEvent(event))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateEventResponse{Id: event.Id}, nil
}

func (s *Server) RemoveEvent(ctx context.Context, req *pb.RemoveEventRequest) (*empty.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "event_id is not specified")
	}
	err := s.app.RemoveEvent(ctx, storage.EventID(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *Server) RemoveOldEvents(ctx context.Context, req *pb.RemoveOldEventsRequest) (*empty.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "date is not specified")
	}
	err := s.app.RemoveOldEvents(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *Server) GetEventsAll(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is not specified")
	}
	events, err := s.app.GetEventsAll(ctx, storage.EventOwnerID(req.GetOwnerId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventsResponse{Events: toPBEvents(events)}, nil
}

func (s *Server) GetEventsForDay(ctx context.Context, req *pb.GetEventsIntervalRequest) (*pb.GetEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is not specified")
	}
	events, err := s.app.GetEventsForDay(ctx, storage.EventOwnerID(req.GetOwnerId()), req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventsResponse{Events: toPBEvents(events)}, nil
}

func (s *Server) GetEventsForWeek(ctx context.Context,
	req *pb.GetEventsIntervalRequest,
) (*pb.GetEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is not specified")
	}
	events, err := s.app.GetEventsForWeek(ctx, storage.EventOwnerID(req.GetOwnerId()), req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventsResponse{Events: toPBEvents(events)}, nil
}

func (s *Server) GetEventsForMonth(ctx context.Context,
	req *pb.GetEventsIntervalRequest,
) (*pb.GetEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is not specified")
	}
	events, err := s.app.GetEventsForMonth(ctx, storage.EventOwnerID(req.GetOwnerId()), req.GetDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventsResponse{Events: toPBEvents(events)}, nil
}

func (s *Server) GetEventsForNotification(ctx context.Context,
	req *pb.GetEventsForNotificationRequest,
) (*pb.GetEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "owner_id is not specified")
	}
	events, err := s.app.GetEventsForNotification(ctx, req.GetStartdate().AsTime(), req.GetEnddate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventsResponse{Events: toPBEvents(events)}, nil
}

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	return &Server{logger: logger, app: app, addr: cfg.GRPC.Address}
}

func (s *Server) Start(_ context.Context) error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.RequestLogInterceptor,
		),
	)
	s.server = server

	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("GRPC Server started %s", s.addr))
	pb.RegisterEventsServiceServer(s.server, s)
	return s.server.Serve(lsn)
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
