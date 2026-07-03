package model

// Stat represents aggregated counters returned by internal stats endpoint.
type Stat struct {
	TotalURLs  int `json:"urls"`
	TotalUsers int `json:"users"`
}

// NewStat creates a Stat with predefined URL and user counters.
func NewStat(totalURLs int, totalUsers int) *Stat {
	return &Stat{
		TotalURLs:  totalURLs,
		TotalUsers: totalUsers,
	}
}

// GetTotalURLs returns total amount of shortened URLs.
func (s *Stat) GetTotalURLs() int {
	return s.TotalURLs
}

// GetTotalUsers returns total amount of users.
func (s *Stat) GetTotalUsers() int {
	return s.TotalUsers
}

// IncrementURLs increases URL counter by count.
func (s *Stat) IncrementURLs(count int) {
	s.TotalURLs += count
}

// IncrementUsers increases user counter by count.
func (s *Stat) IncrementUsers(count int) {
	s.TotalUsers += count
}
