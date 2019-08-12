package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"

	"trace"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type templateHandler struct {
	once       sync.Once
	filename   string
	templating *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templating = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	err := t.templating.Execute(w, data)
	if err != nil {
		log.Println("templateHandler: couldn't parse template. ", err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "The addr of the application")
	doTrace := flag.Bool("trace", false, "Indicates application requires to trace all the events to output")
	flag.Parse()

	setupOAuth()
	setupProfilePictureServer()

	r := newRoom()
	if *doTrace {
		r.tracer = trace.New(os.Stdout)
	}

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	go r.run()

	log.Println("Starting Web Server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Set up OAuth authentication with gomniauth
func setupOAuth() {
	gomniauth.SetSecurityKey("SECURITY KEY")
	gomniauth.WithProviders(
		facebook.New("1320346858115352",
			"f31b3857c34cdc27ffb59da1e83bdfe3",
			"http://blueprints-chat-app.com:8080/auth/callback/facebook",
		),
		github.New("f455ac02dd8d63bea606",
			"f1c4c0eb759c9b0a910ac92c2c7009ae68e4a85a",
			"http://blueprints-chat-app.com:8080/auth/callback/github",
		),
		google.New(
			"986332446747-bf3pl6m0noq47pmsaucdijpqm4rfgjvv.apps.googleusercontent.com",
			"vHgd_mWL0rRuR4o_hl3pqIhX",
			"http://blueprints-chat-app.com:8080/auth/callback/google",
		),
	)
}

// Set up serving user avatars via HTTP
func setupProfilePictureServer() {
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
}
