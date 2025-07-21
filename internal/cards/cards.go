package cards

type CardInput struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// PictureURL — адрес изображения
	PictureURL string `json:"picture_url"`
	// Price — цена
	Price float64 `json:"price"`
}

type CardOutput struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// PictureURL — адрес изображения
	PictureURL string `json:"picture_url"`
	// Price — цена
	Price float64 `json:"price"`
	// Username — автор
	Username string `json:"username"`
	// isOwned — признак принадлежности объявления текущему пользователю
	IsOwned bool
}

type CardsRepo interface {
	// PostACard создает новое объявление
	PostACard(crd *CardInput, userID string) (*CardInput, error)
	// GetCards получает ленту объявлений
	GetCards(params *QueryParams) ([]CardOutput, error)
}
