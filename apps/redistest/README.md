
#### 1.使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。
测试执行1000次
```
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 10 get | grep "<="
95.90% <= 0.1 milliseconds
99.50% <= 0.2 milliseconds
100.00% <= 0.2 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 10 set | grep "<="
0.10% <= 0.1 milliseconds
89.50% <= 0.2 milliseconds
98.70% <= 0.3 milliseconds
99.70% <= 0.4 milliseconds
99.90% <= 0.5 milliseconds
100.00% <= 0.9 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 20 get | grep "<="
0.30% <= 0.1 milliseconds
80.70% <= 0.2 milliseconds
96.10% <= 0.3 milliseconds
99.20% <= 0.4 milliseconds
99.90% <= 0.5 milliseconds
100.00% <= 0.5 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 20 set | grep "<="
96.00% <= 0.1 milliseconds
99.80% <= 0.2 milliseconds
100.00% <= 0.4 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 50 get | grep "<="
0.10% <= 0.1 milliseconds
80.00% <= 0.2 milliseconds
97.40% <= 0.3 milliseconds
99.10% <= 0.4 milliseconds
99.90% <= 0.5 milliseconds
100.00% <= 0.5 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 50 set | grep "<="
98.70% <= 0.1 milliseconds
99.60% <= 0.2 milliseconds
100.00% <= 0.3 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 100 get | grep "<="
0.10% <= 0.1 milliseconds
94.00% <= 0.2 milliseconds
97.10% <= 0.3 milliseconds
99.10% <= 0.4 milliseconds
99.50% <= 0.5 milliseconds
99.60% <= 0.6 milliseconds
99.80% <= 0.8 milliseconds
100.00% <= 1.1 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 100 set | grep "<="
95.60% <= 0.1 milliseconds
99.50% <= 0.2 milliseconds
99.70% <= 0.3 milliseconds
99.90% <= 0.4 milliseconds
100.00% <= 0.4 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 200 get | grep "<="
98.80% <= 0.1 milliseconds
99.80% <= 0.2 milliseconds
100.00% <= 0.2 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 200 set | grep "<="
0.20% <= 0.1 milliseconds
78.10% <= 0.2 milliseconds
96.90% <= 0.3 milliseconds
99.00% <= 0.4 milliseconds
99.60% <= 0.5 milliseconds
99.80% <= 0.6 milliseconds
100.00% <= 0.6 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 1000 get | grep "<="
0.20% <= 0.1 milliseconds
72.00% <= 0.2 milliseconds
96.00% <= 0.3 milliseconds
99.60% <= 0.4 milliseconds
99.90% <= 0.5 milliseconds
100.00% <= 0.5 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 1000 set | grep "<="
98.60% <= 0.1 milliseconds
100.00% <= 0.3 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 5000 get | grep "<="
98.00% <= 0.1 milliseconds
99.90% <= 0.2 milliseconds
100.00% <= 0.4 milliseconds
/data # redis-benchmark -h 127.0.0.1 -p 6379 -c 1 -n 1000 -d 5000 set | grep "<="
98.40% <= 0.1 milliseconds
99.90% <= 0.3 milliseconds
100.00% <= 0.4 milliseconds
```
结论：看不出啥，就感觉1000kb的貌似相对有点慢

#### 2.写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。
写入1万个10kb的key内存前后：
```
used_memory_dataset:64024
used_memory_dataset:224024
平均每个key16kb
```
写入1万个20kb的key内存前后：
```
used_memory_dataset:224024
used_memory_dataset:624024
平均每个key40kb
```
写入1万个50kb的key内存前后：
```
used_memory_dataset:624024
used_memory_dataset:1344024
平均每个key72kb
```
写入1万个100kb的key内存前后：
```
used_memory_dataset:1344024
used_memory_dataset:2624024
平均每个key128kb
```
写入1万个200kb的key内存前后：
```
used_memory_dataset:63552
used_memory_dataset:2463552
平均每个key240kb
```
写入1万个1000kb的key内存前后：
```
used_memory_dataset:2463552
used_memory_dataset:12863552
平均每个key1040kb
```
写入1万个5000kb的key内存前后：
```
used_memory_dataset:12863552
used_memory_dataset:64223552
平均每个key5136kb
```
结论：key比较小时，额外占用空间占比较大（10kb的value实际消耗额外空间高达60%）；key越大，额外占用空间占越小。猜测应该是有一些固定元数据字段大小是固定的，只有个别元数据字段是和value大小相关。
