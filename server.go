package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

var name string

type user struct {
	Name    string
	Phone   string
	Iclass  string
	Message string
}

func Getlogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal(err)
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Fatal(err.Error())
		}
	}
}
func POSTlogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "POST" {
		r.ParseForm()
		Json, err := ioutil.ReadFile("config_example.json")
		if err != nil {
			log.Fatal(err)
		}
		JsonStr := string(Json)
		password := gjson.Get(JsonStr, "password")
		port := gjson.Get(JsonStr, "port")
		url := gjson.Get(JsonStr, "url")
		user := gjson.Get(JsonStr, "user")
		pd := password.String()
		pt := port.String()
		ul := url.String()
		ur := user.String()
		db, err := sql.Open("mysql", ur+":"+pd+"@tcp("+ul+":"+pt+")/test?charset=utf8")
		if err != nil {
			log.Fatal(err)
		}
		//查找数据
		find, err := db.Query("SELECT id FROM userinfo WHERE iclass=?", template.HTMLEscapeString(r.FormValue("iclass"))) //学号
		if err != nil {
			log.Fatal(err)
		}
		user_id := 0
		for find.Next() {
			find.Scan(&user_id)
		}
		if user_id != 0 {
			//更新数据
			up, err := db.Exec("UPDATE userinfo SET phone=?, message=? WHERE id=?", template.HTMLEscapeString(r.FormValue("phone")), template.HTMLEscapeString(r.FormValue("message")), user_id)
			fmt.Println(up)
			if err != nil {
				log.Fatal(err)
				w.Write([]byte("fail"))
			} else {
				w.Write([]byte("success"))
			}
			r.ParseForm()
		} else {
			//插入数据
			stmt, err := db.Prepare("INSERT userinfo SET name=?,phone=?,iclass=?,message=?")

			if err != nil {
				log.Fatal(err)
			}
			r.ParseForm()
			res, err := stmt.Exec(template.HTMLEscapeString(r.FormValue("name")), template.HTMLEscapeString(r.FormValue("phone")), template.HTMLEscapeString(r.FormValue("iclass")), template.HTMLEscapeString(r.FormValue("message")))
			if err != nil {
				log.Fatal(err)
				w.Write([]byte("fail"))
			} else {
				w.Write([]byte("success"))
			}
			fmt.Println(res)
		}
	}
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", Getlogin)
	http.HandleFunc("/login", POSTlogin)
	files := http.FileServer(http.Dir("Public"))
	rtr.PathPrefix("/").Handler(http.StripPrefix("/", files))
	http.Handle("/", rtr)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
