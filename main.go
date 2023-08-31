package main

import (
	"fmt"
	"net/http"
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

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 	}
// }

type Router struct {}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request ) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}