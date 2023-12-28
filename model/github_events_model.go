package model

import (
	"time"
)

type Event struct {
	ID             string    `bson:"_id"`
	Type           string    `bson:"type"`
	CreatedAt      time.Time `bson:"created_at"`
	Public         bool      `bson:"public"`
	RepoFullName   string    `bson:"repo_full_name"`
	RepoUrl        string    `bson:"repo_url"`
	ActorLogin     string    `bson:"actor_login"`
	ActorId        int64     `bson:"actor_id"`
	ActorUrl       string    `bson:"actor_url"`
	ActorAvatarUrl string    `bson:"actor_avatar_url"`
}

type Repo struct {
	ID            string    `bson:"_id"`
	Owner         string    `bson:"owner"`
	Name          string    `bson:"name"`
	Url           string    `bson:"url"`
	Stars         int       `bson:"stars"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}

type User struct {
	ID            int64     `bson:"_id"`
	Login         string    `bson:"login"`
	Url           string    `bson:"url"`
	AvatarUrl     string    `bson:"avatar_url"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}
