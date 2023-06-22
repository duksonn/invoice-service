package domain

type Status int

const (
	ACTIVE Status = iota
	LOCKED
	WAITING_APPROVAL
	ACCEPTED
	REJECTED
)

func (s Status) String() string {
	switch s {
	case ACTIVE:
		return "ACTIVE"
	case LOCKED:
		return "LOCKED"
	case WAITING_APPROVAL:
		return "WAITING_APPROVAL"
	case ACCEPTED:
		return "ACCEPTED"
	case REJECTED:
		return "REJECTED"
	}
	return "unknown"
}
