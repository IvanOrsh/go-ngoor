```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
	bio := `<script>alert("Should never allow such things!");</script>`

  // ...
}
```

using encoded values:

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
	bio := `&lt;script&gt;alert(&quot;Hi!&quot;)&lt;/script&gt;`

  // ...
}
```

- "html/template" package would do that for us!
