package main

import "fmt"
import "hello/lib/files"
import "hello/lib/database/mysql"
import "hello/lib/database/postgres"


func main() {
  fmt.Println("right.")
  files.Walk()
  mysql.Databases()
  postgres.TestInsert()
}
