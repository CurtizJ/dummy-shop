package items

type Item struct {
	Name     string `json:"name"`
	Id       uint64 `json:"id"`
	Category string `json:"category"`
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func FromRow(rows scannable) (Item, error) {
	item := Item{}
	err := rows.Scan(&item.Id, &item.Name, &item.Category)
	return item, err
}
