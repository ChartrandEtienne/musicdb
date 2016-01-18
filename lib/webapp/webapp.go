package webapp

// what should I do here. 
// serve / which is a page that will
// - open a websocket
// - request the list of bands
// - display the list in html

import "net/http"
import "fmt"

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Serve() {
  http.HandleFunc("/testing", handler)
  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
  })
  http.HandleFunc("/example.js", func (w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/example.js")
  })
  http.ListenAndServe(":8080", nil)
}
