package entity

type User struct {
	Id   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
