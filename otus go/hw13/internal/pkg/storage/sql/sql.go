package sql

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/types"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/config"
	"github.com/lenniDespero/otus-golang/hw13/internal/pkg/models"
)

// Storage struct
type Storage struct {
	ConnPool *pgxpool.Pool
	ctx      context.Context
}

//New returns new storage
func New(dbconf *config.DBConfig) (*Storage, error) {
	storage := &Storage{}
	storage.ctx = context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Database)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection configs")
	}
	cfg.MaxConns = 8
	cfg.ConnConfig.TLSConfig = nil
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: 5 * time.Minute,
		Timeout:   1 * time.Second,
	}).DialContext

	pool, err := pgxpool.ConnectConfig(storage.ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	storage.ConnPool = pool
	return storage, nil
}

func (s Storage) Add(event models.Event, ctx context.Context) (string, error) {
	err := s.isInTime(event.DateStarted, event.DateComplete, ctx)
	if err != nil {
		return "", err
	}
	sql := `INSERT INTO calendar.event (title,notice,date_start,date_complete,date_created) VALUES($1,$2,$3,$4,$5) RETURNING id;`
	rows, err := s.ConnPool.Query(ctx, sql, event.Title, event.Notice, event.DateStarted, event.DateComplete, time.Now())
	if err != nil {
		return "", errors.Wrap(err, "failed to update event")
	}
	defer rows.Close()
	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if rows.Err() != nil {
		if rows.Err().Error() == `ERROR: duplicate key value violates unique constraint "event_id_idx" (SQLSTATE 23505)` {
			return "", types.ErrEventIdExists
		}
		return "", rows.Err()
	}
	return result, nil

}

func (s Storage) Edit(id string, event models.Event, ctx context.Context) error {
	sql := `Update calendar.event 
			SET id = $1,
				title = $2,
				notice = $3,
				date_start = $4,
				date_complete = $5,
				date_edited = $6
			WHERE id = $7;`
	rows, err := s.ConnPool.Query(ctx, sql, event.ID, event.Title, event.Notice, event.DateStarted, event.DateComplete, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, "failed to update event")
	}
	defer rows.Close()
	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) GetEvents(ctx context.Context) ([]models.Event, error) {
	sql := "SELECT * FROM calendar.event WHERE deleted = $1;"
	result, err := s.ConnPool.Query(ctx, sql, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to receive events")
	}

	defer result.Close()
	events := make([]models.Event, 0)
	for result.Next() {
		var id, title, notice sql2.NullString
		var dateStart, dateComplete, dateEdited, dateCreated sql2.NullTime
		var deleted bool
		var creatorId, editorId sql2.NullInt64
		if err := result.Scan(&id, &title, &notice, &deleted, &dateStart, &dateComplete, &creatorId, &editorId, &dateCreated, &dateEdited); err != nil {
			return nil, errors.Wrap(err, "failed to scan result into vars")
		}
		ev := models.Event{
			ID:           id.String,
			Title:        title.String,
			Notice:       notice.String,
			Deleted:      deleted,
			DateStarted:  dateStart.Time,
			DateComplete: dateComplete.Time,
		}
		if dateCreated.Valid {
			ev.DateCreated = dateCreated.Time
		}
		if dateEdited.Valid {
			ev.DateEdited = dateEdited.Time
		}
		if editorId.Valid {
			ev.EditorID = editorId.Int64
		}
		if creatorId.Valid {
			ev.CreatorID = creatorId.Int64
		}
		events = append(events, ev)
	}

	if err := result.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to load result")
	}
	if len(events) > 0 {
		return events, nil
	}
	return nil, types.ErrNotFound
}

