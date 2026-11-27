package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// start an S3 multipart upload: POST /key?uploads
// this is the only path that writes the shared a.uploadids map
func remote_upload(remote string, body io.Reader) error {
	req, err := http.NewRequest("POST", remote, body)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("remote_upload: wrong status code %d", resp.StatusCode))
	}
	ioutil.ReadAll(resp.Body) // drain so the connection can be reused
	return nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	reqs := make(chan string, 20000)
	resp := make(chan bool, 20000)
	fmt.Println("starting upload thrasher")

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	// 32 concurrent workers
	for i := 0; i < 32; i++ {
		go func() {
			for key := range reqs {
				if err := remote_upload("http://localhost:3000/"+key+"?uploads", strings.NewReader("")); err != nil {
					fmt.Println("UPLOAD FAILED", err)
					resp <- false
					continue
				}
				resp <- true
			}
		}()
	}

	count := 20000

	start := time.Now()
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("upload-%d", rand.Int())
		reqs <- key
	}

	for i := 0; i < count; i++ {
		if <-resp == false {
			fmt.Println("ERROR on", i)
			os.Exit(-1)
		}
	}

	fmt.Println(count, "multipart uploads in", time.Since(start))
	fmt.Printf("thats %.2f/sec\n", float32(count)/(float32(time.Since(start))/1e9))
}
