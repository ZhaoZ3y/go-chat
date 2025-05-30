package rpc

import (
	"IM/rpc/user/user"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	EtcdHost   = []string{"localhost:2379"}
	UserClient user.UserServiceClient
)

func InitUserClient() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: EtcdHost,
			Key:   "user.rpc",
		},
	})
	UserClient = user.NewUserServiceClient(client.Conn())
}
