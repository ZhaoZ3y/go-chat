package controller

import (
	"IM/api/rpc"
	"IM/pkg/utils/response"
	"IM/rpc/file/file"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strconv"
	"time"
)

// UploadFile 处理文件上传请求。
func UploadFile(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	formFile, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ClientErrorResponse(c, response.ParamErrorCode, "获取上传文件失败: "+err.Error())
		return
	}
	defer formFile.Close()

	fileData, err := io.ReadAll(formFile)
	if err != nil {
		logx.Errorf("读取上传的文件数据失败: %v", err)
		response.ServerErrorResponse(c, "读取文件内容失败")
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rpcResp, err := rpc.FileClient.UploadFile(rpcCtx, &file.UploadFileRequest{
		UserId:      userID,
		FileName:    header.Filename,
		FileSize:    header.Size,
		ContentType: header.Header.Get("Content-Type"),
		FileData:    fileData,
	})
	if err != nil {
		logx.Errorf("RPC UploadFile 调用失败: %v", err)
		response.ServerErrorResponse(c, "上传文件失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}

// GetFileInfo 获取特定文件的元数据信息。
func GetFileInfo(c *gin.Context) {
	_, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	fileID := c.Query("file_id")
	if fileID == "" {
		response.ClientErrorResponse(c, response.ParamErrorCode, "文件ID不能为空")
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FileClient.GetFileInfo(rpcCtx, &file.GetFileInfoRequest{
		FileId: fileID,
	})
	if err != nil {
		logx.Errorf("RPC GetFileInfo 调用失败, file_id %s: %v", fileID, err)
		response.ServerErrorResponse(c, "获取文件信息失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}

// DownloadFile 提供文件下载服务。
func DownloadFile(c *gin.Context) {
	_, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	fileID := c.Query("file_id")
	if fileID == "" {
		response.ClientErrorResponse(c, response.ParamErrorCode, "文件ID不能为空")
		return
	}

	// 为下载设置一个更长的超时时间
	rpcCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rpcResp, err := rpc.FileClient.DownloadFile(rpcCtx, &file.DownloadFileRequest{
		FileId: fileID,
	})
	if err != nil {
		logx.Errorf("RPC DownloadFile 调用失败, file_id %s: %v", fileID, err)
		response.ServerErrorResponse(c, "下载文件失败")
		return
	}

	// 设置 HTTP 响应头，以提示浏览器下载文件
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", rpcResp.FileName))
	c.Header("Content-Type", rpcResp.ContentType)
	c.Header("Content-Length", strconv.FormatInt(rpcResp.FileSize, 10))

	// 将文件的字节数据直接写入 HTTP 响应体
	c.Data(http.StatusOK, rpcResp.ContentType, rpcResp.FileData)
}

// DeleteFile 删除一个文件记录（软删除）。
func DeleteFile(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}

	fileID := c.Query("file_id")
	if fileID == "" {
		response.ClientErrorResponse(c, response.ParamErrorCode, "文件ID不能为空")
		return
	}

	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 将用户ID传递给 RPC 服务以进行权限校验
	rpcResp, err := rpc.FileClient.DeleteFile(rpcCtx, &file.DeleteFileRequest{
		FileId: fileID,
		UserId: userID,
	})
	if err != nil {
		logx.Errorf("RPC DeleteFile 调用失败, file_id %s: %v", fileID, err)
		response.ServerErrorResponse(c, "删除文件失败")
		return
	}
	if !rpcResp.Success {
		response.ClientErrorResponse(c, response.RPCClientErrorCode, "删除文件操作失败")
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "文件删除成功",
	})
}

// GetFileRecord 获取用户的文件记录列表。
func GetFileRecord(c *gin.Context) {
	userID, ok := getAndParseUserID(c)
	if !ok {
		return
	}
	rpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rpcResp, err := rpc.FileClient.GetFileRecord(rpcCtx, &file.GetFileRecordReq{
		UserId: userID,
	})
	if err != nil {
		logx.Errorf("RPC GetFileRecord 调用失败: %v", err)
		response.ServerErrorResponse(c, "获取文件记录失败")
		return
	}

	response.SuccessResponse(c, rpcResp)
}
