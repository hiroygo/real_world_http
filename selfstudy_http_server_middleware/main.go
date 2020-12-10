// https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
// https://tutuz-tech.hatenablog.com/entry/2020/03/23/220326
// https://journal.lampetty.net/entry/implementing-middleware-with-http-package-in-go
// https://blog.lufia.org/entry/2016/08/28/000000
// https://mattn.kaoriya.net/software/lang/go/20160713120926.htm

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// あえて勉強用に書いてるけど、実際に制限するときは
// netutil.LimitedListener がいいはず
type visitorLimit struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type visitorLimitManager struct {
	mu     sync.Mutex
	limits map[string]*visitorLimit
}

func (v *visitorLimitManager) get(ip string) *rate.Limiter {
	v.mu.Lock()
	defer v.mu.Unlock()

	vl, ok := v.limits[ip]
	if ok {
		vl.lastSeen = time.Now()
		return vl.limiter
	}

	lim := rate.NewLimiter(1, 2)
	v.limits[ip] = &visitorLimit{limiter: lim, lastSeen: time.Now()}
	return lim
}

func (v *visitorLimitManager) release(d time.Duration) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for ip, vl := range v.limits {
		if time.Since(vl.lastSeen) > d {
			delete(v.limits, ip)
		}
	}
}

var vlManager = &visitorLimitManager{limits: make(map[string]*visitorLimit)}

func limitRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		limiter := vlManager.get(ip)
		if !limiter.Allow() {
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func requireId(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t, err := r.Cookie("token")
		if err != nil || t.Value == "" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func main() {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				vlManager.release(1 * time.Minute)
			}
		}
	}()

	http.Handle("/", limitRequest(requireId(http.HandlerFunc(handleHello))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
