Using this data

```go
user := User{
  Name: "Bob",
  Email: "bob@example.com"
}
```

We want to render the following:

```html
<body>
  <a her="/account">bob@example</a>
  ...
  <h1>Hello, Bob!</h1>
</body>
```

## Server-side

```html
<body>
  <!-- A link to the users account details -->
  <a href="/account">{{.Email}}</a>
  <!-- ... -->
  <h1>Hello, {{.Name}}</h1>
</body>
```

## API

```json
{
  "name": "Bob",
  "email": "bob@example.com"
}
```
