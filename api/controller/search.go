package controller

import (
	"IM/api/rpc"
	"IM/pkg/model/request"
	"IM/pkg/utils/response"
	"IM/rpc/user/user"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// Search 搜索用户和群组
func Search(c *gin.Context) {
	var req request.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "参数错误: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UserResp, err := rpc.UserClient.SearchUser(ctx, &user.SearchUserRequest{
		Keyword: req.Keyword,
	})
	if err != nil {
		logx.Errorf("搜索用户失败: %v", err)
		response.ClientErrorResponse(c, response.ParamErrorCode, "搜索用户失败: "+err.Error())
		return
	}

	//GroupResp, err := rpc.GroupClient.SearchGroup(ctx, &group.SearchGroupRequest{
	//	Keyword: req.Keyword,
	//})
	//if err != nil {
	//	logx.Errorf("搜索群组失败: %v", err)
	//	response.ClientErrorResponse(c, response.ParamErrorCode, "搜索群组失败: "+err.Error())
	//	return
	//}

	response.SuccessResponse(c, gin.H{
		"users": UserResp.Users,
		//"groups":      GroupResp.Groups,
		"totalUsers": UserResp.Total,
		//"totalGroups": GroupResp.Total,
	})
}
