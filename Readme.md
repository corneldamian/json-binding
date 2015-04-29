# JSON bind for gocraft/web package#

## Installation

    $ go get github.com/corneldamian/json-binding

## Usage

```go
package main

import (
	"net/http"
	"fmt"

	"github.com/corneldamian/json-binding"
	"github.com/gocraft/web"
)

type Context struct {
	BodyJSON interface{}
}

type Authenticate struct {
	Username string
	Password string
}

func Login(ctx *Context, rw web.ResponseWriter, req *web.Request) {
	a := ctx.BodyJSON.(*Authenticate)
	fmt.Fprintf(rw, "User %s, Pass: %s", a.Username, a.Password)
}

func main() {
	web := web.New(Context{}).
		Middleware(binding.Bind(Authenticate{}, nil)).
		Post("/auth/login", Login)

	http.ListenAndServe("localhost:8080", web)
}
```


Inspiration from: http://github.com/opennota/json-binding