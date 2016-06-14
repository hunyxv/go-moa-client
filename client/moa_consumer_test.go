package client

import (
	"github.com/blackbeans/go-moa/proxy"
	// "runtime"

	"github.com/blackbeans/go-moa/core"
	"sync"
	"testing"
	"time"
)

var consumer *MoaConsumer

func init() {

	demo := Demo{make(map[string][]string, 2), "/service/lookup"}
	inter := (*IHello)(nil)
	uinter := (*IUserService)(nil)
	core.NewApplcation("../conf/moa_server.toml", func() []proxy.Service {
		return []proxy.Service{
			proxy.Service{
				ServiceUri: "/service/lookup",
				Instance:   demo,
				Interface:  inter},
			proxy.Service{
				ServiceUri: "/service/moa-admin",
				Instance:   demo,
				Interface:  inter},
			proxy.Service{
				ServiceUri: "/service/user-service",
				Instance:   UserServiceDemo{},
				Interface:  uinter},
			proxy.Service{
				ServiceUri: "/service/user-service-panic",
				Instance:   UserServicePanic{},
				Interface:  uinter}}
	})

	consumer = NewMoaConsumer("../conf/moa_client.toml",
		[]proxy.Service{proxy.Service{
			ServiceUri: "/service/user-service",
			Interface:  &UserService{}}})

}

func TestMakeRpcFunc(t *testing.T) {

	//等待5s注册地址
	time.Sleep(1 * time.Second)

	consumer := NewMoaConsumer("../conf/moa_client.toml",
		[]proxy.Service{proxy.Service{
			ServiceUri: "/service/user-service",
			Interface:  &UserService{}},
			proxy.Service{
				ServiceUri: "/service/user-service-panic",
				Interface:  &UserService{}}})
	time.Sleep(2 * time.Second)
	h := consumer.GetService("/service/user-service").(*UserService)
	a, err := h.GetName("a")
	t.Logf("--------Hello,Buddy|%s|%s\n", a, err)
	if nil != err || a.Uri != "/service/user-service" {
		t.Fail()
	}

	// ---------no return
	h.SetName("a")
	//----no args
	err = h.Ping()
	t.Logf("--------Ping|%s\n", err)
	if nil != err {
		t.Fail()
	}

	_, err = h.Pong()
	t.Logf("--------Pong|%s\n", err)
	if nil != err {
		t.Fail()
	}

	h = consumer.GetService("/service/user-service-panic").(*UserService)
	a, err = h.GetName("a")
	t.Logf("--------Hello,Buddy|%s|error(%s)\n", a, err)
	if nil == err || nil != a {
		t.Fail()
	}
	//---------no return
	h.SetName("a")

	// 暂停一下，不然moa-stat统计打印不出来
	time.Sleep(time.Second * 2)

	consumer.Destory()

}

func TestConsumerPing(t *testing.T) {

	//等待5s注册地址
	time.Sleep(2 * time.Second)

	consumer := NewMoaConsumer("../conf/moa_client.toml",
		[]proxy.Service{proxy.Service{
			ServiceUri: "/service/user-service",
			Interface:  &UserService{}},
			proxy.Service{
				ServiceUri: "/service/user-service-panic",
				Interface:  &UserService{}}})
	//等待5s注册地址
	time.Sleep(5 * time.Second)

	wg := sync.WaitGroup{}
	clone := consumer.clientManager.clientManager.ClientsClone()
	wg.Add(len(clone))
	//等待空闲
	time.Sleep(10 * time.Second)
	for _, c := range clone {
		tmp := c
		go func() {
			defer wg.Done()
			succ := consumer.clientManager.ping(tmp)
			consumer.clientManager.ReleaseClient(tmp)
			t.Logf("[%s] ping %v", c.LocalAddr(), succ)
			if !succ {
				t.Fail()
			}
		}()
	}
	//等待本次的所有的PING—PONG结束
	wg.Wait()
	consumer.Destory()
}

func BenchmarkParallerMakeRpcFunc(b *testing.B) {

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h := consumer.GetService("/service/user-service").(*UserService)
			b.Logf("----------start--------\n")
			a, err := h.GetName("a")
			b.Logf("----------end--------%s\n", a)
			if nil != err || a.Uri != "/service/user-service" {
				b.Fail()
			}
		}
	})

}

func BenchmarkMakeRpcFunc(b *testing.B) {

	b.StopTimer()
	consumer := NewMoaConsumer("../conf/moa_client.toml",
		[]proxy.Service{proxy.Service{
			ServiceUri: "/service/user-service",
			Interface:  &UserService{}}})
	h := consumer.GetService("/service/user-service").(*UserService)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.GetName("a")
		// a, _ := h.GetName("a")
		// b.Logf("--------Hello,Buddy|%s\n", a)
		// if a.Uri != "/service/user-service" {
		// 	b.Fail()
		// }
	}
}
