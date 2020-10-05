package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	_ "github.com/go-sql-driver/mysql"
)

type Information struct {
  StructTitle       string
  StructHost        string
  StructExplanation string
}

func CheckError(err error) {
  if err != nil {
	  log.Fatal(err)
  }
}

func main() {
  http.HandleFunc("/home", topfunc)
  http.HandleFunc("/home/search", resultfunc)
  http.ListenAndServe(":8080", nil)
}

func topfunc(w http.ResponseWriter, r *http.Request) {
  indexFile, _ := template.ParseFiles("index.html")
  indexFile.Execute(w, "index.html")
}

func resultfunc(w http.ResponseWriter, r *http.Request) {
  searchFile, _ := template.ParseFiles("search.html")
  getinp := r.FormValue("keyword")
  myresult := mysqlopenfunc(getinp)
  htmlinsert := struct {
    Mese            string
    InformationSets string
  }{
    Mese: getinp,
    InformationSets: myresult,
  }
  searchFile.ExecuteTemplate(w, "search.html", htmlinsert)
}

func mysqlopenfunc(getinp string)string{
  db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/search")
  CheckError(err)
  defer db.Close()
  //I could only come up with a slightly cumbersome way to write it (LOL
  mysqlstatementLIKE := "SELECT * FROM search WHERE title LIKE'%" + getinp + "%';"
  mysqlsearch, err := db.Query(mysqlstatementLIKE) //Cognitively, a mysql statement!
  CheckError(err)
  defer mysqlsearch.Close()
  var (
    dbtitle     string
    dburl       string
    dbsetu      string
    slicestruct []Information//
  )
  for mysqlsearch.Next(){
    err := mysqlsearch.Scan(&dbtitle, &dburl, &dbsetu)
    CheckError(err)
    slicestruct = append(slicestruct, Information{StructTitle: dbtitle, StructHost: dburl, StructExplanation: dbsetu})
  }
  //The process of producing and returning the html code
  htmlcodes := ``
  for _,v := range slicestruct {
    htmlcodes += `<div class="zyou"><a href="`+ v.StructHost +`"><h3>`+ v.StructTitle +`</h3><a href="`+ v.StructHost +`">`+ v.StructHost +`</a><p>`+ v.StructExplanation +`</p></a></div>`
    //postscript  code: div.zyou/a/h3 && a && p
  }
  return htmlcodes
}