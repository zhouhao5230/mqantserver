package main

import (
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/module/modules"
	//"github.com/liangdas/mqant-modules/tracing"
	_ "net/http/pprof"
	"net/http"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/nats-io/go-nats"
	"sync"
	"github.com/liangdas/mqant/selector"
	"math/rand"
	"fmt"
	"mqantserver/src/server/helloworld"
	"mqantserver/src/server/login"
	"mqantserver/src/server/chat"
	"mqantserver/src/server/user"
	"mqantserver/src/server/hitball"
	"mqantserver/src/webapp"
	"mqantserver/src/server/xaxb"
	"mqantserver/src/server/gate"
)
//func ChatRoute( app module.App,Type string,hash string) (*module.ServerSession){
//	//演示多个服务路由 默认使用第一个Server
//	log.Debug("Hash:%s 将要调用 type : %s",hash,Type)
//	servers:=app.GetServersByType(Type)
//	if len(servers)==0{
//		return nil
//	}
//	return servers[0]
//}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	rs:=registry.DefaultRegistry //etcdv3.NewRegistry()
	nc, err := nats.Connect(nats.DefaultURL,nats.MaxReconnects(10000))
	if err != nil {

	}
	app := mqant.CreateApp(
		module.Nats(nc),
		module.Registry(rs),
	)
	app.Options().Selector.Init(selector.SetStrategy(func(services []*registry.Service) selector.Next{
		var nodes []*registry.Node

		// Filter the nodes for datacenter
		for _, service := range services {
			for _, node := range service.Nodes {
				nodes = append(nodes, node)
				//if node.Metadata["type"] == "helloworld" {
				//	nodes = append(nodes, node)
				//}
			}
		}

		var mtx sync.Mutex
		//log.Info("services[0] $v",services[0].Nodes[0])
		return func() (*registry.Node, error) {
			mtx.Lock()
			defer mtx.Unlock()
			if len(nodes)==0{
				return nil,fmt.Errorf("no node")
			}
			index := rand.Intn(int(len(nodes)))
			return nodes[index], nil
		}
	}))
	//app.Route("Chat",ChatRoute)
	app.Run(true, //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		modules.MasterModule(),
		hitball.Module(),
		mgate.Module(),  //这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		helloworld.Module(),
		login.Module(), //这是用户登录验证模块
		chat.Module(),
		user.Module(),
		webapp.Module(),
		xaxb.Module(),
		//tracing.Module(), //很多初学者不会改文件路径，先移除了
	)  //这是聊天模块

}

