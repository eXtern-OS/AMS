package AMS

import (
	"crypto/sha1"
	"encoding/base64"
	token "github.com/eXtern-OS/TokenMaster"
	"net/http"
	"strconv"
	"time"
)

func makehash(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Don't forget to init beatrix in main code.

func Init(mongoUri, mysqlURI string) {
	URI = mongoUri
	SQL_URI = mysqlURI
	token.Init(mongoUri)
	initDb()
}

func GetToken(login, password, ip string) (int, string) {
	if login == "" || password == "" {
		return http.StatusNoContent, ""
	}

	hashed, uid := GetPasswordHashed(login, password)

	if hashed == "" || uid == "" {
		return http.StatusNonAuthoritativeInfo, ""
	}

	if makehash(password) == hashed {
		code, t := token.NewToken(ip, uid)
		return code, t.TokenId
	} else {
		return http.StatusUnauthorized, ""
	}
}

func Register(name, username, login, avatarurl, pwd, website, email string) bool {
	var acc = Account{
		Login:      login,
		Password:   makehash(pwd),
		Username:   username,
		Name:       name,
		AvatarURL:  avatarurl,
		Developer:  false,
		Patreon:    false,
		Registered: time.Now().Format(time.RFC1123Z),
		Website:    website,
		Email:      email,
	}

	if CheckIfExists(email) {
		return false
	}

	acc.UID = makehash(acc.Password + login + acc.Registered + strconv.Itoa(random(1000, 2000)))
	return UpdateDB(acc)
}
