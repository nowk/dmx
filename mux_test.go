package dmx

import (
	"fmt"
	"gopkg.in/nowk/assert.v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func hfunc(s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(s))
	})
}

func TestMethodPatternDuplicationPanics(t *testing.T) {
	for _, v := range []string{
		"",
		"/posts",
	} {
		s := v
		if s == "" {
			s = "/"
		}

		{
			mux := New()
			assert.Panic(t, fmt.Sprintf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v, hfunc(""), "POST", "POST")
				})
		}
		{
			mux := New()
			assert.Panic(t, fmt.Sprintf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v+"/", hfunc(""), "POST")
					mux.Add(v, hfunc(""), "POST")
				})
		}
		{
			mux := New()
			assert.Panic(t, fmt.Sprintf("error: mux: POST %s is already defined", s),
				func() {
					mux.Add(v, hfunc(""), "POST")
					mux.Add(v+"/", hfunc(""), "POST")
				})
		}
	}
}

func TestDispatchesToMatchingResource(t *testing.T) {
	mux := New()
	mux.Add("/posts/:post_id/comments/:id", hfunc(""), "PUT", "PATCH")
	mux.Add("/posts/:post_id/comments", hfunc(""), "POST")
	mux.Add("/posts/:post_id/comments", hfunc(""), "GET")
	mux.Add("/posts/:id", hfunc(""), "PUT", "PATCH")
	mux.Add("/posts", hfunc(""), "POST")
	mux.Add("/posts", hfunc(""), "GET")
	mux.Add("/", hfunc(""), "GET")

	for k, v := range map[string][]struct {
		u string
		c int
	}{
		"GET": {
			{"/", 200},
			{"/posts", 200},
			{"/posts/123", 405},
			{"/posts/123/comments", 200},
			{"/posts/123/comments/456", 405},
			{"/posts/123/author", 404},
		},
		"POST": {
			{"/", 405},
			{"/posts", 200},
			{"/posts/123", 405},
			{"/posts/123/comments", 200},
			{"/posts/123/comments/456", 405},
			{"/posts/123/author", 404},
		},
		"PUT": {
			{"/", 405},
			{"/posts", 405},
			{"/posts/123", 200},
			{"/posts/123/comments", 405},
			{"/posts/123/comments/456", 200},
			{"/posts/123/author", 404},
		},
		"PATCH": {
			{"/", 405},
			{"/posts", 405},
			{"/posts/123", 200},
			{"/posts/123/comments", 405},
			{"/posts/123/comments/456", 200},
			{"/posts/123/author", 404},
		},
		"DELETE": {
			{"/", 405},
			{"/posts", 405},
			{"/posts/123", 405},
			{"/posts/123/comments", 405},
			{"/posts/123/comments/456", 405},
			{"/posts/123/author", 404},
		},
	} {
		for _, r := range v {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(k, fmt.Sprintf("http://www.com%s", r.u), nil)
			if err != nil {
				t.Fatal(err)
			}
			mux.ServeHTTP(w, req)

			assert.Equal(t, r.c, w.Code, k, " ", r.u)
		}
	}
}

func TestNamedParamValues(t *testing.T) {
	var ph = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		fmt.Fprintf(w, "post_id=%s&id=%s", q.Get(":post_id"), q.Get(":id"))
	})

	mux := New()
	mux.Add("/posts/:post_id/tags/:id", ph, "GET")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://www.com/posts/123/tags/456", nil)
	if err != nil {
		t.Fatal(err)
	}

	mux.ServeHTTP(w, req)
	assert.Equal(t, "post_id=123&id=456", w.Body.String())
}

func TestHandlerFuncShortcuts(t *testing.T) {
	var fn = func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello %s!", req.Method)))
	}

	mux := New()
	mux.GetFunc("/say-you", fn)
	mux.HeadFunc("/say-you", fn)
	mux.PostFunc("/say-you", fn)
	mux.PutFunc("/say-you", fn)
	mux.PatchFunc("/say-you", fn)
	mux.DelFunc("/say-you", fn)

	for _, v := range []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
	} {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(v, "/say-you", nil)
		if err != nil {
			t.Fatal(err)
		}
		mux.ServeHTTP(w, req)
		assert.Equal(t, fmt.Sprintf("Hello %s!", v), w.Body.String())
	}
}
