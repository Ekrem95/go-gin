package main

import (
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

func SetupRouter() *gin.Engine {
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
	r.GET("/api/getpostbyusername/:name", getPostByUsername)
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

// resp, err := http.PostForm("http://example.com/form",
// 	url.Values{"key": {"Value"}, "id": {"123"}})

func Setup() {
	r := SetupRouter()
	r.Run()
}

func TestDB(t *testing.T) {
	Setup()
	MySQL()
}

func TestSignin(t *testing.T) {
	testRouter := SetupRouter()

	// type SigninForm struct {
	// 	Username string `form:"username" json:"username"`
	// 	Password string `form:"password" json:"password"`
	// }
	//
	// var signinForm SigninForm
	//
	// signinForm.Username = testEmail
	// signinForm.Password = testPassword
	//
	// data, _ := json.Marshal(signinForm)

	form := url.Values{}
	form.Add("username", testUsername)
	form.Add("password", testPassword)

	req, error := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// req, error := http.NewRequest("POST", "/login", bytes.NewBufferString(string(data)))
	// req.Header.Set("Content-Type", "application/json")

	if error != nil {
		fmt.Println(error)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	signinCookie = resp.Header().Get("Set-Cookie")

	assert.Equal(t, resp.Code, 200)
	fmt.Println(resp.Body)
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
