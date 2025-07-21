package cards

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	PerPage  int
	Offset   int
	SortBy   string
	Order    string
	PriceMin *float64
	PriceMax *float64
	Username *string
}

// GetCards получает ленту объявлений
func (repo *CardsDBRepository) GetCards(params *QueryParams) ([]CardOutput, error) {
	baseQuery := `
        SELECT
            c.title,
            c.card_text,
            c.picture_url,
            c.price,
            u.username
        FROM cards c
        JOIN users u ON u.id = c.user_id
    `

	var whereClauses []string
	var args []interface{}
	argPos := 1

	if params.PriceMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("c.price >= $%d", argPos))
		args = append(args, *params.PriceMin)
		argPos++
	}
	if params.PriceMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("c.price <= $%d", argPos))
		args = append(args, *params.PriceMax)
		argPos++
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY c.%s %s", params.SortBy, params.Order)
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, params.PerPage, params.Offset)

	rows, err := repo.dtb.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []CardOutput
	for rows.Next() {
		var c CardOutput
		if err := rows.Scan(
			&c.Title,
			&c.Text,
			&c.PictureURL,
			&c.Price,
			&c.Username,
		); err != nil {
			return nil, err
		}
		if params.Username != nil && c.Username == *params.Username {
			c.IsOwned = true
		} else {
			c.IsOwned = false
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
