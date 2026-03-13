package model

type Urls struct {
	Original string `json:"original_url"`
	ShortURL string `json:"short_url"`
}

func NewUrls(original, shortURL string) *Urls {
	return &Urls{
		Original: original,
		ShortURL: shortURL,
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
