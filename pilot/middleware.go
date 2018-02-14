package haedal

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Middleware func(next HandlerFunc) HandlerFunc

func LogHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		next(c)

		log.Printf("[%s] %q %v\n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t),
		)
	}
}

func RecoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(
					c.ResponseWriter,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()
		next(c)
	}
}

func StaticMiddleware(dir string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		var (
			dir       = http.Dir(dir)
			indexFile = "index.html"
		)

		return func(c *Context) {
			if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
				next(c)
				return
			}

			file := c.Request.URL.Path
			f, err := dir.Open(file)
			if err != nil {
				fmt.Println(err)
				next(c)
				return
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				next(c)
				return
			}

			if fi.IsDir() {
				if !strings.HasSuffix(c.Request.URL.Path, "/") {
					http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
					return
				}

				file = path.Join(file, indexFile)

				f, err := dir.Open(file)
				if err != nil {
					next(c)
					return
				}
				defer f.Close()

				fi, err = f.Stat()
				if err != nil || fi.IsDir() {
					next(c)
					return
				}
			}

			http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
		}
	}
}
