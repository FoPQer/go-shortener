package model

type Urls struct {
	Original string `json:"original_url"`
	ShortURL string `json:"short_url"`
	UserID   string `json:"user_id,omitempty"`
	Deleted  bool   `json:"-"`
}

func NewUrls(original, shortURL string) *Urls {
	return &Urls{
		Original: original,
		ShortURL: shortURL,
		Deleted:  false,
	}
}

func (u *Urls) GetOriginal() string {
	return u.Original
}

func (u *Urls) SetOriginal(original string) {
	u.Original = original
}

func (u *Urls) GetShortURL() string {
	return u.ShortURL
}

func (u *Urls) SetShortURL(shortURL string) {
	u.ShortURL = shortURL
}

func (u *Urls) GetUserID() string {
	return u.UserID
}

func (u *Urls) SetUserID(userID string) {
	u.UserID = userID
}

func (u *Urls) IsDeleted() bool {
	return u.Deleted
}

func (u *Urls) SetDeleted(deleted bool) {
	u.Deleted = deleted
}
