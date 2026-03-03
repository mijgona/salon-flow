package client

import (
	"github.com/mijgona/salon-crm/internal/pkg/errs"
	"time"

	"github.com/google/uuid"
)

// Note represents a note attached to a client profile.
type Note struct {
	text      string
	authorID  uuid.UUID
	createdAt time.Time
}

// NewNote creates a Note value object with validation.
func NewNote(text string, authorID uuid.UUID) (Note, error) {
	if text == "" {
		return Note{}, errs.NewErrValueRequired("note text")
	}
	if authorID == uuid.Nil {
		return Note{}, errs.NewErrValueRequired("author ID")
	}
	return Note{
		text:      text,
		authorID:  authorID,
		createdAt: time.Now(),
	}, nil
}

// MustNewNote creates a Note or panics.
func MustNewNote(text string, authorID uuid.UUID) Note {
	n, err := NewNote(text, authorID)
	if err != nil {
		panic(err)
	}
	return n
}

// RestoreNote creates a Note from persisted data (no validation).
func RestoreNote(text string, authorID uuid.UUID, createdAt time.Time) Note {
	return Note{text: text, authorID: authorID, createdAt: createdAt}
}

// Text returns the note text.
func (n Note) Text() string { return n.text }

// AuthorID returns the author's UUID.
func (n Note) AuthorID() uuid.UUID { return n.authorID }

// CreatedAt returns the creation timestamp.
func (n Note) CreatedAt() time.Time { return n.createdAt }
