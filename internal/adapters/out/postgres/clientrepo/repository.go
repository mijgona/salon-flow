package clientrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// PostgresClientRepository is a pgx + Squirrel implementation of ClientRepository.
type PostgresClientRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresClientRepository creates a new repository.
func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{pool: pool}
}

func (r *PostgresClientRepository) conn(tx interface{}) interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (interface{ RowsAffected() int64 }, error)
} {
	if t := postgres.ExtractTx(tx); t != nil {
		return &txWrapper{t}
	}
	return &poolWrapper{r.pool}
}

func (r *PostgresClientRepository) Add(ctx context.Context, tx interface{}, c *client.Client) error {
	prefsJSON, _ := json.Marshal(prefsToMap(c.Preferences()))
	allergiesJSON, _ := json.Marshal(allergiesToSlice(c.Allergies()))
	notesJSON, _ := json.Marshal(notesToSlice(c.Notes()))

	query, args, err := psql.Insert("clients").
		Columns(
			"id", "tenant_id", "phone", "email", "first_name", "last_name",
			"birthday", "preferences", "allergies", "notes", "source", "registered_at",
		).
		Values(
			c.ID(), c.TenantID().UUID(), c.ContactInfo().Phone().String(), c.ContactInfo().Email(),
			c.ContactInfo().FirstName(), c.ContactInfo().LastName(),
			birthdayToNullable(c.Birthday()), prefsJSON, allergiesJSON, notesJSON,
			c.Source().String(), c.RegisteredAt(),
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresClientRepository) Update(ctx context.Context, tx interface{}, c *client.Client) error {
	prefsJSON, _ := json.Marshal(prefsToMap(c.Preferences()))
	allergiesJSON, _ := json.Marshal(allergiesToSlice(c.Allergies()))
	notesJSON, _ := json.Marshal(notesToSlice(c.Notes()))

	query, args, err := psql.Update("clients").
		Set("phone", c.ContactInfo().Phone().String()).
		Set("email", c.ContactInfo().Email()).
		Set("first_name", c.ContactInfo().FirstName()).
		Set("last_name", c.ContactInfo().LastName()).
		Set("birthday", birthdayToNullable(c.Birthday())).
		Set("preferences", prefsJSON).
		Set("allergies", allergiesJSON).
		Set("notes", notesJSON).
		Where(sq.Eq{"id": c.ID()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.exec(ctx, tx, query, args...)
}

func (r *PostgresClientRepository) Get(ctx context.Context, tx interface{}, id uuid.UUID) (*client.Client, error) {
	query, args, err := psql.Select(
		"id", "tenant_id", "phone", "email", "first_name", "last_name",
		"birthday", "preferences", "allergies", "notes", "source", "registered_at",
	).From("clients").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanClient(ctx, tx, query, args...)
}

func (r *PostgresClientRepository) FindByPhone(ctx context.Context, tx interface{}, tenantID model.TenantID, phone model.PhoneNumber) (*client.Client, error) {
	query, args, err := psql.Select(
		"id", "tenant_id", "phone", "email", "first_name", "last_name",
		"birthday", "preferences", "allergies", "notes", "source", "registered_at",
	).From("clients").
		Where(sq.Eq{"tenant_id": tenantID.UUID(), "phone": phone.String()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	return r.scanClient(ctx, tx, query, args...)
}

func (r *PostgresClientRepository) FindByTenant(ctx context.Context, tx interface{}, tenantID model.TenantID, limit, offset int) ([]*client.Client, error) {
	query, args, err := psql.Select(
		"id", "tenant_id", "phone", "email", "first_name", "last_name",
		"birthday", "preferences", "allergies", "notes", "source", "registered_at",
	).From("clients").
		Where(sq.Eq{"tenant_id": tenantID.UUID()}).
		OrderBy("registered_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	rows, err := r.query(ctx, tx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*client.Client
	for rows.Next() {
		c, err := scanClientRow(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

// --- internal helpers ---

func (r *PostgresClientRepository) exec(ctx context.Context, tx interface{}, sql string, args ...any) error {
	if t := postgres.ExtractTx(tx); t != nil {
		_, err := t.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *PostgresClientRepository) query(ctx context.Context, tx interface{}, sql string, args ...any) (pgx.Rows, error) {
	if t := postgres.ExtractTx(tx); t != nil {
		return t.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func (r *PostgresClientRepository) scanClient(ctx context.Context, tx interface{}, sql string, args ...any) (*client.Client, error) {
	rows, err := r.query(ctx, tx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil // not found
	}
	return scanClientRow(rows)
}

func scanClientRow(rows pgx.Rows) (*client.Client, error) {
	var (
		id            uuid.UUID
		tenantID      uuid.UUID
		phone         string
		email         string
		firstName     string
		lastName      string
		birthday      *time.Time
		prefsJSON     []byte
		allergiesJSON []byte
		notesJSON     []byte
		source        string
		registeredAt  time.Time
	)

	err := rows.Scan(
		&id, &tenantID, &phone, &email, &firstName, &lastName,
		&birthday, &prefsJSON, &allergiesJSON, &notesJSON,
		&source, &registeredAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan client: %w", err)
	}

	tid := model.MustNewTenantID(tenantID)
	ph := model.MustNewPhoneNumber(phone)
	ci := client.MustNewContactInfo(ph, email, firstName, lastName)

	var bd model.Birthday
	if birthday != nil {
		bd, _ = model.NewBirthday(*birthday)
	}

	prefs := parsePreferences(prefsJSON)
	allergies := parseAllergies(allergiesJSON)
	notes := parseNotes(notesJSON)

	return client.RestoreClient(
		id, tid, ci, bd, prefs, allergies, notes,
		nil, nil, // photos, visitRecords — stored in separate tables or loaded lazily
		client.ClientSource(source), registeredAt,
	), nil
}

// --- JSON mapping helpers ---

func prefsToMap(p client.Preferences) map[string]interface{} {
	favs := make([]string, 0, len(p.FavoriteServices()))
	for _, f := range p.FavoriteServices() {
		favs = append(favs, f.String())
	}
	return map[string]interface{}{
		"preferred_master_id": p.PreferredMasterID().String(),
		"favorite_services":   favs,
		"channel":             string(p.Channel()),
	}
}

func parsePreferences(data []byte) client.Preferences {
	if len(data) == 0 {
		return client.EmptyPreferences()
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return client.EmptyPreferences()
	}
	masterID, _ := uuid.Parse(fmt.Sprint(m["preferred_master_id"]))
	var favs []uuid.UUID
	if arr, ok := m["favorite_services"].([]interface{}); ok {
		for _, v := range arr {
			if uid, err := uuid.Parse(fmt.Sprint(v)); err == nil {
				favs = append(favs, uid)
			}
		}
	}
	ch := client.CommunicationChannel(fmt.Sprint(m["channel"]))
	return client.NewPreferences(masterID, favs, ch)
}

type allergyJSON struct {
	Substance string `json:"substance"`
	Severity  string `json:"severity"`
}

func allergiesToSlice(allergies []client.Allergy) []allergyJSON {
	result := make([]allergyJSON, 0, len(allergies))
	for _, a := range allergies {
		result = append(result, allergyJSON{Substance: a.Substance(), Severity: string(a.AllergyLevel())})
	}
	return result
}

func parseAllergies(data []byte) []client.Allergy {
	if len(data) == 0 {
		return nil
	}
	var items []allergyJSON
	if err := json.Unmarshal(data, &items); err != nil {
		return nil
	}
	result := make([]client.Allergy, 0, len(items))
	for _, item := range items {
		a, err := client.NewAllergy(item.Substance, client.Severity(item.Severity))
		if err == nil {
			result = append(result, a)
		}
	}
	return result
}

type noteJSON struct {
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

func notesToSlice(notes []client.Note) []noteJSON {
	result := make([]noteJSON, 0, len(notes))
	for _, n := range notes {
		result = append(result, noteJSON{Text: n.Text(), AuthorID: n.AuthorID().String(), CreatedAt: n.CreatedAt()})
	}
	return result
}

func parseNotes(data []byte) []client.Note {
	if len(data) == 0 {
		return nil
	}
	var items []noteJSON
	if err := json.Unmarshal(data, &items); err != nil {
		return nil
	}
	result := make([]client.Note, 0, len(items))
	for _, item := range items {
		authorID, _ := uuid.Parse(item.AuthorID)
		result = append(result, client.RestoreNote(item.Text, authorID, item.CreatedAt))
	}
	return result
}

func birthdayToNullable(b model.Birthday) *time.Time {
	if b.IsZero() {
		return nil
	}
	t := b.Time()
	return &t
}

// --- tx/pool wrappers to satisfy interface uniformly ---

type txWrapper struct {
	tx pgx.Tx
}

func (w *txWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.tx.QueryRow(ctx, sql, args...)
}

func (w *txWrapper) Exec(ctx context.Context, sql string, args ...any) (interface{ RowsAffected() int64 }, error) {
	tag, err := w.tx.Exec(ctx, sql, args...)
	return tag, err
}

type poolWrapper struct {
	pool *pgxpool.Pool
}

func (w *poolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.pool.QueryRow(ctx, sql, args...)
}

func (w *poolWrapper) Exec(ctx context.Context, sql string, args ...any) (interface{ RowsAffected() int64 }, error) {
	tag, err := w.pool.Exec(ctx, sql, args...)
	return tag, err
}
