package files

import "fmt"
import "io/ioutil"

func NoOp() {
}

func Files() { fmt.Printf("hello, world\n")
  contents, _ := ioutil.ReadDir("/home/usr/music")
  // fmt.Println("wtf", contents)
  for _, element := range contents {
    fmt.Println("er", element.Name())
  }
}

