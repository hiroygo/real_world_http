`curl -d @main.go localhost:8080`
```
POST / HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 277
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.68.0

package mainimport (	"log"	"net/http"	"os")func main() {	file, err := os.Open("main.go")	if err != nil {		panic(err)	}	resp, err := http.Post("http://localhost:8080", "text/plain", file)	if err != nil {		panic(err)	}	defer resp.Body.Close()	log.Println("Status:", resp.Status)}
```
