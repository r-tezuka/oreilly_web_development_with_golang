package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"

	"github.com/oreilly_web_development_with_golang/trace"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

//現在アクティブなAvatarの実装
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UserAuthAvatar,
	UseGravatar}

//temp1は１つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

//ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.temp1.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈します。
	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		facebook.New("11095562301-6gf18o1mutb0t8v2nacduc4sscq6shkl.apps.googleusercontent.com", "cHpWef3nbDcwGiBjPvr8mw9w", "http://localhost:8080/auth/callback/facebook"),
		github.New("11095562301-6gf18o1mutb0t8v2nacduc4sscq6shkl.apps.googleusercontent.com", "cHpWef3nbDcwGiBjPvr8mw9w", "http://localhost:8080/auth/callback/github"),
		google.New("11095562301-6gf18o1mutb0t8v2nacduc4sscq6shkl.apps.googleusercontent.com", "cHpWef3nbDcwGiBjPvr8mw9w", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	//チャットルームを開始します
	go r.run()
	//Webサーバーを起動します
	log.Println("Webサーバーを開始します。ポート : ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
