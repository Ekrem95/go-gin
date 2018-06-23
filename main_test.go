package main

import (
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

	"github.com/ekrem95/go-gin/db"
	"github.com/ekrem95/go-gin/router"
	"github.com/stretchr/testify/assert"
)

var active bool
var testRouter = router.Default()
var signinCookie string
var testUsername = "test"
var testPassword = "123456"
var testArticleTitle = "Test article title"

type UserForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func clearDatabase() error {
	var table string
	var tables []string

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	for _, v := range tables {
		if _, err = db.Exec("TRUNCATE TABLE " + v); err != nil {
			return err
		}
	}
	return nil
}

func setup() {
	if !active {
		os.Setenv("ENV", "TEST")
		if err := db.CheckSQLConnection(); err != nil {
			log.Fatal(err)
		}
		if err := clearDatabase(); err != nil {
			log.Fatal(err)
		}
		active = true
	}
}

type Request struct {
	method       string
	addr         string
	form         url.Values
	signinCookie string
}

func request(r Request) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(r.method, r.addr, strings.NewReader(r.form.Encode()))
	if err != nil {
		return nil, err
	}
	req.PostForm = r.form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if r.signinCookie != "" {
		req.Header.Set("Cookie", r.signinCookie)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	return resp, nil
}

func TestSignup(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("username", testUsername)
	form.Add("password", testPassword)

	resp, err := request(Request{"POST", "/signup", form, ""})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 200)
}

func TestInvalidSignup(t *testing.T) {
	setup()

	form := url.Values{}
	resp, err := request(Request{"POST", "/signup", form, ""})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 400)
}

func TestInvalidLogin(t *testing.T) {
	setup()

	form := url.Values{}
	resp, err := request(Request{"POST", "/login", form, ""})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 400)
}

func TestLogin(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("username", testUsername)
	form.Add("password", testPassword)

	resp, err := request(Request{"POST", "/login", form, ""})
	if err != nil {
		t.Error(err)
	}

	signinCookie = resp.Header().Get("Set-Cookie")
	assert.Equal(t, resp.Code, 200)
}

func TestChangePassword(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("current", testPassword)
	form.Add("newPassword", "newPassword")

	req, err := http.NewRequest("POST", "/changepassword", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", signinCookie)

	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestChangePasswordFail(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("current", testPassword)
	form.Add("newPassword", "newPassword")

	resp, err := request(Request{"POST", "/changepassword", form, ""})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Code, 400)
}

func TestGetArticle(t *testing.T) {
	setup()
	title := "Lion's Head Caves Adventure"

	form := url.Values{}
	form.Add("title", title)
	form.Add("src", "src.co")
	form.Add("description", "Description")
	form.Add("posted_by", testUsername)

	resp, err := request(Request{"POST", "/add", form, signinCookie})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 200)

	var res struct {
		ID int `json:"id"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println(err)
	}

	articleID := res.ID //Lion's Head Caves Adventure

	addr := fmt.Sprintf("/api/postbyid/%d", articleID)

	resp, err = request(Request{"GET", addr, nil, signinCookie})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data db.Post `json:"post"`
	}
	var p Data
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, p.Data.Title, title)
}

func TestArticleNotFound(t *testing.T) {
	setup()

	articleID := 123456
	addr := fmt.Sprintf("/api/postbyid/%d", articleID)
	resp, err := request(Request{"GET", addr, nil, signinCookie})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Code, 400)
}

func TestGetArticles(t *testing.T) {
	setup()

	addr := fmt.Sprintf("/api/posts")
	resp, err := request(Request{"GET", addr, nil, signinCookie})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Code, 200)

	type Data struct {
		Data []db.Post `json:"posts"`
	}

	var p Data
	if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Error(err)
	}

	for _, v := range p.Data {
		if reflect.ValueOf(v).Kind() != reflect.Struct {
			t.Error("Wrong data type.")
		}
	}

}

func TestAddArticle(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("title", testArticleTitle)
	form.Add("src", "src.co")
	form.Add("description", "Test Description")
	form.Add("posted_by", testUsername)

	resp, err := request(Request{"POST", "/add", form, signinCookie})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 200)
}

func TestPostComment(t *testing.T) {
	setup()

	var id string
	error := db.QueryRowScan("select id from posts where title ='"+testArticleTitle+"' limit 1", &id)
	if error != nil {
		t.Error(error)
	}

	form := url.Values{}
	form.Add("text", "Comment")
	form.Add("post_id", id)
	form.Add("sender", testUsername)

	resp, err := request(Request{"POST", "/comment", form, ""})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 200)
}

func TestDeletePostByID(t *testing.T) {
	var id string
	// var post Post
	error := db.QueryRowScan("select id from posts where title ='"+testArticleTitle+" 'limit 1", &id)
	if error != nil {
		t.Error(error)
	}

	setup()

	form := url.Values{}
	form.Add("id", id)
	form.Add("user", testUsername)

	resp, err := request(Request{"POST", "/delete/" + id, form, signinCookie})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Code, 200)
}

func TestCreateInvalidArticle(t *testing.T) {
	setup()

	form := url.Values{}
	form.Add("title", testArticleTitle)

	resp, err := request(Request{"POST", "/add", form, signinCookie})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.Code, 400)
}
