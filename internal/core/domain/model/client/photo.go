package client

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"
)

// PhotoType represents the type of photo.
type PhotoType string

const (
	PhotoTypeBefore  PhotoType = "before"
	PhotoTypeAfter   PhotoType = "after"
	PhotoTypeProfile PhotoType = "profile"
)

// Photo represents a photo attached to a client profile.
type Photo struct {
	url        string
	photoType  PhotoType
	uploadedAt time.Time
}

// NewPhoto creates a Photo value object with validation.
func NewPhoto(url string, photoType PhotoType) (Photo, error) {
	if url == "" {
		return Photo{}, errs.NewErrValueRequired("photo URL")
	}
	if photoType == "" {
		photoType = PhotoTypeProfile
	}
	return Photo{
		url:        url,
		photoType:  photoType,
		uploadedAt: time.Now(),
	}, nil
}

// MustNewPhoto creates a Photo or panics.
func MustNewPhoto(url string, photoType PhotoType) Photo {
	p, err := NewPhoto(url, photoType)
	if err != nil {
		panic(err)
	}
	return p
}

// RestorePhoto creates a Photo from persisted data.
func RestorePhoto(url string, photoType PhotoType, uploadedAt time.Time) Photo {
	return Photo{url: url, photoType: photoType, uploadedAt: uploadedAt}
}

// URL returns the photo URL.
func (p Photo) URL() string { return p.url }

// Type returns the photo type.
func (p Photo) Type() PhotoType { return p.photoType }

// UploadedAt returns the upload timestamp.
func (p Photo) UploadedAt() time.Time { return p.uploadedAt }
