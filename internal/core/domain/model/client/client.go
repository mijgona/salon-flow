package client

import (
	"errors"
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Client is the aggregate root for client management.
type Client struct {
	*ddd.BaseAggregate[uuid.UUID]
	tenantID     model.TenantID
	contactInfo  ContactInfo
	birthday     model.Birthday
	preferences  Preferences
	allergies    []Allergy
	notes        []Note
	photos       []Photo
	visitRecords []VisitRecord
	source       ClientSource
	registeredAt time.Time
}

// NewClient creates a new Client aggregate with validation and raises ClientRegistered event.
func NewClient(
	tenantID model.TenantID,
	contactInfo ContactInfo,
	source ClientSource,
	referredByClientID uuid.UUID,
) (*Client, error) {
	if !source.IsValid() {
		return nil, errs.NewErrValueMustBe("client source", "a valid source")
	}

	id := uuid.New()
	c := &Client{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		tenantID:      tenantID,
		contactInfo:   contactInfo,
		preferences:   EmptyPreferences(),
		allergies:     make([]Allergy, 0),
		notes:         make([]Note, 0),
		photos:        make([]Photo, 0),
		visitRecords:  make([]VisitRecord, 0),
		source:        source,
		registeredAt:  time.Now(),
	}

	c.RaiseDomainEvent(NewClientRegistered(
		id,
		tenantID.UUID(),
		contactInfo.FirstName(),
		contactInfo.LastName(),
		contactInfo.Phone().String(),
		source,
		referredByClientID,
	))

	return c, nil
}

// RestoreClient rehydrates a Client from the database (no validation, no events).
func RestoreClient(
	id uuid.UUID,
	tenantID model.TenantID,
	contactInfo ContactInfo,
	birthday model.Birthday,
	preferences Preferences,
	allergies []Allergy,
	notes []Note,
	photos []Photo,
	visitRecords []VisitRecord,
	source ClientSource,
	registeredAt time.Time,
) *Client {
	return &Client{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		tenantID:      tenantID,
		contactInfo:   contactInfo,
		birthday:      birthday,
		preferences:   preferences,
		allergies:     allergies,
		notes:         notes,
		photos:        photos,
		visitRecords:  visitRecords,
		source:        source,
		registeredAt:  registeredAt,
	}
}

// UpdateProfile updates the client's contact info, birthday, and preferences.
func (c *Client) UpdateProfile(contactInfo ContactInfo, birthday model.Birthday, preferences Preferences) {
	c.contactInfo = contactInfo
	c.birthday = birthday
	c.preferences = preferences
}

// AddAllergy adds an allergy with deduplication by substance.
func (c *Client) AddAllergy(allergy Allergy) error {
	for _, existing := range c.allergies {
		if existing.Equal(allergy) {
			return errors.New("allergy for this substance already exists")
		}
	}
	c.allergies = append(c.allergies, allergy)
	return nil
}

// AddVisitRecord adds a visit record to the client.
func (c *Client) AddVisitRecord(record VisitRecord) {
	c.visitRecords = append(c.visitRecords, record)
}

// AddNote adds a note to the client.
func (c *Client) AddNote(note Note) {
	c.notes = append(c.notes, note)
}

// AddPhoto adds a photo to the client.
func (c *Client) AddPhoto(photo Photo) {
	c.photos = append(c.photos, photo)
}

// TotalVisits returns the total number of visits.
func (c *Client) TotalVisits() int {
	return len(c.visitRecords)
}

// TotalSpent returns the total amount spent across all visits.
func (c *Client) TotalSpent() model.Money {
	total := model.ZeroMoney()
	for _, vr := range c.visitRecords {
		sum, err := total.Add(vr.Price())
		if err == nil {
			total = sum
		}
	}
	return total
}

// TotalSpentDecimal returns the total spent as a decimal (convenience).
func (c *Client) TotalSpentDecimal() decimal.Decimal {
	return c.TotalSpent().Amount()
}

// Getters
func (c *Client) TenantID() model.TenantID    { return c.tenantID }
func (c *Client) ContactInfo() ContactInfo    { return c.contactInfo }
func (c *Client) Birthday() model.Birthday    { return c.birthday }
func (c *Client) Preferences() Preferences    { return c.preferences }
func (c *Client) Allergies() []Allergy        { return c.allergies }
func (c *Client) Notes() []Note               { return c.notes }
func (c *Client) Photos() []Photo             { return c.photos }
func (c *Client) VisitRecords() []VisitRecord { return c.visitRecords }
func (c *Client) Source() ClientSource        { return c.source }
func (c *Client) RegisteredAt() time.Time     { return c.registeredAt }
