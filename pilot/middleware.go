package haedal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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

// TODO: need to be encapsulized
const VerifyMessage = "verified"

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			c.Redirect("/login")
			return
		} else if err != nil {
			c.RenderErr(http.StatusInternalServerError, err)
		} else if verify(VerifyMessage, v.Value) {
			next(c)
			return
		}

		c.Redirect("/login")
	}
}

func verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

func ParseFormHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				c.Params[k] = v[0]
			}
		}
		next(c)
	}
}
