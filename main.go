package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const faqHtml = `
<h1>FAQ Page</h1>
<p><b>Q</b>: Is there a free version?</p>
<p><b>A</b>: Yes! We offer a free trial for 30 days on any paid plans.</p>
<br />
<p><b>Q</b>: What are your support hours?</p>
<p><b>A</b>: We have support staff answering email 24/7, though response times may be a bit slower on weekdays.</p>
<br />
<p><b>Q</b>: How do I contact support?</p>
<p><b>A</b>: Email us - support@lenslocked.com.</p>
`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // redundant
	w.WriteHeader(http.StatusOK) // redundant
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:john@example.com\">john</a></p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, faqHtml)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}