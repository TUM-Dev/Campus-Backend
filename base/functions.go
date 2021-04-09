package base

import (
	"bytes"
	"github.com/divan/gorilla-xmlrpc/xml"
	"net/http"
)

func XmlRpcCall(url string, method string, args struct{ Who string }) (reply struct{ Message []string }, err error) {
	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := http.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return struct{ Message []string }{}, err
	}
	defer resp.Body.Close()

	err = xml.DecodeClientResponse(resp.Body, &reply)
	return reply, err
}
