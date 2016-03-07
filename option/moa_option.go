package option

import (
	"errors"
	log "github.com/blackbeans/log4go"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"time"
)

type HostPort struct {
	Hosts string
}

//配置信息
type Option struct {
	Env struct {
		RunMode      string
		BindAddress  string
		RegistryType string
		AppName      string
		AppSecretKey string
	}

	//使用的环境
	Registry map[string]HostPort //momokeeper的配置
	Clusters map[string]Cluster  //各集群的配置
}

//----------------------------------------
//Cluster配置
type Cluster struct {
	Env             string //当前环境使用的是dev还是online
	ProcessTimeout  int    //处理超时 5 s单位
	PoolSizePerHost int    //5
	// ReadBufferSize   int    //=16 * 1024 //读取缓冲大小
	// WriteBufferSize  int    //=16 * 1024 //写入缓冲大小
	// WriteChannelSize int    //=1000 //写异步channel长度
	// ReadChannelSize  int    //=1000 //读异步channel长度
}

//---------最终需要的ClientCOption
type ClientOption struct {
	AppName         string
	AppSecretKey    string
	RegistryType    string
	RegistryHosts   string
	ProcessTimeout  time.Duration
	PoolSizePerHost int
	// maxDispatcherSize int           //=8000//最大分发处理协程数
	// readBufferSize    int           //=16 * 1024 //读取缓冲大小
	// writeBufferSize   int           //=16 * 1024 //写入缓冲大小
	// writeChannelSize  int           //=1000 //写异步channel长度
	// readChannelSize   int           //=1000 //读异步channel长度
	// idleDuration      time.Duration //=60s //连接空闲时间
}

func LoadConfiruation(path string) (*ClientOption, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buff, rerr := ioutil.ReadAll(f)
	if nil != rerr {
		return nil, rerr
	}
	// log.DebugLog("application", "LoadConfiruation|Parse|toml:%s", string(buff))
	//读取配置
	var option Option
	err = toml.Unmarshal(buff, &option)
	if nil != err {
		log.ErrorLog("application", "LoadConfiruation|Parse|FAIL|%s", err)
		return nil, err
	}

	cluster, ok := option.Clusters[option.Env.RunMode]
	if !ok {
		return nil, errors.New("no cluster config for " + option.Env.RunMode)
	}

	reg, exist := option.Registry[option.Env.RunMode]
	if !exist {
		return nil, errors.New("no reg  for " + option.Env.RunMode + ":" + cluster.Env)
	}

	// if cluster.MaxDispatcherSize <= 0 {
	// 	cluster.MaxDispatcherSize = 8000 //最大分发处理协程数
	// }

	// if cluster.ReadBufferSize <= 0 {
	// 	cluster.ReadBufferSize = 16 * 1024 //读取缓冲大小
	// }

	// if cluster.WriteBufferSize <= 0 {
	// 	cluster.WriteBufferSize = 16 * 1024 //写入缓冲大小
	// }

	// if cluster.WriteChannelSize <= 0 {
	// 	cluster.WriteChannelSize = 1000 //写异步channel长度
	// }

	// if cluster.ReadChannelSize <= 0 {
	// 	cluster.ReadChannelSize = 1000 //读异步channel长度

	// }

	//拼装为可用的MOA参数
	mop := &ClientOption{}
	mop.AppName = option.Env.AppName
	mop.AppSecretKey = option.Env.AppSecretKey
	mop.RegistryType = option.Env.RegistryType
	mop.RegistryHosts = reg.Hosts
	mop.ProcessTimeout = time.Duration(int64(cluster.ProcessTimeout) * int64(time.Second))
	mop.PoolSizePerHost = cluster.PoolSizePerHost
	// mop.maxDispatcherSize = cluster.MaxDispatcherSize //最大分发处理协程数
	// mop.readBufferSize = cluster.ReadBufferSize       //读取缓冲大小
	// mop.writeBufferSize = cluster.WriteBufferSize     //写入缓冲大小
	// mop.writeChannelSize = cluster.WriteChannelSize   //写异步channel长度
	// mop.readChannelSize = cluster.ReadChannelSize     //读异步channel长度
	// mop.idleDuration = 60 * time.Second
	return mop, nil

}
