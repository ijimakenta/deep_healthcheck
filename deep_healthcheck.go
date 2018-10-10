package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "database/sql"
  _ "github.com/lib/pq"
  "gopkg.in/yaml.v2"
  "gopkg.in/alecthomas/kingpin.v2"
)


type Config struct {
  URL       []string `yaml:"url"`
  Postgres  Postgres
}
type Postgres struct {
  User      string `yaml:"user"`
  Password  string `yaml:"password"`
  Host      string `yaml:"host"`
  DBname    string `yaml:"dbname"`
  SSLmode   string `yaml:"sslmode"`
}

var c Config

var (
  listenAddress  = kingpin.Flag("listen-address", "Address on which to expose metrics and web interface.").Default(":1234").String()
  configPath     = kingpin.Flag("config-path", "Path under which to yml path").Default("").String()
)

func handler(w http.ResponseWriter, r *http.Request) {
  var flag bool

  for i := len(c.URL); i > 0 ; i-- {
    var idx = i - 1

    resp, err := http.Get(c.URL[idx])
    if err != nil {
      fmt.Printf("%s NG\n", err.Error())
  	}

    if resp.StatusCode != 200 {
      defer resp.Body.Close()
      desc := fmt.Sprintf("URL: %s StatusCode: %s", c.URL[idx], resp.StatusCode)
      fmt.Printf("%s NG\n",desc)
      flag = true

    } else {
      defer resp.Body.Close()
      fmt.Printf("URL: %s StatusCocde: %s OK\n", c.URL[idx], resp.StatusCode)
    }
  }

  dbconnect := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
    c.Postgres.User ,c.Postgres.Password, c.Postgres.Host, c.Postgres.DBname, c.Postgres.SSLmode)
  db, err := sql.Open("postgres", dbconnect)
  if err != nil {
    fmt.Printf("%s NG\n", err.Error())
    flag = true
  }
  defer db.Close()

  rows, err := db.Query("select 1;")
  if err != nil {
    fmt.Printf("%s NG\n", err.Error())
    flag = true

  } else {
    for rows.Next() {
      var data string
      rows.Scan(&data)
      fmt.Printf("Postgres: select 1; result:%s OK\n", data)
    }
    defer rows.Close()
  }

  if flag == true {
    fmt.Println("Health check NG")
    http.Error(w, "Health check NG", http.StatusInternalServerError)
    return
  }

  fmt.Println("Health check OK")
  w.Header().Set("Content-Type","text/plain")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Health check OK"))
}

func main() {
  kingpin.Parse()
  fmt.Printf("listen address: %s, configration path: %s\n", *listenAddress, *configPath)

  source, err := ioutil.ReadFile(*configPath)
  if err != nil {
    panic(err)
  }
  fmt.Printf("source: %+v\n", string(source))

  err = yaml.Unmarshal(source, &c)
  if err != nil {
    panic(err)
  }
  fmt.Printf("Value: %+v\n", c)

  http.HandleFunc("/health/check", handler)
  http.ListenAndServe(*listenAddress, nil)
}
