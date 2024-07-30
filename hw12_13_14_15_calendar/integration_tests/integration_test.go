//go:build integration
// +build integration

package integration_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/pb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TestData struct {
	eventId1, eventId2, eventId3 int64
	ownerId                      int64
}
type CalendarSuite struct {
	suite.Suite
	ctx         context.Context
	serviceConn *grpc.ClientConn
	client      pb.EventsServiceClient
	testData    TestData
}

func (s *CalendarSuite) SetupSuite() {
	calendarHost := os.Getenv("CALENDAR_SERVER_HOST")
	if calendarHost == "" {
		calendarHost = "127.0.0.1:50051"
	}
	conn, err := grpc.NewClient(calendarHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().Nil(err)

	s.ctx = context.Background()
	s.serviceConn = conn
	s.client = pb.NewEventsServiceClient(s.serviceConn)
}

func (s *CalendarSuite) SetupTest() {
	var seed int64 = time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	s.T().Log("seed:", seed)
	s.testData.ownerId = rand.Int63()
	s.T().Log("ownerId:", s.testData.ownerId)

	s.createEvents()
}

func (s *CalendarSuite) getEventPB(eventTimeString string) *pb.Event {
	eventTime, err := time.Parse("02.01.2006 15:04:05", eventTimeString)
	if err != nil {
		return nil
	}
	return &pb.Event{
		OwnerId:      s.testData.ownerId,
		Title:        "test event",
		Description:  "test event description",
		Start:        timestamppb.New(eventTime),
		Duration:     durationpb.New(0),
		TimeToNotify: durationpb.New(0),
	}
}

func (s *CalendarSuite) getGetEventsRequest(dateTime string) *pb.GetEventsIntervalRequest {
	datetime, err := time.Parse("02.01.2006 15:04:05", dateTime)
	s.Require().NoError(err)
	return &pb.GetEventsIntervalRequest{
		OwnerId: s.testData.ownerId,
		Date:    timestamppb.New(datetime),
	}
}

func (s *CalendarSuite) createEvents() {
	eventReq1 := s.getEventPB("13.05.2024 12:00:00")
	s.Require().NotNil(eventReq1)
	ev1Id, err := s.client.CreateEvent(s.ctx, eventReq1)
	s.Require().NoError(err)
	s.Require().NotNil(ev1Id)
	s.testData.eventId1 = ev1Id.GetId()

	eventReq2 := s.getEventPB("16.05.2024 20:00:00")
	ev2Id, err := s.client.CreateEvent(s.ctx, eventReq2)
	s.Require().NoError(err)
	s.Require().NotNil(ev2Id)
	s.testData.eventId2 = ev2Id.GetId()

	eventReq3 := s.getEventPB("23.05.2024 23:00:00")
	ev3Id, err := s.client.CreateEvent(s.ctx, eventReq3)
	s.Require().NoError(err)
	s.Require().NotNil(ev3Id)
	s.testData.eventId3 = ev3Id.GetId()
}

func (s *CalendarSuite) TestCreateEvents() {
	eventReq1 := s.getEventPB("28.07.2024 12:00:00")
	s.Require().NotNil(eventReq1)
	ev1Id, err := s.client.CreateEvent(s.ctx, eventReq1)
	s.Require().NoError(err)
	s.Require().NotNil(ev1Id)
	s.testData.eventId1 = ev1Id.GetId()

	eventReq2 := s.getEventPB("29.07.2024 20:00:00")
	ev2Id, err := s.client.CreateEvent(s.ctx, eventReq2)
	s.Require().NoError(err)
	s.Require().NotNil(ev2Id)
	s.testData.eventId1 = ev2Id.GetId()

	eventReq3 := s.getEventPB("30.07.2024 23:00:00")
	ev3Id, err := s.client.CreateEvent(s.ctx, eventReq3)
	s.Require().NoError(err)
	s.Require().NotNil(ev3Id)
	s.testData.eventId1 = ev3Id.GetId()

	s.T().Log("test create events ok")
}

func (s *CalendarSuite) TestSelectEvents() {
	testDate := "13.05.2024 10:00:00"

	resp, err := s.client.GetEventsForDay(s.ctx, s.getGetEventsRequest(testDate))
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Events, 1)
	s.Require().Equal(s.testData.eventId1, resp.Events[0].Id)
	expTime, err := time.Parse("02.01.2006 15:04:05", "13.05.2024 12:00:00")
	s.Require().NoError(err)
	s.Require().Equal(expTime, resp.Events[0].GetStart().AsTime())

	resp2, err := s.client.GetEventsForWeek(s.ctx, s.getGetEventsRequest(testDate))
	s.Require().NoError(err)
	s.Require().NotNil(resp2)
	s.Require().Len(resp2.Events, 2)

	resp3, err := s.client.GetEventsForMonth(s.ctx, s.getGetEventsRequest(testDate))
	s.Require().NoError(err)
	s.Require().NotNil(resp3)
	s.Require().Len(resp3.Events, 3)

	resp4, err := s.client.GetEventsAll(s.ctx, &pb.GetEventsRequest{OwnerId: s.testData.ownerId})
	s.Require().NoError(err)
	s.Require().NotNil(resp4)
	s.Require().Len(resp4.Events, 3)

	s.T().Log("test select events ok")
}

func (s *CalendarSuite) checkCountForWeek(date string, expectedLen int) {
	resp, err := s.client.GetEventsForWeek(s.ctx, s.getGetEventsRequest(date))
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.Events, expectedLen)
}

func (s *CalendarSuite) TestChangeEvents() {
	testDate := "26.05.2024 12:00:00"

	s.checkCountForWeek(testDate, 1)

	req := s.getEventPB("26.05.2024 12:00:00")
	s.Require().NotNil(req)
	req.Id = s.testData.eventId1

	resp1, err := s.client.UpdateEvent(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp1)
	s.Require().Equal(s.testData.eventId1, resp1.Id)

	s.checkCountForWeek(testDate, 2)

	_, err = s.client.RemoveEvent(s.ctx, &pb.RemoveEventRequest{Id: s.testData.eventId1})
	s.Require().NoError(err)

	s.checkCountForWeek(testDate, 1)

	s.T().Log("test change events ok")
}

func (s *CalendarSuite) TearDownTest() {
	_, err := s.client.RemoveEvent(s.ctx, &pb.RemoveEventRequest{Id: s.testData.eventId1})
	s.Require().NoError(err)
	_, err = s.client.RemoveEvent(s.ctx, &pb.RemoveEventRequest{Id: s.testData.eventId2})
	s.Require().NoError(err)
	_, err = s.client.RemoveEvent(s.ctx, &pb.RemoveEventRequest{Id: s.testData.eventId3})
	s.Require().NoError(err)
}

func (s *CalendarSuite) TearDownSuite() {
	s.serviceConn.Close()
}

func TestCalendar(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
