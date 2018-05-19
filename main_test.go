package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ekrem95/go-gin/router"

	"github.com/ekrem95/go-gin/db"
	"github.com/gin-gonic/gin"
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

func testRouter() *gin.Engine {
	os.Setenv("ENV", "TEST")
	if err := db.TestSQLConnection(); err != nil {
		log.Fatal(err)
	}

	return router.Default()
}

func TestSignup(t *testing.T) {
	testRouter := testRouter()

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
	_, err := db.Exec("delete from users where username=? limit 1", testUsername)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInvalidSignup(t *testing.T) {
	testRouter := testRouter()

	var form db.User

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
	testRouter := testRouter()

	var form db.User

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
	testRouter := testRouter()

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
	testRouter := testRouter()

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
	testRouter := testRouter()

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
	testRouter := testRouter()
	title := "Lion's Head Caves Adventure"

	// ----------------------
	var post = db.Post{Title: title, Description: "Description", PostedBy: testUsername, Src: "src"}

	data, _ := json.Marshal(post)

	req, err := http.NewRequest("POST", "/add", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", signinCookie)

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 200)

	var res struct {
		ID int `json:"id"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println(err)
	}

	// ----------------------

	articleID := res.ID //Lion's Head Caves Adventure

	url := fmt.Sprintf("/api/postbyid/%d", articleID)

	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", signinCookie)

	if err != nil {
		fmt.Println(err)
	}

	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data db.Post `json:"post"`
	}
	var p Data
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Equal(t, p.Data.Title, title)
}

func TestArticleNotFound(t *testing.T) {
	testRouter := testRouter()

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
	testRouter := testRouter()

	url := fmt.Sprintf("/api/posts")

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", signinCookie)

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data []db.Post `json:"posts"`
	}

	var p Data
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range p.Data {
		if reflect.ValueOf(v).Kind() != reflect.Struct {
			t.Error("Wrong data type.")
		}
	}

}

func TestAddArticle(t *testing.T) {
	testRouter := testRouter()

	var post db.Post

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
	testRouter := testRouter()

	var id string
	error := db.QueryRowScan("select id from posts where title ='"+testArticleTitle+"' limit 1", &id)
	if error != nil {
		t.Error(error)
	}

	comment := &db.Comment{Text: "Comment", PostID: id, Sender: testUsername}

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
	_, err := db.Exec("delete from posts where title=?", testArticleTitle)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeletePostByID(t *testing.T) {
	defer DeletePost()
	var id string
	// var post Post
	error := db.QueryRowScan("select id from posts where title ='"+testArticleTitle+" 'limit 1", &id)
	if error != nil {
		t.Error(error)
	}

	testRouter := testRouter()

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
	testRouter := testRouter()

	var post db.Post

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
