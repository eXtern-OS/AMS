package AMS

type Account struct {
	Login    string `bson:"login" json:"login"`
	Password string `bson:"password" json:"password"`
	UID      string `bson:"uid" json:"uid"`

	Username   string `bson:"username" json:"username"`
	Name       string `bson:"name" json:"name"`
	AvatarURL  string `bson:"avatar_url" json:"avatar_url"`
	Developer  bool   `bson:"developer" json:"developer"`
	Patreon    bool   `bson:"patreon" json:"patreon"`
	Registered string `bson:"registered" json:"registered"`
	Website    string `bson:"website" json:"website"`
	Email      string `bson:"email" json:"email"`

	// To be implemented - payment info
}
