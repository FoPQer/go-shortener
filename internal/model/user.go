package model

// User represents an application user and their associated URLs.
type User struct {
	ID   string  `json:"id"`
	Urls []*Urls `json:"urls"`
}

// NewUser creates a new User with an initialized URL collection.
func NewUser(id string) *User {
	return &User{
		ID:   id,
		Urls: make([]*Urls, 0),
	}
}

// GetID returns the user identifier.
func (u *User) GetID() string {
	return u.ID
}

// SetID sets the user identifier.
func (u *User) SetID(id string) {
	u.ID = id
}

// GetURLs returns URLs associated with the user.
func (u *User) GetURLs() []*Urls {
	return u.Urls
}

// SetURLs replaces URLs associated with the user.
func (u *User) SetURLs(urls []*Urls) {
	u.Urls = urls
}

// AddURL appends a URL to the user's URL collection.
func (u *User) AddURL(url *Urls) {
	u.Urls = append(u.Urls, url)
}
