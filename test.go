package main

import (
	"fmt"
	"github.com/go-redis/redis"
	//"time"
	"strconv"
	"time"
)
/*
请求协议：
	* <参数数量> CR LF
	$ <参数 1 的字节数量> CR LF
	  <参数 1 的数据> CR LF
	...
	$ <参数 N 的字节数量> CR LF
	<参数 N 的数据> CR LF
*/
func client_lpush(client *redis.ClusterClient) {
	for i := 0; i < 99999999; i++ {
		client.LPush("test_list", i)
		//time.Sleep(time.Second)
	}
}

func client_rpop(client *redis.ClusterClient) {
	for {
		pop := client.BRPop(0, "test_list")
		fmt.Println("rpop:   ", pop.Val())
	}
}
func main() {
	g_redis := redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1" + ":" + strconv.Itoa(6379),
		Password:    "",
		DB:          0,
		PoolSize:    1024,
		MaxRetries:  3,
		IdleTimeout: 10 * time.Minute,
	})
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":6379", ":6380", ":6381", ":6382", ":6383", ":6384"},
	})

	client2 := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":6379", ":6380", ":6381", ":6382", ":6383", ":6384"},
	})
	pong, err := client2.Ping().Result()
	fmt.Println(pong, err)

	pong, err = client.Ping().Result()
	fmt.Println(pong, err)

	client.HSet("statistic_192.168.2.106_7780", "10003", "333")
	hgetall := client.HGetAll("statistic_192.168.2.106_7780")

	fmt.Println(hgetall.Val())
	fmt.Println(hgetall.Val()["10002"])

	ret := client.Get("123")
	fmt.Println(ret.Result())
	//go client_lpush(client)
	//go client_rpop(client)

	//client.SAdd("play_id_6", 777, 888, 111, 222, 333)
	//val, err := client.SMembers("play_id_6").Result()
	//if err == redis.Nil {
	//	fmt.Println("play_id_6 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("play_id_6", len(val), val)
	//	fmt.Println("select 4 players: ", val[0], val[1], val[2], val[3])
	//
	//	//	用集合的话，多线程读的时候再去删除会有问题，应该使用队列，lpush,rpop
	//	// Redis直接提供的命令都是原子操作，包括lpush、rpop、blpush、brpop等。
	//	// 但多条命令一起不一定原子
	//	cnt := client.SRem("play_id_6", val[0], val[1], val[2], val[3])
	//	fmt.Println(cnt.Val())
	//
	//	pop := client.RPop("mylist")
	//	fmt.Println("rpop: ", len(pop.Val()), pop.Val())
	//
	//}
	//for {
	//	pop := client2.BRPop(0,"test_list")
	//	fmt.Println("rpop_2: ", pop.Val())
	//}
}
