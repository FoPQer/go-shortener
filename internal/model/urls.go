package model

// Urls represents a shortened URL entity in the domain model.
// generate:reset
type Urls struct {
	Original string `json:"original_url"`
	ShortURL string `json:"short_url"`
	UserID   string `json:"user_id,omitempty"`
	Deleted  bool   `json:"-"`
}

// NewUrls creates a new Urls entity with default non-deleted state.
func NewUrls(original, shortURL string) *Urls {
	return &Urls{
		Original: original,
		ShortURL: shortURL,
		Deleted:  false,
	}
}

// GetOriginal returns the original URL.
func (u *Urls) GetOriginal() string {
	return u.Original
}

// SetOriginal sets the original URL.
func (u *Urls) SetOriginal(original string) {
	u.Original = original
}

// GetShortURL returns the short URL token.
func (u *Urls) GetShortURL() string {
	return u.ShortURL
}

// SetShortURL sets the short URL token.
func (u *Urls) SetShortURL(shortURL string) {
	u.ShortURL = shortURL
}

// GetUserID returns the owner user ID.
func (u *Urls) GetUserID() string {
	return u.UserID
}

// SetUserID sets the owner user ID.
func (u *Urls) SetUserID(userID string) {
	u.UserID = userID
}

// IsDeleted reports whether URL is marked as deleted.
func (u *Urls) IsDeleted() bool {
	return u.Deleted
}

// SetDeleted sets the deleted flag.
func (u *Urls) SetDeleted(deleted bool) {
	u.Deleted = deleted
}
