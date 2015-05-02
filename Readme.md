# JSON bind for gocraft/web package#

- decode json body to your object  
- encode your object to json and send it as a response

## Installation

    $ go get github.com/corneldamian/json-binding

## Usage

```go
package main

import (
	"net/http"

	"github.com/corneldamian/json-binding"
	"github.com/gocraft/web"
)

type Context struct {
	RequestJSON interface{}
	ResponseJSON interface{}
	ResponseStatus int
}

type Authenticate struct {
	Username string
	Password string
}

func Login(ctx *Context, rw web.ResponseWriter, req *web.Request) {
	a := ctx.RequestJSON.(*Authenticate)
	ctx.ResponseJSON = binding.SuccessResponse("User " + a.Username + " Pass: " +  a.Password)
	ctx.ResponseStatus = http.StatusUnauthorized
}

func main() {
	web := web.New(Context{}).
		Middleware(binding.Response(nil)).		
		Middleware(binding.Request(Authenticate{}, nil)).
		Post("/auth/login", Login)

	http.ListenAndServe("localhost:8080", web)
}
```


Inspiration from: http://github.com/opennota/json-binding