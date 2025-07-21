package cards

import "fmt"

// PostACard создает новое объявление
func (repo *CardsDBRepository) PostACard(crd *CardInput, userID string) (*CardInput, error) {
	query := `INSERT INTO cards (title, card_text, picture_url, price, user_id)
	         VALUES ($1, $2, $3, $4, $5);`

	_, err := repo.dtb.Exec(query, crd.Title, crd.Text, crd.PictureURL, crd.Price, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: создание объявления: %v", err)
	}
	return crd, nil
}
