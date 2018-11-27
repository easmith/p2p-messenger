package types

import "fmt"

type Response struct {
	Status string
	Headers map[string]string
	Body string
}

func New(body string) (response Response) {
	response.Status = "HTTP/1.1 200 OK"
	response.Headers = make(map[string]string)
	response.Body = body
	return
}

func (r Response) String() string{
	result := r.Status + "\r\n"

	r.Headers["Content-length"] = fmt.Sprintf("%v", len([]byte(r.Body)))

	for k, v  := range r.Headers {
		result += fmt.Sprintf("%s: %s\r\n", k, v)
	}



	result += "\r\n" + r.Body
	return result
}

