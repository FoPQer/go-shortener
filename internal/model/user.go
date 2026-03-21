package model

type User struct {
	ID   string  `json:"id"`
	Urls []*Urls `json:"urls"`
}

func NewUser(id string) *User {
	return &User{
		ID:   id,
		Urls: make([]*Urls, 0),
	}
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetUrls() []*Urls {
	return u.Urls
}

func (u *User) AddUrl(url *Urls) {
	u.Urls = append(u.Urls, url)
}