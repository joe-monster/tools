package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const MinuteStrFormat = "2006-01-02 15:04"

var (
	Addr   = flag.String("addr", "", "format: 192.168.8.208")
	Minute = flag.String("minute", "", "format: "+MinuteStrFormat)
)

func init() {
	flag.Parse()
	if *Addr == "" || *Minute == "" {
		log.Fatalln("-h for help")
	}
}

func main() {

	deviceIds, err := getDeviceIds()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	var wg sync.WaitGroup
	wg.Add(len(deviceIds))
	for _, id := range deviceIds {
		go func(id string) {
			defer wg.Done()
			if err := makeData(id); err != nil {
				log.Printf("%+v", err)
			}
		}(id)
	}
	wg.Wait()

	os.Exit(0)
}

func getDeviceIds() ([]string, error) {
	params := url.Values{}
	params.Set("Page", "1")
	params.Set("Limit", "10000")
	params.Set("Status", "11") //11是正常的状态码

	var u url.URL
	u.Scheme = "http"
	u.Host = *Addr
	u.Path = "/api/device/v1/collector/pagelist"
	u.RawQuery = params.Encode()
	addr := u.String()

	resp, err := httpGet(addr, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	//处理状态码
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	//处理返回结果
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("resp.Body:%v", resp.Body))
	}

	type collectorList struct {
		//ID                string    `json:"Id"`
		//ProducerID        string    `json:"ProducerId"`
		//ProducerSerail    string    `json:"ProducerSerial"`
		//ProducerStatus    string    `json:"ProducerStatus"`
		Serial string `json:"Serial"`
		//OtaVersion        string    `json:"OtaVersion"`  // ota 版本
		//OtaProgress       int       `json:"OtaProgress"` // ota 进度
		//OtaSpeed          float64   `json:"OtaSpeed"`    // ota 速度
		//Node              string    `json:"Node"`
		//Client            string    `json:"Client"`
		//Location          string    `json:"Location"`
		//Status            string    `json:"Status"`
		//Type              string    `json:"Type"` // S/C/T
		//HardwareVersion   string    `json:"HardwareVersion"`
		//FirmwareVersion   string    `json:"FirmwareVersion"`
		//Time              string    `json:"Time"`
		//AudioSampleRate   string    `json:"AudioSampleRate"`
		//VibrateSampleRate string    `json:"VibrateSampleRate"`
		//CreateDatetime    time.Time `json:"CreateDatetime"`
	}
	type collectorListResponse struct {
		Total string           `json:"Total"`
		Rows  []*collectorList `json:"Rows"`
	}
	type ListCollectorResponse struct {
		Code       int                    `json:"code"`
		Msg        string                 `json:"msg"`
		ResponseAt string                 `json:"response_at"`
		Data       *collectorListResponse `json:"data"`
	}

	var data ListCollectorResponse
	if err := json.Unmarshal(byteData, &data); err != nil {
		return nil, errors.WithStack(err)
	}

	if data.Code != 200 {
		return nil, errors.New(data.Msg)
	}

	//提取设备id数据
	var ids []string
	for _, v := range data.Data.Rows {
		ids = append(ids, v.Serial)
	}

	return ids, nil
}

func makeData(id string) error {
	var u url.URL
	u.Scheme = "http"
	u.Host = *Addr
	u.Path = "/api/data/v1/history/raw"

	params := url.Values{}
	params.Set("sensorid", id)
	params.Set("type", "audio")

	t, err := time.Parse(MinuteStrFormat, *Minute)
	if err != nil {
		return errors.Wrap(err, *Minute)
	}

	params.Set("from", fmt.Sprintf("%d", t.Unix()*1e3))
	params.Set("to", fmt.Sprintf("%d", t.Add(1*time.Minute).Unix()*1e3))

	u.RawQuery = params.Encode()

	addr := u.String()

	resp, err := httpGet(addr, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	//处理状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		return errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	//处理返回结果
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("resp.Body:%v", resp.Body))
	}

	log.Println(string(byteData))

	return nil
}

//GET请求构造
func httpGet(url string, headers map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, url)
	}

	//增加header选项
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	//发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("req:%v", req))
	}

	return resp, nil
}
