package scheduling

import (
	"github.com/mijgona/salon-crm/internal/core/domain/model"
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"

	"github.com/google/uuid"
)

// ServiceInfo holds information about the service being performed.
type ServiceInfo struct {
	serviceID uuid.UUID
	name      string
	duration  time.Duration
	basePrice model.Money
}

// NewServiceInfo creates a ServiceInfo value object with validation.
func NewServiceInfo(serviceID uuid.UUID, name string, duration time.Duration, basePrice model.Money) (ServiceInfo, error) {
	if serviceID == uuid.Nil {
		return ServiceInfo{}, errs.NewErrValueRequired("service ID")
	}
	if name == "" {
		return ServiceInfo{}, errs.NewErrValueRequired("service name")
	}
	if duration <= 0 {
		return ServiceInfo{}, errs.NewErrValueMustBe("duration", "positive")
	}
	return ServiceInfo{
		serviceID: serviceID,
		name:      name,
		duration:  duration,
		basePrice: basePrice,
	}, nil
}

// MustNewServiceInfo creates a ServiceInfo or panics.
func MustNewServiceInfo(serviceID uuid.UUID, name string, duration time.Duration, basePrice model.Money) ServiceInfo {
	si, err := NewServiceInfo(serviceID, name, duration, basePrice)
	if err != nil {
		panic(err)
	}
	return si
}

// ServiceID returns the service UUID.
func (si ServiceInfo) ServiceID() uuid.UUID { return si.serviceID }

// Name returns the service name.
func (si ServiceInfo) Name() string { return si.name }

// Duration returns the service duration.
func (si ServiceInfo) Duration() time.Duration { return si.duration }

// BasePrice returns the base price.
func (si ServiceInfo) BasePrice() model.Money { return si.basePrice }
