package network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/pkg/errors"
)

func Post(url string, headers map[string]string, data url.Values) ([]byte, error) {

	var req *http.Request
	var err error

	if data == nil {
		req, err = http.NewRequest("POST", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("data:%v", data))
	}

	//增加header选项
	for k,v := range headers {
		req.Header.Set(k, v)
	}

	//发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("req:%v", req))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	//处理返回结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("resp.Body:%v", resp.Body))
	}

	return body, nil
}

func PostJson(url string, headers map[string]string, data []byte) ([]byte,error) {

	var req *http.Request
	var err error

	if data == nil {
		req, err = http.NewRequest("POST", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("data:%v", data))
	}

	//增加header选项
	req.Header.Set("Content-Type", "application/json")
	for k,v := range headers {
		req.Header.Set(k, v)
	}

	//发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("req:%v", req))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	//处理返回结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("resp.Body:%v", resp.Body))
	}

	return body, nil
}

//GET请求构造
func Get(url string, headers map[string]string) ([]byte, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("url:%s", url))
	}

	//增加header选项
	for k,v := range headers {
		req.Header.Add(k, v)
	}

	//发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("req:%v", req))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	//处理返回结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("resp.Body:%v", resp.Body))
	}

	return body, nil
}
