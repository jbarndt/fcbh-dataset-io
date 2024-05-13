package fetch

import "dataset"

type DBPUser struct {
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func GetDBPUser() (DBPUser, dataset.Status) {
	var status dataset.Status
	var u DBPUser
	u.UserId = 99
	u.Username = `GaryNGriswold`
	u.FirstName = `Gary`
	u.LastName = `Griswold`
	u.Email = `gary@shortsands.com`
	return u, status
}

func GetTestUser() (DBPUser, dataset.Status) {
	var status dataset.Status
	var u DBPUser
	u.UserId = 99
	u.Username = `GaryNGriswold`
	u.FirstName = `Gary`
	u.LastName = `Griswold`
	u.Email = `gary@shortsands.com`
	return u, status
}
