package utils

import (
	"sort"
	"strings"
	"github.com/valyala/fasthttp"
	"errors"
	"sync"
	"fmt"
	"time"
)

const (
	HTTP_TIMEOUT = 2
)

var (
	bufPool = &sync.Pool{New:func() interface{} {
		return make([]byte, 2 * 1024)
	}}
	//transport = http.Transport{
	//	Dial: func(network, addr string) (net.Conn, error) {
	//		deadline := time.Now().Add((HTTP_TIMEOUT + 2) * time.Second)
	//		c, err := net.DialTimeout(network, addr, HTTP_TIMEOUT*time.Second)
	//		if err != nil {
	//			return nil, err
	//		}
	//		c.SetDeadline(deadline)
	//		return c, nil
	//	},
	//	DisableKeepAlives: true,
	//}
)

type HttpClient struct {
	client *fasthttp.Client
}

func NewHttpClient() HttpClient {
	return HttpClient{
		client: &fasthttp.Client{ReadTimeout: HTTP_TIMEOUT * time.Second, WriteTimeout:HTTP_TIMEOUT * time.Second},
	}
}

//map to query string & sort by key
func httpBuildQuery(params map[string]string) string {
	list := make([]string, 0, len(params))
	buffer := make([]string, 0, len(params))
	for key := range params {
		list = append(list, key)
	}
	sort.Strings(list)
	for _, key := range list {
		value := params[key]
		buffer = append(buffer, key)
		buffer = append(buffer, "=")
		buffer = append(buffer, value)
		buffer = append(buffer, "&")
	}
	buffer = buffer[:len(buffer) - 1]
	return strings.Join(buffer, "")
}

func (b *HttpClient) Get(url string) ([]byte, error) {
	buf, _ := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	code, body, err := b.client.Get(buf, url)
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, errors.New(fmt.Sprintf("server return code %d", code))
	}


	return body, nil
}
