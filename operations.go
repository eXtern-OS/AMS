package AMS

import (
	"context"
	beatrix "github.com/eXtern-OS/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func UpdateDB(acc Account) bool {

	var collection = newCollection("accounts")

	_, err := collection.InsertOne(context.Background(), acc)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error inserting into collection", "UPDATEDB")
		return false
	}
	return true

}

func GetPasswordHashed(login, password string) (string, string) {
	var collection = newCollection("accounts")

	filter := bson.M{"login": login}
	var acc Account
	var pwd, uid string
	err := collection.FindOne(context.Background(), filter).Decode(&acc)
	if err != nil {
		return "", ""
	} else {
		pwd = acc.Password
		uid = acc.UID
	}

	return pwd, uid
}

func GetUserByID(uid string) Account {
	var collection = newCollection("accounts")

	filter := bson.M{"uid": uid}

	var acc Account

	err := collection.FindOne(context.Background(), filter).Decode(&acc)
	if err != nil {
		log.Println(err)
		// Seems nothing found - better search on wordpress
		go beatrix.SendError("Error creating new mongo client", "GETUSERBYID")
		return Account{}
	}
	return acc
}

func UpdateDatabase(name, username, avatarurl, password, uid string) {
	var collection = newCollection("accounts")
	// Updating objects
	filter := bson.M{"uid": uid}
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"avatar_url": avatarurl,
			"password":   makehash(password),
			"username":   username,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error updating database", "AMS.UPDATEDATABASE")
	}
	return
}

func CheckIfExists(email string) bool {
	var collection = newCollection("accounts")
	fiter := bson.M{"email": email}

	var a Account
	err := collection.FindOne(context.Background(), fiter).Decode(&a)
	return err == nil
}
