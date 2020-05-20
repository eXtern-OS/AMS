package token

import (
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"time"
)

type Token struct {
	TokenId string `bson:"id"`
	Expires int64  `bson:"expires"`
	Ip      string `bson:"ip"`
	Uid     string `bson:"uid"`
}

func makehash(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func NewToken(ip, uid string) (int, Token) {
	token := Token{
		TokenId: makehash(ip + uid + strconv.Itoa(int(time.Now().Unix()))),
		Expires: time.Now().Unix() + 24*60*60, // It's valid for 24 hours
		Ip:      ip,
		Uid:     uid,
	}
	return PutToken(token), token
}
