package services

import "gits/internal/model/html"

const tmp = `

**A middleware handler** is simply an *http.Handler* that wraps another *http.Handler* to do some pre- and/or post-processing of the request. It's called *"middleware"* because it sits in the middle between the Go web server and the actual handler. 

In this scheme general scheme. We could see how work middleware

![image](general-scheme.png)

We could see where we need to put our code for make pre- and/or post-processing

![image](pre-andpost-processing.png)

# Logging Middleware

It is the responsibility of the console to write the details of each request. We could see go code on exapmle

![image](loggingMiddleware.png)

# Chaining Middleware

Because every middleware constructor both accepts and returns an http.Handler, you can chain multiple middleware handlers together. We could see chaining middleware using gorilla-mux library on exapmle. The handler chain will be stopped if your middleware doesn't call. See like in line codes 28-31.

![image](chainingMiddleware.png)

![image](accessMiddleware.png)

# Middleware and Request-Scoped Values
We could pass values from middle ware to handler with using context. Let's see it on code above. In line code 35 we prepare user variable in 48 line code we create new conext wrapped our user variable then in 49 line code we create reaquest with our created context

# Handling CORS Requests gorilla/mux

Gorilla mux have own realization middleware CORSMethodMiddleware. Let's see how it work we must allow options request for all api endpoints. and add middleware that configure response header Access-Control-Allow-Methods
 
![image](setupCORS.png)

Let's check cron. What are its properties? The example web browser makes the first OPTIONS request and parses the header. If we make a get request and the Access-Control-Allow-Methods does not contain GET, the web browser will not do the request.

![image](testCron.png)

`

type Publisher interface {
	Article() (*html.Article, error)
}

type publisher struct {
	md MD
}

func NewPublisher(md MD) Publisher {
	return &publisher{md: md}
}

func (p *publisher) Article() (*html.Article, error) {
	content, err := p.md.RenderMdToHTML([]byte(tmp))
	return &html.Article{
		Name:    "Gin",
		Content: string(content),
	}, err
}
