// 2.写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。
package main

import (
	"bytes"
	"flag"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-redis/redis"
	"strings"
)

var (
	Size = flag.Int("size", 10, "")
)

func main() {
	flag.Parse()

	monitor()
	write()
	monitor()
}

func write() {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint("localhost:6379"),
	})
	defer client.Close()
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	// 一共写1万个key
	n := 10000

	//制作数据
	var b bytes.Buffer
	for i := 0; i < *Size; i++ {
		b.Write([]byte("1"))
	}
	value := b.String()

	for i := 0; i < n; i++ {
		if err := client.Set(fmt.Sprintf("test-%d-%d", *Size, i), value, 0).Err(); err != nil {
			panic(err)
		}
	}
}

func monitor() {

	monitor, err := redigo.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer monitor.Close()

	//查看info memory
	monitor.Send("MULTI")
	monitor.Send("info")
	info, err := redigo.Strings(monitor.Do("EXEC"))
	if err != nil {
		panic(err)
	}

	infoSlice := strings.Split(info[0], "\n\r")
	for i, v := range infoSlice {
		if i == 2 {
			vSlice := strings.Split(v, "\n")
			for _, vv := range vSlice {
				vvSlice := strings.Split(vv, ":")
				if vvSlice[0] == "used_memory_dataset" {
					fmt.Println(vv)
				}
			}
		}
	}

}
