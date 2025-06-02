package rpc

import (
	"IM/rpc/group/group"
	"IM/rpc/notify/notification"
	"IM/rpc/user/user"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	EtcdHost     = []string{"etcd:2379"}
	UserClient   user.UserServiceClient
	GroupClient  group.GroupServiceClient
	NotifyClient notification.NotificationServiceClient
)

func init() {
	InitUserClient()
	InitGroupClient()
	InitNotifyClient()
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

func InitNotifyClient() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: EtcdHost,
			Key:   "notify.rpc",
		},
	})
	NotifyClient = notification.NewNotificationServiceClient(client.Conn())
}
