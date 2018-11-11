package connection

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

// PseudoRequest : Pseudo vesion of http.Request
type PseudoRequest struct {
	Method           string
	URL              string
	Body             []byte
	Header           http.Header
	TransferEncoding []string
	Close            bool
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
	Response         *PseudoResponse //*http.Response
}

// PseudoResponse : Pseudo version of http.Response
type PseudoResponse struct {
	Status           string // e.g. "200 OK"
	StatusCode       int    // e.g. 200
	Proto            string // e.g. "HTTP/1.0"
	ProtoMajor       int    // e.g. 1
	ProtoMinor       int    // e.g. 0
	Header           http.Header
	Body             []byte //io.ReadCloser
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Trailer          http.Header
	Request          *PseudoRequest //*http.Request
	TLS              *tls.ConnectionState
}

// RequestToPseudoRequest : Converts http.Request to PseudoRequest
func RequestToPseudoRequest(req *http.Request) (preq *PseudoRequest, err error) {
	preq.Method = req.Method
	preq.URL = req.URL.String()
	_, err = req.Body.Read(preq.Body)
	if err != nil {
		return nil, err
	}
	preq.Header = req.Header
	preq.TransferEncoding = req.TransferEncoding
	preq.Close = req.Close
	preq.Form = req.Form
	preq.PostForm = req.PostForm
	preq.MultipartForm = req.MultipartForm
	preq.Trailer = req.Trailer
	preq.RemoteAddr = req.RemoteAddr
	preq.RequestURI = req.RequestURI
	preq.TLS = req.TLS
	if req.Response != nil {
		pres, err := ResponseToPseudoResponse(req.Response)
		if err != nil {
			return nil, err
		}
		preq.Response = pres
	}
	return preq, nil
}

// ResponseToPseudoResponse : Converts http.Response to PseudoResponse
func ResponseToPseudoResponse(res *http.Response) (pres *PseudoResponse, err error) {
	pres.Status = res.Status
	pres.StatusCode = res.StatusCode
	pres.Proto = res.Proto
	pres.ProtoMajor = res.ProtoMajor
	pres.ProtoMinor = res.ProtoMinor
	pres.Header = res.Header
	_, err = res.Body.Read(pres.Body)
	if err != nil {
		return nil, err
	}
	pres.ContentLength = res.ContentLength
	pres.TransferEncoding = res.TransferEncoding
	pres.Close = res.Close
	pres.Uncompressed = res.Uncompressed
	pres.Trailer = res.Trailer
	if res.Request != nil {
		preq, err := RequestToPseudoRequest(res.Request)
		if err != nil {
			return nil, err
		}
		pres.Request = preq
	}
	pres.TLS = res.TLS
	return pres, nil
}

// PseudoRequestToRequest : Converts PseudoRequest to http.Request
func PseudoRequestToRequest(preq *PseudoRequest) (req *http.Request, err error) {
	bodyBuffer := bytes.NewBuffer(preq.Body)
	req, err = http.NewRequest(preq.Method, preq.URL, bodyBuffer)
	if err != nil {
		return nil, err
	}
	req.Header = preq.Header
	req.TransferEncoding = preq.TransferEncoding
	req.Close = preq.Close
	req.Form = preq.Form
	req.PostForm = preq.PostForm
	req.MultipartForm = preq.MultipartForm
	req.Trailer = preq.Trailer
	req.RemoteAddr = preq.RemoteAddr
	req.RequestURI = preq.RequestURI
	req.TLS = preq.TLS
	if preq.Response != nil {
		res, err := PseudoResponseToResponse(preq.Response)
		if err != nil {
			return nil, err
		}
		req.Response = res
	}
	return req, nil
}

// PseudoResponseToResponse : Converts PseudoResponse to http.Response
func PseudoResponseToResponse(pres *PseudoResponse) (res *http.Response, err error) {
	res.Status = pres.Status
	res.StatusCode = pres.StatusCode
	res.Proto = pres.Proto
	res.ProtoMajor = pres.ProtoMajor
	res.ProtoMinor = pres.ProtoMinor
	res.Header = pres.Header
	bodyBuffer := bytes.NewBuffer(pres.Body)
	res.Body = ioutil.NopCloser(bodyBuffer)
	res.ContentLength = pres.ContentLength
	res.TransferEncoding = pres.TransferEncoding
	res.Close = pres.Close
	res.Uncompressed = pres.Uncompressed
	res.Trailer = pres.Trailer
	if pres.Request != nil {
		req, err := PseudoRequestToRequest(pres.Request)
		if err != nil {
			return nil, err
		}
		res.Request = req
	}
	res.TLS = pres.TLS
	return res, nil
}
