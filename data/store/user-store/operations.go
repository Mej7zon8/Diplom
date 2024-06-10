package user_store

import (
	"errors"
	"gorm.io/gorm"
	"messenger/data/entities"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

// Create creates a new record in the database for the specified user.
func (s *Store) Create(update func(user *entities.User)) (entities.User, error) {
	var item = entities.NewUser()
	update(item)

	return *item, s.table.Transaction(func(tx *gorm.DB) error {
		e := tx.Create(item).Error
		if e != nil && e.Error() == "UNIQUE constraint failed: users.id" {
			e = ErrUserAlreadyExists
		}
		return e
	})
}

func (s *Store) GetAll() ([]entities.User, error) {
	var items []entities.User
	return items, s.table.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&items).Error
	})
}

// Read retrieves the record from the database for the specified user.
func (s *Store) Read(userRef entities.UserRef) (*entities.User, error) {
	var item entities.User
	return &item, s.table.Transaction(func(tx *gorm.DB) error {
		e := tx.First(&item, userRef).Error
		if e != nil && e.Error() == "record not found" {
			return ErrUserNotFound
		}
		return e
	})
}

// Update updates the record in the database for the specified user.
func (s *Store) Update(userRef entities.UserRef, update func(user *entities.User)) error {
	return s.table.Transaction(func(tx *gorm.DB) error {
		var item, err = s.Read(userRef)
		if err != nil {
			return err
		}

		update(item)

		return tx.Save(item).Error
	})
}

func (s *Store) Delete(userRef entities.UserRef) error {
	return s.table.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&entities.User{}, userRef).Error
	})
}
