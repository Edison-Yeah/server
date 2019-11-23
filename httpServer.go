package main

import (
    "flag"
    "github.com/aymerick/raymond"
    "github.com/fvbock/endless"
    "github.com/go-zoo/bone"
    "log"
    "net/http"
    "os"
    "syscall"
)

var (
    //homeTpl, _ = raymond.ParseFile("home.hbs")
    homeTpl = raymond.MustParse(`<html>
<head>
<title>test</title>
</head>
</body>
<div class="entry">
<h1></h1>
<div class="body">

</div>
</div>
</body>
</html>
`)
)

func homeHandler(rw http.ResponseWriter, req *http.Request) {
    ctx := map[string]string{"greet": "hello", "name": "world"}
    result := homeTpl.MustExec(ctx)
    rw.Write([]byte(result))
}
func varHandler(rw http.ResponseWriter, req *http.Request) {
    varr := bone.GetValue(req, "var")
    test := bone.GetValue(req, "test")

    rw.Write([]byte(varr + " " + test))
}
func Handler404(rw http.ResponseWriter, req *http.Request) {
    rw.Write([]byte("These are not resources you're looking for ..."))
}
func restartHandler(rw http.ResponseWriter, req *http.Request) {
    syscall.Kill(syscall.Getppid(), syscall.SIGHUP)
    rw.Write([]byte("restarted"))
}
func main() {
    flag.Parse()
    mux := bone.New()
    // Custom 404
    mux.NotFoundFunc(Handler404)
    // Handle with any http method, Handle takes http.Handler as argument.
    mux.Handle("/index", http.HandlerFunc(homeHandler))
    mux.Handle("/index/:var/info/:test", http.HandlerFunc(varHandler))
    // Get, Post etc... takes http.HandlerFunc as argument.
    mux.Post("/home", http.HandlerFunc(homeHandler))
    mux.Get("/home/:var", http.HandlerFunc(varHandler))
    mux.GetFunc("/test/*", func(rw http.ResponseWriter, req *http.Request) {
        rw.Write([]byte(req.RequestURI))
    })
    mux.Get("/restart", http.HandlerFunc(restartHandler))
    err := endless.ListenAndServe(":4242", mux)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("Server on 4242 stopped")
    os.Exit(0)
}
