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
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Src         string `json:"src"`
	Description string `json:"description"`
	Likes       int    `json:"likes"`
	PostedBy    string `json:"postedBy"`
}

// Comment type
type Comment struct {
	Text   string `json:"text"`
	PostID string `json:"post_id"`
	Time   int64  `json:"time"`
	Sender string `json:"sender"`
}

// Like type
type Like struct {
	PostID string `json:"postID"`
	User   string `json:"user"`
}
