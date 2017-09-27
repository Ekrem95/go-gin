package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

var signinCookie string

var testUsername = "test"
var testPassword = "123456"
var testArticleTitle = "Test article title"

type UserForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func SetupRouter() *gin.Engine {
	MySQL()
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("session_test", store))
	r.LoadHTMLGlob("../templates/*")

	r.StaticFS("/src", http.Dir("../src"))
	r.StaticFile("/favicon.ico", "../templates/favicon.ico")

	// socketio
	server, socketErr := socketio.NewServer(nil)
	if socketErr != nil {
		log.Fatal(socketErr)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("msg", func(msg *Message) {
			so.BroadcastTo("chat", "dist", msg)
			RedisSaveMsg(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	r.GET("/", common)
	r.GET("/signup", common)
	r.GET("/login", common)
	r.GET("/add", common)
	r.GET("/upload", common)
	r.GET("/user", getUser)
	r.GET("/messages", RedisGetMsgs)
	r.GET("/api/posts", getPosts)
	r.GET("/api/postbyid/:id", getPostByID)
	r.GET("/api/commentsbyid/:id", getCommentsByID)
	r.GET("/p/*all", common)
	r.GET("/myposts", common)
	r.GET("/api/getpostbyusername/:name", getPostsByUsername)
	r.GET("/edit/:id", common)
	r.GET("/changepassword", common)
	r.GET("/get_likes/:id", getLikes)

	r.POST("/signup", signup)
	r.POST("/login", login)
	r.POST("/logout", logout)
	r.POST("/add", addPost)
	r.POST("/comment", postComment)
	r.POST("/upload", uploadFile)
	r.POST("/edit/:id", editPost)
	r.POST("/delete/:id", deletePostByID)
	r.POST("/changepassword", changePassword)
	r.POST("/post_likes", postLikes)

	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/", gin.WrapH(server))

	return r
}

func Setup() {
	r := SetupRouter()
	r.Run()
}

func TestSignup(t *testing.T) {
	testRouter := SetupRouter()

	form := url.Values{}
	form.Add("username", testUsername)
	form.Add("password", testPassword)

	req, error := http.NewRequest("POST", "/signup", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func DeleteUser() {
	_, err = db.Exec("delete from users where username=? limit 1", testUsername)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInvalidSignup(t *testing.T) {
	testRouter := SetupRouter()

	var form User

	data, _ := json.Marshal(form)
	req, error := http.NewRequest("POST", "/signup", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400) //406
}

func TestInvalidLogin(t *testing.T) {
	testRouter := SetupRouter()

	var form User

	data, _ := json.Marshal(form)
	req, error := http.NewRequest("POST", "/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400) //406
}

func TestLogin(t *testing.T) {
	testRouter := SetupRouter()

	form := url.Values{}
	form.Add("username", testUsername)
	form.Add("password", testPassword)

	req, error := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	signinCookie = resp.Header().Get("Set-Cookie")

	assert.Equal(t, resp.Code, 200)
}

func TestChangePassword(t *testing.T) {
	testRouter := SetupRouter()

	form := url.Values{}
	form.Add("current", testPassword)
	form.Add("newPassword", "newPassword")

	req, error := http.NewRequest("POST", "/changepassword", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestChangePasswordFail(t *testing.T) {
	testRouter := SetupRouter()

	form := url.Values{}
	form.Add("current", testPassword)
	form.Add("newPassword", "newPassword")

	req, error := http.NewRequest("POST", "/changepassword", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400)
}

func TestGetArticle(t *testing.T) {
	testRouter := SetupRouter()

	articleID := 12 //Lion's Head Caves Adventure

	url := fmt.Sprintf("/api/postbyid/%d", articleID)

	req, error := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data Post `json:"post"`
	}
	var p Data
	error = json.NewDecoder(resp.Body).Decode(&p)
	if error != nil {
		fmt.Println(error)
		return
	}

	assert.Equal(t, p.Data.Title, "Lion's Head Caves Adventure")
}

func TestArticleNotFound(t *testing.T) {
	testRouter := SetupRouter()

	articleID := 123456 //Lion's Head Caves Adventure

	url := fmt.Sprintf("/api/postbyid/%d", articleID)

	req, error := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)
}

func TestGetArticles(t *testing.T) {
	testRouter := SetupRouter()
	// var post Post

	url := fmt.Sprintf("/api/posts")

	req, error := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data []Post `json:"posts"`
	}

	var p Data
	error = json.NewDecoder(resp.Body).Decode(&p)
	if error != nil {
		fmt.Println(error)
		return
	}

	for _, v := range p.Data {
		if reflect.ValueOf(v).Kind() != reflect.Struct {
			t.Error("Wrong data type.")
		}
	}

}

func TestCreateArticle(t *testing.T) {
	// defer DeletePost()
	testRouter := SetupRouter()

	var post Post

	post.Title = testArticleTitle
	post.Description = "Test description"
	post.PostedBy = testUsername
	post.Src = "source"

	data, _ := json.Marshal(post)

	req, error := http.NewRequest("POST", "/add", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 200)
}

func TestPostComment(t *testing.T) {
	testRouter := SetupRouter()

	var id string
	error := db.QueryRow("select id from posts where title =? limit 1", testArticleTitle).Scan(&id)
	if error != nil {
		t.Error(error)
	}

	comment := &Comment{Text: "Comment", PostID: id, Sender: testUsername}

	data, _ := json.Marshal(comment)

	req, error := http.NewRequest("POST", "/comment", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 200)
}

func DeletePost() {
	_, err = db.Exec("delete from posts where title=?", testArticleTitle)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeletePostByID(t *testing.T) {
	defer DeletePost()
	var id string
	// var post Post
	error := db.QueryRow("select id from posts where title =? limit 1", testArticleTitle).Scan(&id)
	if error != nil {
		t.Error(error)
	}

	testRouter := SetupRouter()

	form := url.Values{}
	form.Add("id", id)
	form.Add("user", testUsername)

	req, error := http.NewRequest("POST", "/delete/"+id, strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestCreateInvalidArticle(t *testing.T) {
	defer DeleteUser()
	testRouter := SetupRouter()

	var post Post

	post.Title = testArticleTitle

	data, _ := json.Marshal(post)

	req, error := http.NewRequest("POST", "/add", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", signinCookie)

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400) //406
}
