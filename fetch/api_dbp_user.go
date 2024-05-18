package fetch

import (
	"dataset"
	"dataset/request"
)

type DBPUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetDBPUser(req request.Request) (DBPUser, dataset.Status) {
	var status dataset.Status
	var u DBPUser
	u.Username = req.Username
	u.Email = req.Email
	return u, status
}

func GetTestUser() (DBPUser, dataset.Status) {
	var status dataset.Status
	var u DBPUser
	u.Username = `GaryNTest`
	u.Email = `gary@shortsands.com`
	return u, status
}
