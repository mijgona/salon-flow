package client

// ClientSource represents how a client was acquired.
type ClientSource string

const (
	ClientSourceOnlineBooking ClientSource = "online_booking"
	ClientSourceAdminEntry    ClientSource = "admin_entry"
	ClientSourceReferral      ClientSource = "referral"
	ClientSourceWalkIn        ClientSource = "walk_in"
)

// IsValid checks if the source is a known value.
func (cs ClientSource) IsValid() bool {
	switch cs {
	case ClientSourceOnlineBooking, ClientSourceAdminEntry, ClientSourceReferral, ClientSourceWalkIn:
		return true
	}
	return false
}

// String returns the string representation.
func (cs ClientSource) String() string { return string(cs) }
