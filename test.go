package main

import (
	"fmt"
	"github.com/go-redis/redis"
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
// 对于订阅功能，仅仅是将客户端和频道名(组)记录在某个数据结构中，
// 当有其他客户端向某个频道执行发布功能时，检查数据结构中那些订阅了该频道的客户端，
// 并向其发送消息。并不涉及到数据的存储
// 说明：在redis集群服务端中，其中的一台机器会向集群中所有的机器无脑地广播发布订阅信息
// 而且在redis集群中，每一台机器都会维护着集群中所有机器的fd等信息，可以参见redis源码
func client_subscribe(client *redis.ClusterClient) {
	// 返回值 sub 类似于一个处于订阅模式下的客户端
	// 处于订阅模式下的客户端，只能向Redis服务器发送PING、SUBSCRIBE、UNSUBSCRIBE、PSUBSCRIBE、PUNSUBSCRIBE命令
	sub := client.Subscribe("mychannel1", "mychannel2")
	sub.Subscribe("mychannel3") // 订阅某一个频道
	sub.Unsubscribe("mychannel2") // 退订某一个频道
	for {
		msg, err := sub.ReceiveTimeout(0)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("receive msg:", msg)
	}
}

func main() {
	//g_redis := redis.NewClient(&redis.Options{
	//	Addr:        "127.0.0.1" + ":" + strconv.Itoa(6379),
	//	Password:    "",
	//	DB:          0,
	//	PoolSize:    1024,
	//	MaxRetries:  3,
	//	IdleTimeout: 10 * time.Minute,
	//})
	// 实现的是客户端分片，将分片工作放在业务程序端，
	// 程序代码根据预先设置的路由规则，直接对多个Redis实例进行分布式访问
	// 是一种静态分片技术。Redis实例的增减，都得手工调整分片程序。
	// 这种分片机制的性能比代理式更好（少了一个中间分发环节）。
	// 但缺点是升级麻烦，对研发人员的个人依赖性强——需要有较强的程序开发能力做后盾
	// 其实这个redis客户端不是这样做的了啊，因为所有的分片信息都是来自于redis服务端了，自身并不需要分片。
	client := redis.NewClusterClient(&redis.ClusterOptions{
		//Addrs: []string{":6379", ":6380", ":6381", ":6382", ":6383", ":6384"},
		Addrs: []string{":6379"},
	})

	// 在新建一个集群客户端时，客户端会自动向redis集群服务端请求集群状态信息，节点信息等
	// cluster info,cluster nodes,并通过这些信息来组建客户端关于服务端集群的相关映射，
	// 那么在客户端就可以根据一致性哈希直接访问到对应的节点。
	// 同时为了检测集群状态变化导致节点不命中的问题，客户端也会在每次命令执行返回moved或者ask时，重新获取
	// 集群状态信息，并更新到最新的集群状态。
	client2 := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":6379", ":6380", ":6381", ":6382", ":6383", ":6384"},
	})
	pong, err := client2.Ping().Result()
	fmt.Println(pong, err)

	//pong, err = client.Ping().Result()
	//fmt.Println(pong, err)
	//
	//client.HSet("statistic_192.168.2.106_7780", "10003", "333")
	//hgetall := client.HGetAll("statistic_192.168.2.106_7780")
	//
	//fmt.Println(hgetall.Val())
	//fmt.Println(hgetall.Val()["10002"])

	ret := client.Get("123")
	fmt.Println(ret.Result())

	ret = client.Get("foo")
	fmt.Println(ret.Result())

	go client_subscribe(client)
	time.Sleep(1 * time.Second)
	fmt.Println("send msg:")
	client.Publish("mychannel1", "hello")
	client.Publish("mychannel2", "world.")
	client.Publish("mychannel3", "lizhanbin")

	select {}
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
