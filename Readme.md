# StdinToHttp

This small library will convert a Golang io.Reader to a net/http Request. This allows for the ease of storing raw HTTP requests and allowing native Golang to do the work of converting them into a Request object.

### Usage

```
package main

import (
	"fmt"
	"github.com/wdahlenburg/StdinToHttp"
	"log"
	"net/http"
	"os"
)

func main() {
	request, err := StdinToHttp.ReadStdin(os.Stdin, true)
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}

	fmt.Printf("[%d]\n", resp.StatusCode)
}
```