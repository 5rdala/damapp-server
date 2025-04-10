package domain

type FriendshipStatus int

const (
	FriendshipStatusPending FriendshipStatus = iota
	FriendshipStatusRejected
	FriendshipStatusAccepted
)

func (s FriendshipStatus) String() string {
	switch s {
	case FriendshipStatusPending:
		return "pending"
	case FriendshipStatusRejected:
		return "rejected"
	case FriendshipStatusAccepted:
		return "accepted"
	default:
		return "unknown"
	}
}

// UserID1 is the sender of the Friendship request
// UserID2 is the recipient of the Friendship request
type Friendship struct {
	ID        uint64           `json:"id" gorm:"primaryKey"`
	UserID1   uint64           `json:"user_id_1"`
	UserID2   uint64           `json:"user_id_2"`
	Status    FriendshipStatus `json:"status"`
	CreatedAt uint64           `json:"created_at"`
}
