package utils

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type Client struct {
	base http.Client
}

func NewClient(c http.Client) *Client {
	return &Client{
		base: c,
	}
}

func (c *Client) Post(ctx context.Context, url string, req []byte) ([]byte, error) {

	reqBody := bytes.NewBuffer(req)

	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	addRequestIDToHeader(ctx, request)

	return executeRequest(c, request)
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url,
		http.NoBody)
	if err != nil {
		return nil, err
	}

	addRequestIDToHeader(ctx, req)

	return executeRequest(c, req)
}

func addRequestIDToHeader(ctx context.Context, r *http.Request) {

	requestID := middleware.GetReqID(ctx)

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(middleware.RequestIDHeader, requestID)
}

func executeRequest(c *Client, r *http.Request) ([]byte, error) {
	resp, err := c.base.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http request failed with status %v, error: %v", resp.StatusCode, string(body))
	}

	return body, nil
}

//import (
//"bytes"
//"errors"
//"fmt"
//"github.com/ugorji/go/codec"
//"io/ioutil"
//"log"
//"net/http"
//)
//
//var authErr = errors.New("authentication failed. request is not authorized")
//var jsonHandle codec.JsonHandle
//
//type Client struct {
//	authToken string
//}
//
//func NewClient() *Client {
//	return &Client{authToken: ""}
//}
//
//func (c *Client) SetAuthToken(token string) {
//	log.Println("set auth-token to the http-client")
//	c.authToken = token
//}
//
//func (c *Client) SendGetRequest(url string, queryParam map[string]string) (body []byte, httpStatus int, err error) {
//	req, err := c.createGetRequest(url, queryParam)
//	if err != nil {
//		return nil, 0, err
//	}
//	log.Println("sending GET request to url " + url)
//	res, err := c.sendRequest(req)
//	if err != nil {
//		return nil, 0, err
//	}
//	defer res.Body.Close()
//
//	resBody, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return resBody, res.StatusCode, nil
//}
//
//func (c *Client) SendPostRequest(reqBody interface{}, url string) (body []byte, httpStatus int, err error) {
//	req, err := c.createPostRequest(reqBody, url)
//	if err != nil {
//		return nil, 0, err
//	}
//	log.Println("sending POST request to url " + url)
//	res, err := c.sendRequest(req)
//	if err != nil {
//		return nil, 0, err
//	}
//	defer res.Body.Close()
//
//	resBody, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return resBody, res.StatusCode, nil
//}
//func (c *Client) createPostRequest(body interface{}, url string) (*http.Request, error) {
//	jsonBody := []byte{}
//	enc := codec.NewEncoderBytes(&jsonBody, &jsonHandle)
//	err := enc.Encode(body)
//	if err != nil {
//		return nil, err
//	}
//
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
//	if err != nil {
//		return nil, err
//	}
//	if len(c.authToken) > 0 {
//		req.Header.Set("authorization", "Bearer "+c.authToken)
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//
//	return req, nil
//}
//
//func (c *Client) createGetRequest(url string, queryParam map[string]string) (*http.Request, error) {
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	if len(c.authToken) > 0 {
//		req.Header.Set("authorization", "Bearer "+c.authToken)
//	}
//
//	if len(queryParam) > 0 {
//		q := req.URL.Query()
//
//		for key, value := range queryParam {
//			q.Add(key, value)
//		}
//		req.URL.RawQuery = q.Encode()
//	}
//
//	return req, nil
//}
//func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
//	client := http.Client{}
//	return client.Do(req)
//}
//
//func CheckForSuccessfulResponse(httpStatus int) (err error) {
//	if httpStatus < 200 || httpStatus > 299 {
//		if httpStatus == 401 {
//			err = authErr
//		} else {
//			err = fmt.Errorf("http response with status code %d", httpStatus)
//		}
//
//		log.Printf("error on http-response, err: %v", err)
//		return err
//	}
//
//	return nil
//}
