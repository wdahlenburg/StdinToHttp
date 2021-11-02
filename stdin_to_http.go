package StdinToHttp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ReadStdin(reader io.Reader) (*http.Request, error) {
	scanner := bufio.NewScanner(reader)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	request, err := http.ReadRequest(bufio.NewReader(scannerToReader(scanner)))
	if err != nil {
		return nil, err
	}

	// Fix up the url since only the uri is set
	//https://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go
	u, err := url.Parse(fmt.Sprintf("http://%s%s", request.Host, request.RequestURI))
	if err != nil {
		return nil, err
	}
	request.URL = u
	request.RequestURI = ""

	return request, nil
}

func scannerToReader(scanner *bufio.Scanner) io.Reader {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		for scanner.Scan() {
			writer.Write(scanner.Bytes())
			writer.Write([]byte("\n"))
		}
	}()

	return reader
}
