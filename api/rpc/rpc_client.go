package rpc

import (
	"IM/rpc/friend/friend"
	"IM/rpc/group/group"
	"IM/rpc/user/user"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	EtcdHost     = []string{"etcd:2379"}
	UserClient   user.UserServiceClient
	GroupClient  group.GroupServiceClient
	FriendClient friend.FriendServiceClient
)

func init() {
	InitUserClient()
	InitGroupClient()
	InitFriendClient()
}

func InitUserClient() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: EtcdHost,
			Key:   "user.rpc",
		},
	})
	UserClient = user.NewUserServiceClient(client.Conn())
}

func InitGroupClient() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: EtcdHost,
			Key:   "group.rpc",
		},
	})
	GroupClient = group.NewGroupServiceClient(client.Conn())
}

func InitFriendClient() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: EtcdHost,
			Key:   "friend.rpc",
		},
	})
	FriendClient = friend.NewFriendServiceClient(client.Conn())
}
