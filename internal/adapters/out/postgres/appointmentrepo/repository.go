package appointmentrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresAppointmentRepository is a pgx + Squirrel implementation of AppointmentRepository.
type PostgresAppointmentRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresAppointmentRepository creates a new repository.
func NewPostgresAppointmentRepository(pool *pgxpool.Pool) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{pool: pool}
}

func (r *PostgresAppointmentRepository) Add(ctx context.Context, tx interface{}, a *scheduling.Appointment) error {
	query, args, err := psql.Insert("appointments").
		Columns(
			"id", "tenant_id", "client_id", "master_id", "salon_id",
			"service_id", "service_name", "service_duration",
			"start_time", "end_time", "status",
			"price_amount", "price_currency", "source", "comment", "created_at",
		).
		Values(
			a.ID(), a.TenantID().UUID(), a.ClientID(), a.MasterID(), a.SalonID(),
			a.ServiceInfo().ServiceID(), a.ServiceInfo().Name(),
			fmt.Sprintf("%d minutes", int(a.ServiceInfo().Duration().Minutes())),
			a.TimeSlot().StartTime(), a.TimeSlot().EndTime(),
			a.Status().String(),
			a.Price().Amount(), a.Price().Currency(),
			a.Source().String(), a.Comment(), a.CreatedAt(),
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) Update(ctx context.Context, tx interface{}, a *scheduling.Appointment) error {
	query, args, err := psql.Update("appointments").
		Set("status", a.Status().String()).
		Set("start_time", a.TimeSlot().StartTime()).
		Set("end_time", a.TimeSlot().EndTime()).
		Set("price_amount", a.Price().Amount()).
		Set("comment", a.Comment()).
		Where(sq.Eq{"id": a.ID()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) Get(ctx context.Context, tx interface{}, id uuid.UUID) (*scheduling.Appointment, error) {
	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointment(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) FindByClientID(ctx context.Context, tx interface{}, clientID uuid.UUID) ([]*scheduling.Appointment, error) {
	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"client_id": clientID}).
		OrderBy("start_time DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointments(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) FindByMasterAndDate(ctx context.Context, tx interface{}, masterID uuid.UUID, date time.Time) ([]*scheduling.Appointment, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"master_id": masterID}).
		Where(sq.GtOrEq{"start_time": dayStart}).
		Where(sq.Lt{"start_time": dayEnd}).
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointments(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) FindByDateRange(ctx context.Context, tx interface{}, tenantID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"tenant_id": tenantID}).
		Where(sq.GtOrEq{"start_time": from}).
		Where(sq.Lt{"start_time": to}).
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointments(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) FindByMasterDateRange(ctx context.Context, tx interface{}, masterID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"master_id": masterID}).
		Where(sq.GtOrEq{"start_time": from}).
		Where(sq.Lt{"start_time": to}).
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointments(ctx, tx, query, args...)
}

func (r *PostgresAppointmentRepository) FindBySalonDateRange(ctx context.Context, tx interface{}, salonID uuid.UUID, from, to time.Time) ([]*scheduling.Appointment, error) {
	query, args, err := psql.Select(appointmentColumns()...).
		From("appointments").
		Where(sq.Eq{"salon_id": salonID}).
		Where(sq.GtOrEq{"start_time": from}).
		Where(sq.Lt{"start_time": to}).
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanAppointments(ctx, tx, query, args...)
}

// --- internal helpers ---

func appointmentColumns() []string {
	return []string{
		"id", "tenant_id", "client_id", "master_id", "salon_id",
		"service_id", "service_name", "service_duration",
		"start_time", "end_time", "status",
		"price_amount", "price_currency", "source", "comment", "created_at",
	}
}

func (r *PostgresAppointmentRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *PostgresAppointmentRepository) query(ctx context.Context, tx interface{}, sql string, args ...any) (pgx.Rows, error) {
	if t := postgres.ExtractTx(tx); t != nil {
		return t.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func (r *PostgresAppointmentRepository) scanAppointment(ctx context.Context, tx interface{}, sql string, args ...any) (*scheduling.Appointment, error) {
	rows, err := r.query(ctx, tx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	return scanAppointmentRow(rows)
}

func (r *PostgresAppointmentRepository) scanAppointments(ctx context.Context, tx interface{}, sql string, args ...any) ([]*scheduling.Appointment, error) {
	rows, err := r.query(ctx, tx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*scheduling.Appointment
	for rows.Next() {
		a, err := scanAppointmentRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

func scanAppointmentRow(rows pgx.Rows) (*scheduling.Appointment, error) {
	var (
		id              uuid.UUID
		tenantID        uuid.UUID
		clientID        uuid.UUID
		masterID        uuid.UUID
		salonID         uuid.UUID
		serviceID       uuid.UUID
		serviceName     string
		serviceDuration string // INTERVAL comes as string
		startTime       time.Time
		endTime         time.Time
		status          string
		priceAmount     decimal.Decimal
		priceCurrency   string
		source          string
		comment         *string
		createdAt       time.Time
	)

	err := rows.Scan(
		&id, &tenantID, &clientID, &masterID, &salonID,
		&serviceID, &serviceName, &serviceDuration,
		&startTime, &endTime, &status,
		&priceAmount, &priceCurrency, &source, &comment, &createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan appointment: %w", err)
	}

	tid := model.MustNewTenantID(tenantID)
	price := model.MustNewMoney(priceAmount, priceCurrency)
	duration := endTime.Sub(startTime)
	serviceInfo := scheduling.MustNewServiceInfo(serviceID, serviceName, duration, price)
	timeSlot := scheduling.MustNewTimeSlot(startTime, endTime)

	var commentStr string
	if comment != nil {
		commentStr = *comment
	}

	return scheduling.RestoreAppointment(
		id, tid, clientID, masterID, salonID,
		serviceInfo, timeSlot,
		scheduling.AppointmentStatus(status),
		price,
		scheduling.BookingSource(source),
		commentStr,
		createdAt,
	), nil
}
