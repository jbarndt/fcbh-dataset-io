package fetch

import (
	log "dataset/logger"
	"dataset/request"
)

type DBPUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetDBPUser(req request.Request) (DBPUser, *log.Status) {
	var u DBPUser
	u.Username = req.Username
	u.Email = req.Email
	return u, nil
}

func GetTestUser() (DBPUser, *log.Status) {
	var u DBPUser
	u.Username = `GaryNTest`
	u.Email = `gary@shortsands.com`
	return u, nil
}
