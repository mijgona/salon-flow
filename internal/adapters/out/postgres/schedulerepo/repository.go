package schedulerepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresScheduleRepository is a pgx + Squirrel implementation of MasterScheduleRepository.
type PostgresScheduleRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresScheduleRepository creates a new repository.
func NewPostgresScheduleRepository(pool *pgxpool.Pool) *PostgresScheduleRepository {
	return &PostgresScheduleRepository{pool: pool}
}

func (r *PostgresScheduleRepository) Add(ctx context.Context, tx interface{}, ms *scheduling.MasterSchedule) error {
	bookedJSON, _ := json.Marshal(slotsToJSON(ms.BookedSlots()))
	blockedJSON, _ := json.Marshal(slotsToJSON(ms.BlockedSlots()))

	query, args, err := psql.Insert("master_schedules").
		Columns(
			"id", "master_id", "salon_id", "schedule_date",
			"work_start", "work_end", "break_start", "break_end",
			"booked_slots", "blocked_slots",
		).
		Values(
			ms.ID(), ms.MasterID(), ms.SalonID(), ms.Date(),
			ms.WorkingHours().StartTime().Format("15:04"),
			ms.WorkingHours().EndTime().Format("15:04"),
			nullableTime(ms.WorkingHours().BreakStart()),
			nullableTime(ms.WorkingHours().BreakEnd()),
			bookedJSON, blockedJSON,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresScheduleRepository) Update(ctx context.Context, tx interface{}, ms *scheduling.MasterSchedule) error {
	bookedJSON, _ := json.Marshal(slotsToJSON(ms.BookedSlots()))
	blockedJSON, _ := json.Marshal(slotsToJSON(ms.BlockedSlots()))

	query, args, err := psql.Update("master_schedules").
		Set("booked_slots", bookedJSON).
		Set("blocked_slots", blockedJSON).
		Set("work_start", ms.WorkingHours().StartTime().Format("15:04")).
		Set("work_end", ms.WorkingHours().EndTime().Format("15:04")).
		Set("break_start", nullableTime(ms.WorkingHours().BreakStart())).
		Set("break_end", nullableTime(ms.WorkingHours().BreakEnd())).
		Where(sq.Eq{"id": ms.ID()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresScheduleRepository) GetByMasterAndDate(ctx context.Context, tx interface{}, masterID uuid.UUID, date time.Time) (*scheduling.MasterSchedule, error) {
	scheduleDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	query, args, err := psql.Select(
		"id", "master_id", "salon_id", "schedule_date",
		"work_start", "work_end", "break_start", "break_end",
		"booked_slots", "blocked_slots",
	).From("master_schedules").
		Where(sq.Eq{"master_id": masterID, "schedule_date": scheduleDate}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	rows, err := r.query(ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	return scanMasterScheduleRow(rows, date.Location())
}

// --- internal helpers ---

func (r *PostgresScheduleRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *PostgresScheduleRepository) query(ctx context.Context, tx interface{}, sql string, args ...any) (pgx.Rows, error) {
	if t := postgres.ExtractTx(tx); t != nil {
		return t.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func scanMasterScheduleRow(rows pgx.Rows, loc *time.Location) (*scheduling.MasterSchedule, error) {
	var (
		id           uuid.UUID
		masterID     uuid.UUID
		salonID      uuid.UUID
		scheduleDate time.Time
		workStart    string
		workEnd      string
		breakStart   *string
		breakEnd     *string
		bookedJSON   []byte
		blockedJSON  []byte
	)

	err := rows.Scan(
		&id, &masterID, &salonID, &scheduleDate,
		&workStart, &workEnd, &breakStart, &breakEnd,
		&bookedJSON, &blockedJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("scan master_schedule: %w", err)
	}

	ws := parseTimeOfDay(workStart, scheduleDate, loc)
	we := parseTimeOfDay(workEnd, scheduleDate, loc)
	var bs, be time.Time
	if breakStart != nil && breakEnd != nil {
		bs = parseTimeOfDay(*breakStart, scheduleDate, loc)
		be = parseTimeOfDay(*breakEnd, scheduleDate, loc)
	}

	workingHours := scheduling.MustNewWorkingHours(ws, we, bs, be)
	bookedSlots := parseSlotsJSON(bookedJSON, scheduleDate, loc)
	blockedSlots := parseSlotsJSON(blockedJSON, scheduleDate, loc)

	return scheduling.RestoreMasterSchedule(
		id, masterID, salonID, scheduleDate,
		workingHours, bookedSlots, blockedSlots,
	), nil
}

// --- JSON slot helpers ---

type slotJSON struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func slotsToJSON(slots []scheduling.TimeSlot) []slotJSON {
	result := make([]slotJSON, 0, len(slots))
	for _, s := range slots {
		result = append(result, slotJSON{
			Start: s.StartTime().Format(time.RFC3339),
			End:   s.EndTime().Format(time.RFC3339),
		})
	}
	return result
}

func parseSlotsJSON(data []byte, date time.Time, loc *time.Location) []scheduling.TimeSlot {
	if len(data) == 0 {
		return nil
	}
	var items []slotJSON
	if err := json.Unmarshal(data, &items); err != nil {
		return nil
	}
	result := make([]scheduling.TimeSlot, 0, len(items))
	for _, item := range items {
		start, err1 := time.Parse(time.RFC3339, item.Start)
		end, err2 := time.Parse(time.RFC3339, item.End)
		if err1 == nil && err2 == nil {
			result = append(result, scheduling.MustNewTimeSlot(start, end))
		}
	}
	return result
}

func parseTimeOfDay(s string, date time.Time, loc *time.Location) time.Time {
	t, _ := time.Parse("15:04", s)
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, loc)
}

func nullableTime(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format("15:04")
	return &s
}
