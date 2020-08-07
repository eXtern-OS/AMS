package AMS

import (
	"context"
	"database/sql"
	wphash "github.com/GerardSoleCa/wordpress-hash-go"
	beatrix "github.com/eXtern-OS/Beatrix"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
)

func SearchAndMerge(login, password string) (string, string) {
	db, err := sql.Open("mysql", SQL_URI)

	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error  connecting to mysqlDB", "SEARCHANDMERGE")
		return "", ""
	}
	passwordCommand := "SELECT user_pass from `wp_users` WHERE user_login = \"" + login + "\""
	var acc Account
	err = db.QueryRow(passwordCommand).Scan(&acc.Password)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error searching for password", "SEARCHANDMERGE")
		return "", ""
	}

	if !wphash.CheckPassword(password, acc.Password) {
		return "", ""
	}

	command := "SELECT user_login, user_nicename, user_email, user_registered, user_url, display_name from `wp_users` WHERE user_login = \"" + login + "\""

	err = db.QueryRow(command).Scan(&acc.Login, &acc.Username, &acc.Email, &acc.Registered, &acc.Website, &acc.Name)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error searching for userdata", "SEARCHANDMERGE")
		return "", ""
	}

	acc.Password = makehash(password)
	acc.UID = makehash(password + login + acc.Registered + strconv.Itoa(random(1000, 2000)))

	if !UpdateDB(acc) {
		return "", ""
	} else {
		return acc.Password, acc.UID
	}

}

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
		log.Println(err)
		// Seems nothing found - better search on wordpress
		pwd, uid = SearchAndMerge(login, password)
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
