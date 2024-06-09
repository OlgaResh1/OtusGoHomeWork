package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct { // TODO
	db  *sql.DB
	dsn string
}

func New(ctx context.Context, cfg config.Config) (storage *Storage, err error) {
	if len(cfg.Sql.Dsn) == 0 {
		return nil, fmt.Errorf("dsn not defined")
	}
	storage = &Storage{dsn: cfg.Sql.Dsn}

	storage.db, err = sql.Open("pgx", storage.dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to load driver: %w", err)
	}

	err = storage.db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return storage, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.EventId, error) {
	if len(event.Title) == 0 || event.StartDateTime.IsZero() {
		return 0, storage.ErrNotValidEvent
	}
	tx, err := s.db.Begin()
	if err == nil {
		return 0, err
	}
	defer tx.Rollback()

	// check exist event
	query := `select id	from events where owner_id=$1 and begin_datetime ='$2'`
	var existId storage.EventId
	err = s.db.QueryRowContext(ctx, query, event.OwnerId, event.StartDateTime).Scan(&existId)
	if err == nil {
		return 0, storage.ErrDateBusy
	}

	query = `insert into events(owner_id, title, description, begin_datetime, duration)
	values($1, $2, $3, $4, $5) RETURNING id`

	var id storage.EventId
	err = s.db.QueryRowContext(ctx, query, event.OwnerId, event.Title,
		event.Description, event.StartDateTime, event.Duration).Scan(&id)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return storage.EventId(id), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id storage.EventId, event storage.Event) error {
	query := `update events set title = $1, description = $2, 
				begin_datetime = $3, duration = $4, notify=$5 where id = $6`

	result, err := s.db.ExecContext(ctx, query,
		event.Title,
		event.Description,
		event.StartDateTime,
		event.Duration.Abs(),
		event.TimeToNotify.Abs(),
		id,
	)
	if err != nil {
		return err
	}
	idUpdated, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if idUpdated != 1 {
		return fmt.Errorf("update error")
	}
	return err
}

func (s *Storage) RemoveEvent(ctx context.Context, id storage.EventId) error {
	query := `delete from events where id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventsAll(ctx context.Context, ownerId storage.EventOwnerId) ([]storage.Event, error) {
	query := `select id, owner_id, title, description, begin_datetime, duration, notify
		from events where owner_id=$1 order by begin_datetime`

	rows, err := s.db.QueryContext(ctx, query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []storage.Event
	for rows.Next() {
		var event storage.Event
		var duration, notify sql.NullInt64
		var description sql.NullString

		err := rows.Scan(&event.Id, &event.OwnerId, &event.Title, &description, &event.StartDateTime, &duration, &notify)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			event.Description = description.String
		}
		if duration.Valid {
			event.Duration = time.Duration(duration.Int64)
		}
		if notify.Valid {
			event.TimeToNotify = time.Duration(notify.Int64)
		}
		result = append(result, event)
	}
	return result, nil
}

func (s *Storage) getEventsByInterval(ctx context.Context, ownerId storage.EventOwnerId, beginDT time.Time, endDT time.Time) ([]storage.Event, error) {
	query := `select id, owner_id, title, description, begin_datetime, duration, notify
		from events where owner_id=$1 and begin_datetime >=$2 and  begin_datetime <$3 order by begin_datetime;`

	rows, err := s.db.QueryContext(ctx, query, ownerId, beginDT, endDT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []storage.Event
	for rows.Next() {
		var event storage.Event
		var duration, notify sql.NullInt64
		var description sql.NullString

		err := rows.Scan(&event.Id, &event.OwnerId, &event.Title, &description, &event.StartDateTime, &duration, &notify)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			event.Description = description.String
		}
		if duration.Valid {
			event.Duration = time.Duration(duration.Int64)
		}
		if notify.Valid {
			event.TimeToNotify = time.Duration(notify.Int64)
		}
		result = append(result, event)
	}
	return result, nil
}

func (s *Storage) GetEventsForDay(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	startDT := date.Truncate(24 * time.Hour)
	nextDay := date.AddDate(0, 0, 1)

	return s.getEventsByInterval(ctx, ownerId, startDT, nextDay)
}

func (s *Storage) GetEventsForWeek(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	weekday := int(date.Weekday())
	sunday := date.AddDate(0, 0, -weekday)
	nextSunday := date.AddDate(0, 0, 7)

	return s.getEventsByInterval(ctx, ownerId, sunday, nextSunday)
}

func (s *Storage) GetEventsForMonth(ctx context.Context, ownerId storage.EventOwnerId, date time.Time) ([]storage.Event, error) {
	day := int(date.Day())
	firstOfMonth := date.AddDate(0, 0, -day)
	nextMonth := date.AddDate(0, 1, 0)

	return s.getEventsByInterval(ctx, ownerId, firstOfMonth, nextMonth)
}
