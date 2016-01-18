package postgres

import "database/sql"
import "io/ioutil"
import "fmt"
import "encoding/json"
import "errors"
import _ "github.com/lib/pq"

// once again I could do something really evil and have an open connection as a module global
// TODO: investigate best practices
// update: I tried something, namely have the settings in a json file somewhere.
// This raises a number of concerns;
// in any normal context I would be using a framework which would have a clear process for dealing with initialization and management of pools of resources and so on.
// Here, I'm not too sure how to proceed. 
// Let's do something: global_db will be a global containing the db connection in the postgres package, and this package will assume that the connection is open and ready to use. Main will call the function that will initialise that connection. 

var global_db *sql.DB

type Settings struct {
  Database string
  User string
  Password string
}

func getDatabaseOpts() (*Settings, error) {
  // assumes that it's called from $GOPATH/src/musicdb
  // kinda makes the idea of looking for music in the current directory moot
  // maybe read those from a regular dotfile?
  file, err := ioutil.ReadFile("./database_opts.json")
  if err != nil {
    fmt.Println("sorry, but this program only works when called from the $GOPATH/src/musicdb directory")
    return nil, err
  }
  var settings map[string]string
  err = json.Unmarshal(file, &settings)
  if err != nil {
    return nil, err
  }
  database, database_exists := settings["database"]
  user, user_exists := settings["user"]
  password, password_exists := settings["password"]
  if !(database_exists && user_exists && password_exists) {
    return nil, errors.New("must specify database, user and password in ./database_opts")
  }
  return &Settings{database, user, password}, nil
}


func InitDatabase() {
  settings, err := getDatabaseOpts()
  if err != nil {
    fmt.Println("error with database settings: ", err)
    panic("not recovering from that")
  }
  // TODO: find best practices for storing config data
  db, err := sql.Open("postgres", "user=" + settings.User + " password=" + settings.Password + " dbname=" + settings.Database + " sslmode=disable")
  if err != nil {
    panic(err.Error()) 
  }
  global_db = db
}

func InsertAlbumReturnId(album_name string, band_id int) int {
  stmtIns, err := global_db.Prepare("INSERT INTO album (name, band_id) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  // TODO: get type inference working in vim
  // if possible
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(album_name, band_id)
  if err != nil {
    panic(err.Error()) 
  }
  var album_id int
  err = global_db.QueryRow("SELECT last_value FROM album_id_seq").Scan(&album_id)
  if err != nil {
    panic(err.Error()) 
  }
  return album_id
}

func InsertBandReturnId(band string) int {
  stmtIns, err := global_db.Prepare("INSERT INTO band (name) VALUES ($1)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(band)
  if err != nil {
    panic(err.Error()) 
  }
  var band_id int
  err = global_db.QueryRow("SELECT last_value FROM band_id_seq").Scan(&band_id)
  if err != nil {
    panic(err.Error()) 
  }
  return band_id
}

// all these inserters should do some error management but eh
// also TODO: find what kind of ORM Go hackers like to use
func InsertTrack(track_name string, album_id int) {
  stmtIns, err := global_db.Prepare("INSERT INTO track (name, album_id) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(track_name, album_id)
  if err != nil {
    panic(err.Error()) 
  }
}
