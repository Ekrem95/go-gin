package main

// User type
type User struct {
	Name     string
	Password string
}

// Message type
type Message struct {
	Text   string
	Time   string
	Sender string
}

// Post type
type Post struct {
	id          int
	title       string
	src         string
	description string
	// comments    []string
	likes int
}
