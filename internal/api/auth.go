package api

type UserSingleton struct {
	userID  int
	isAdmin bool
}

func loadUserData() UserSingleton {
	userData := UserSingleton{
		userID:  1,
		isAdmin: false,
	}
	return userData
}

func singleton() (int, bool, error) {
	userData := loadUserData()
	return userData.userID, userData.isAdmin, nil
}
