package repository

import "space/internal/app/ds"

func (r *Repository) GetByEmail(email string) (ds.User, error) {
	var u ds.User
	err := r.db.First(&u, "email = ?", email).Error

	if err != nil {
		return ds.User{}, err
	}

	return u, nil
}

func (r *Repository) AddUser(user ds.User) (int, error) {
	result := r.db.Create(&user)

	if err := result.Error; err != nil {
		return 0, err
	}
	return int(user.ID), nil
}
