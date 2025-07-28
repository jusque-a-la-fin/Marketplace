package cards

type CardInput struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// ImageURL — адрес изображения
	ImageURL string `json:"image_url"`
	// Price — цена
	Price float64 `json:"price"`
}

type CardOutput struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// ImageURL — адрес изображения
	ImageURL string `json:"image_url"`
	// Price — цена
	Price float64 `json:"price"`
	// Username — автор
	Username string `json:"username"`
	// isOwned — признак принадлежности объявления текущему пользователю
	IsOwned bool `json:"is_owned,omitempty"`
}

type CardsRepo interface {
	// PostACard создает новое объявление
	PostACard(crd *CardInput, userID string) (*CardInput, error)
	// GetCards получает ленту объявлений
	GetCards(params *QueryParams) ([]CardOutput, error)
}
