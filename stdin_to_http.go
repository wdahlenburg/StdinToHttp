package StdinToHttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ReadStdin(reader io.Reader, tls bool) (*http.Request, error) {
	scanner := bufio.NewReader(reader)

	// Parse the request line
	requestLine, err := scanner.ReadString('\n')
	if err != nil {
		return nil, err
	}
	requestLine = strings.TrimSpace(requestLine)

	// Parse the method, URL, and protocol from the request line
	parts := bytes.Split([]byte(requestLine), []byte(" "))
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid HTTP Request")
	}
	method := string(parts[0])
	uri := string(parts[1])
	protocol := string(parts[2])

	// Parse the HTTP headers from the file
	headers, err := readHeaders(scanner)
	if err != nil {
		return nil, err
	}

	body := bytes.Buffer{}
	_, err = body.ReadFrom(scanner)
	if err != nil {
		return nil, err
	}

	// Strip trailing newline
	body = *bytes.NewBuffer(bytes.TrimSuffix(body.Bytes(), []byte("\n")))

	request, err := http.NewRequest(method, uri, &body)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	request.Proto = protocol
	request.Host = request.Header.Get("Host")

	// Fix up the url since only the uri is set
	//https://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go
	var newUrl string
	if tls {
		newUrl = fmt.Sprintf("https://%s%s", request.Host, uri)
	} else {
		newUrl = fmt.Sprintf("http://%s%s", request.Host, uri)
	}
	u, err := url.Parse(newUrl)
	if err != nil {
		return nil, err
	}
	request.URL = u
	request.RequestURI = ""

	return request, nil
}

func readHeaders(reader *bufio.Reader) (http.Header, error) {
	var headers http.Header = http.Header{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		// Split the line into a key and value
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Malformed header: %q", line)
		}

		// Add the header to the http.Header object
		key := http.CanonicalHeaderKey(parts[0])
		value := strings.TrimSpace(parts[1])
		headers.Add(key, value)
	}

	return headers, nil
}
