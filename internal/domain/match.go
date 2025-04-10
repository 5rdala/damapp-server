package domain

type MatchStatus int

const (
	MatchStatusWaiting MatchStatus = iota
	MatchStatusStarted
	MatchStatusFinished
	MatchStatusStopped
)

func (s MatchStatus) String() string {
	switch s {
	case MatchStatusWaiting:
		return "waiting"
	case MatchStatusStarted:
		return "started"
	case MatchStatusFinished:
		return "finished"
	case MatchStatusStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Match struct {
	ID         uint64      `json:"id"`
	Code       int         `json:"code"`
	Player1ID  uint64      `json:"player_1_id"`
	Player2ID  uint64      `json:"player_2_id"`
	Winner     *uint64     `json:"winner"`
	Status     MatchStatus `json:"status"`
	CreatedAt  int64       `json:"created_at"`
	StartedAt  *int64      `json:"started_at"`
	FinishedAt *int64      `json:"finished_at"`
}
