package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	mygin "github.com/wmsx/pkg/gin"
	"github.com/wmsx/store_api/handler/minio"
	"github.com/wmsx/store_api/setting"
	proto "github.com/wmsx/store_svc/proto/store"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path"
	"time"
)

const (
	storeSvcName = "wm.sx.svc.store"
)

type StoreHandler struct {
	storeClient proto.StoreService
}

func NewStoreHandler(client client.Client) *StoreHandler {
	return &StoreHandler{
		storeClient: proto.NewStoreService(storeSvcName, client),
	}
}

func (s *StoreHandler) UploadAvatar(c *gin.Context) {
	var (
		header  *multipart.FileHeader
		session *sessions.Session
		err     error
	)
	if header, err = c.FormFile("avatar"); err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	c.PostForm("")

	var file multipart.File
	if file, err = header.Open(); err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	if session, err = mygin.GetSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	mengerId := session.Values["id"].(int64)

	objectName := fmt.Sprintf("%d_%d%s", mengerId, time.Now().Unix(), path.Ext(header.Filename))
	if err = minio.UploadFile(minio.AvatarBulk, objectName, header.Size, file); err != nil {
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	log.Info("上传成功 filename: ", header.Filename, " size: ", header.Size)

	c.String(http.StatusOK, fmt.Sprintf("%s/%s/%s", setting.MinIOSetting.Endpoint, minio.AvatarBulk, objectName))
	return
}

func (s *StoreHandler) UploadFiles(c *gin.Context) {
	var (
		err              error
		formData         *multipart.Form
		session          *sessions.Session
		saveStoreInfoRes *proto.SaveStoreInfoResponse
	)

	if formData, err = c.MultipartForm(); err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	category := formData.Value["category"][0]
	headers := formData.File["files"]

	if session, err = mygin.GetSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	mengerId := session.Values["id"].(int64)

	var storeInfos []*proto.StoreInfo
	for _, header := range headers {
		var file multipart.File
		if file, err = header.Open(); err != nil {
			c.String(http.StatusBadRequest, "参数错误")
			return
		}

		randInt := rand.Int()
		objectName := fmt.Sprintf("%d/%d_%d%s", randInt, mengerId, time.Now().Unix(), path.Ext(header.Filename))
		if err = minio.UploadFile(category, objectName, header.Size, file); err != nil {
			c.String(http.StatusInternalServerError, "服务器异常")
			return
		}
		log.Info("上传成功 filename: ", header.Filename, " size: ", header.Size)

		storeInfos = append(storeInfos, &proto.StoreInfo{
			Filename:   header.Filename,
			BulkName:   category,
			ObjectName: objectName,
			Size:       header.Size,
		})
	}

	saveStoreInfoRequest := &proto.SaveStoreInfoRequest{StoreInfos: storeInfos}
	if saveStoreInfoRes, err = s.storeClient.SaveStoreInfo(context.Background(), saveStoreInfoRequest); err != nil {
		log.Error("调用SaveStoreInfo接口失败 err: ", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	c.JSON(http.StatusOK, saveStoreInfoRes.Name2IdMap)
	return
}
