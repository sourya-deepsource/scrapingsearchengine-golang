package main

import (
  "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Error(err error){
  if err != nil{
    log.Fatal(err)
  }
}

type Information struct{
  Title       string
  Host        string
  Explanation string
}

func main(){
  url := "http://www.tohoho-web.com/"
	geturlfunc(url)
	//mysql -> urlカラム -> 関数呼び出し(名(geturlfunc))
}
func geturlfunc(url string){
  resp, err := http.Get(url)
	Error(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Error(err)
	z := html.NewTokenizer(strings.NewReader(string(body)))
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, v := range t.Attr {
					sh := strings.HasPrefix(v.Val, "#")
					if sh == false {
						http := strings.HasPrefix(v.Val, "http")
						switch http {
						case false:
							host := strings.HasPrefix(v.Val, url)
							switch host {
							case false:
								sc := strings.HasSuffix(v.Val, "/")
								switch sc {
								case false:
									sc2 := strings.HasPrefix(v.Val, "/")
									switch sc2 {
									case true:
										acquisitionfunc(url + "/" + v.Val)
									case false:
										acquisitionfunc(url + v.Val)
									}
								case true:
									acquisitionfunc(url + v.Val)
								}
							default:
								acquisitionfunc(v.Val)
							}
						case true:
							acquisitionfunc(v.Val)
						}
					}
				}
			}
		}
	}
}

func acquisitionfunc(url string){
  resp, err := http.Get(url)
	Error(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Error(err)
	z := html.NewTokenizer(strings.NewReader(string(body)))
	titles := titlegetfunc(z)
	descriptions := descriptiongetfunc(z)
	ke := Information{titles, url, descriptions}
	ke.databaseinsertfunc()
}

func titlegetfunc(ztitle *html.Tokenizer) (titlereturn string){
  for {
		tt := ztitle.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := ztitle.Token()
			if t.Data == "title" {
				ztitle.Next()
				i := ztitle.Token()
				return i.Data
			}
		}
	}
}

func descriptiongetfunc(zdescription *html.Tokenizer) (descriptionreturn string){
  for {
		tt := zdescription.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := zdescription.Token()
			c := t.Data == "p"
			switch c {
			case true:
				zdescription.Next()
				i := zdescription.Token()
				return i.Data
			case false:
				if t.Data == "meta" {
					for _, v := range t.Attr {
						if v.Key == "name" {
							if v.Val == "description" {
								for _, v := range t.Attr {
									if v.Key == "content" {
										return v.Val
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (accept Information) databaseinsertfunc(){
  db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/search")
	Error(err)
	defer db.Close()
	informationinsert, err := db.Prepare("INSERT INTO search(title,url,setu) VALUES(?,?,?)")
	Error(err)
	informationinsert.Exec(accept.Title, accept.Host, accept.Explanation)
}