package scheduling

// BookingSource represents how an appointment was booked.
type BookingSource string

const (
	BookingSourceOnline BookingSource = "online"
	BookingSourceAdmin  BookingSource = "admin"
)

// IsValid checks if the booking source is a known value.
func (bs BookingSource) IsValid() bool {
	switch bs {
	case BookingSourceOnline, BookingSourceAdmin:
		return true
	}
	return false
}

// String returns the string representation.
func (bs BookingSource) String() string { return string(bs) }
