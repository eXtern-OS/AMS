package AMS

import (
	"crypto/sha1"
	"encoding/base64"
	token "github.com/eXtern-OS/TokenMaster"
	"net/http"
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

func Register() {}
