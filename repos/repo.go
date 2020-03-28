package repos

import (
	. "github.com/CurtizJ/dummy-shop/items"
)

type Repo interface {
	Get(id uint64) (*Item, error)
	Add(item *Item) error
	Update(newItem *Item) error
	Delete(id uint64) error

	List(length, offset uint64) ([]Item, error)
	ListAll() ([]Item, error)
}