func (s Storage) GetEventByID(id string, ctx context.Context) ([]models.Event, error) {
	sql := "SELECT * FROM calendar.event WHERE deleted = $1 and id = $2;"
	result, err := s.ConnPool.Query(ctx, sql, false, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to receive events")
	}

	defer result.Close()
	events := make([]models.Event, 0)
	for result.Next() {
		var id, title, notice sql2.NullString
		var dateStart, dateComplete, dateEdited, dateCreated sql2.NullTime
		var deleted bool
		var creatorId, editorId sql2.NullInt64
		if err := result.Scan(&id, &title, &notice, &deleted, &dateStart, &dateComplete, &creatorId, &editorId, &dateCreated, &dateEdited); err != nil {
			return nil, errors.Wrap(err, "failed to scan result into vars")
		}
		ev := models.Event{
			ID:           id.String,
			Title:        title.String,
			Notice:       notice.String,
			Deleted:      deleted,
			DateStarted:  dateStart.Time,
			DateComplete: dateComplete.Time,
		}
		if dateCreated.Valid {
			ev.DateCreated = dateCreated.Time
		}
		if dateEdited.Valid {
			ev.DateEdited = dateEdited.Time
		}
		if editorId.Valid {
			ev.EditorID = editorId.Int64
		}
		if creatorId.Valid {
			ev.CreatorID = creatorId.Int64
		}
		events = append(events, ev)
	}

	if err := result.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to load result")
	}
	if len(events) > 0 {
		return events, nil
	}
	return nil, types.ErrNotFound
}

func (s Storage) Delete(id string, ctx context.Context) error {
	sql := `UPDATE calendar.event SET deleted = $1 where id = $2;`
	_, err := s.ConnPool.Exec(ctx, sql, true, id)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) GetEventsByStartPeriod(timeBefore string, timeLength string, ctx context.Context) ([]models.Event, error) {
	timeBeforeInt, err := strconv.Atoi(timeBefore)
	if err != nil {
		return nil, err
	}
	timeLengthInt, err := strconv.Atoi(timeLength)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	dateBefore := now.Local().Add(time.Minute * time.Duration(timeBeforeInt))
	dateAfter := now.Local().Add(time.Minute * time.Duration(timeBeforeInt+timeLengthInt))
	sql := `
		SELECT * FROM calendar.event WHERE deleted = $1
		AND date_start between $2 and $3;`
	rows, err := s.ConnPool.Query(ctx, sql, false, dateBefore, dateAfter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]models.Event, 0)
	for rows.Next() {
		var id, title, notice sql2.NullString
		var dateStart, dateComplete, dateEdited, dateCreated sql2.NullTime
		var deleted bool
		var creatorId, editorId sql2.NullInt64
		if err := rows.Scan(&id, &title, &notice, &deleted, &dateStart, &dateComplete, &creatorId, &editorId, &dateCreated, &dateEdited); err != nil {
			return nil, errors.Wrap(err, "failed to scan result into vars")
		}
		ev := models.Event{
			ID:           id.String,
			Title:        title.String,
			Notice:       notice.String,
			Deleted:      deleted,
			DateStarted:  dateStart.Time,
			DateComplete: dateComplete.Time,
		}
		if dateCreated.Valid {
			ev.DateCreated = dateCreated.Time
		}
		if dateEdited.Valid {
			ev.DateEdited = dateEdited.Time
		}
		if editorId.Valid {
			ev.EditorID = editorId.Int64
		}
		if creatorId.Valid {
			ev.CreatorID = creatorId.Int64
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to load result")
	}
	return events, nil

}

func (s Storage) isInTime(dateStarted time.Time, dateComplete time.Time, ctx context.Context) error {
	sql := `
		SELECT COUNT(*) from calendar.event where (
			($1 between date_start and date_complete)
			or ($2 between date_start and date_complete)
		);
	`
	rows, err := s.ConnPool.Query(ctx, sql, dateStarted, dateComplete)
	if err != nil {
		return err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return err
		}
	}
	if count != 0 {
		return types.ErrDateBusy
	}
	return nil
}

func (s Storage) ClearDB() error {
	sql := `delete from calendar.event;`
	_, err := s.ConnPool.Exec(s.ctx, sql)
	if err != nil {
		return err
	}
	return nil
}
