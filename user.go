package main

type User struct {
	ID         UserID     `json:"id"`
	Properties []Property `json:"properties"`
}
