package scheduler

import (
	"context"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/pb"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Scheduler) checkEventsToNotify(ctx context.Context) error {
	conn, err := grpc.NewClient(s.calendarAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewEventsServiceClient(conn)

	req := &pb.GetEventsForNotificationRequest{
		Startdate: timestamppb.New(time.Now().Add(-s.timePeriod)),
		Enddate:   timestamppb.New(time.Now()),
	}
	var events *pb.GetEventsResponse
	if events, err = client.GetEventsForNotification(ctx, req); err != nil {
		return err
	}
	for _, pbEvent := range events.Events {
		s.notifyQueue <- Notification{
			ID:            storage.EventID(pbEvent.GetId()),
			OwnerID:       storage.EventOwnerID(pbEvent.GetOwnerId()),
			Title:         pbEvent.GetTitle(),
			StartDateTime: pbEvent.GetStart().AsTime(),
		}
	}

	return nil
}

func (s *Scheduler) checkOldEvents(ctx context.Context) error {
	conn, err := grpc.NewClient(s.calendarAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewEventsServiceClient(conn)
	req := &pb.RemoveOldEventsRequest{
		Date: timestamppb.New(time.Now().Add(-s.eventsExpiration)),
	}
	if _, err = client.RemoveOldEvents(ctx, req); err != nil {
		return err
	}

	return nil
}
