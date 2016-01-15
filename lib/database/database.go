package database

import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func NoOp() {
}

func Databases() {
  db, err := sql.Open("mysql", "usr:access@/driver_test")
  if err != nil {
    panic(err.Error())
  }
  rows, err := db.Query("select * from foobar")
  if err != nil {
    panic(err.Error()) 
  }
  cols, err := rows.Columns()
  if err != nil {
    panic(err.Error()) 
  }

  values := make([]sql.RawBytes, len(cols))

  scanArgs := make([]interface{}, len(values))

  for i := range values {
    scanArgs[i] = &values[i]
  }
  for rows.Next() {
    // get RawBytes from data
    err = rows.Scan(scanArgs...)
    if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
    }
    // Now do something with the data.
    // Here we just print each column as a string.
    var value string
    for i, col := range values {
      // Here we can check if the value is nil (NULL value)
      if col == nil {
        value = "NULL"
      } else {
        value = string(col)
      }
      fmt.Println(cols[i], ": ", value)
    }
    fmt.Println("-----------------------------------")
  }
  fmt.Println("res: ", cols)
}
