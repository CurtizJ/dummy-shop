package repos

import (
	"errors"

	. "github.com/CurtizJ/dummy-shop/api/items"
	. "github.com/CurtizJ/dummy-shop/errors"
)

type MemoryRepo struct {
	items     map[uint64]*Item
	currentId uint64
}

func NewMemoryRepo() Repo {
	return &MemoryRepo{make(map[uint64]*Item, 0), 0}
}

func (repo *MemoryRepo) Get(Id uint64) (*Item, error) {
	if res, found := repo.items[Id]; found {
		return res, nil
	}

	return nil, &ItemNotFoundError{Id}
}

func (repo *MemoryRepo) Add(item *Item) error {
	item.Id = repo.currentId
	repo.currentId++
	repo.items[item.Id] = item
	return nil
}

func (repo *MemoryRepo) Update(newItem *Item) error {
	if oldItem, found := repo.items[newItem.Id]; found {
		*oldItem = *newItem
		return nil
	}

	return &ItemNotFoundError{newItem.Id}
}

func (repo *MemoryRepo) Delete(Id uint64) error {
	if _, found := repo.items[Id]; found {
		delete(repo.items, Id)
		return nil
	}

	return &ItemNotFoundError{Id}
}

func (repo *MemoryRepo) ListAll() ([]Item, error) {
	res := make([]Item, 0, len(repo.items))
	for _, item := range repo.items {
		res = append(res, *item)
	}

	return res, nil
}

func (repo *MemoryRepo) List(length, offset uint64) ([]Item, error) {
	return nil, errors.New("Not implemented")
}
