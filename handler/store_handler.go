package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	mygin "github.com/wmsx/pkg/gin"
	"github.com/wmsx/store_api/setting"
	proto "github.com/wmsx/store_svc/proto/store"
	"math/rand"
	"mime/multipart"
	"path"
	"strconv"
	"time"
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
		header *multipart.FileHeader
		err    error
	)
	app := mygin.Gin{C: c}

	if header, err = c.FormFile("avatar"); err != nil {
		app.LogicErrorResponse("参数错误")
		return
	}

	var file multipart.File
	if file, err = header.Open(); err != nil {
		app.LogicErrorResponse("参数错误")
		return
	}
	defer file.Close()


	mengerId, err := strconv.ParseInt(c.GetHeader("uid"), 10, 64)
	if err != nil {
		app.ServerErrorResponse()
		return
	}

	objectName := fmt.Sprintf("%d_%d%s", mengerId, time.Now().Unix(), path.Ext(header.Filename))
	if err = UploadFile(c, AvatarBulk, objectName, header.Size, file); err != nil {
		log.Error("上传头像失败 err: ", err)
		app.ServerErrorResponse()
		return
	}
	log.Info("上传成功 filename: ", header.Filename, " size: ", header.Size)

	app.Response(fmt.Sprintf("%s/%s/%s", setting.MinIOSetting.Endpoint, AvatarBulk, objectName))
	return
}

func (s *StoreHandler) UploadFiles(c *gin.Context) {
	var (
		err              error
		formData         *multipart.Form
		saveStoreInfoRes *proto.SaveStoreInfoResponse
	)
	app := mygin.Gin{C: c}

	if formData, err = c.MultipartForm(); err != nil {
		app.LogicErrorResponse("参数错误")
		return
	}

	category := formData.Value["category"][0]
	headers := formData.File["files"]

	mengerId, err := strconv.ParseInt(c.GetHeader("uid"), 10, 64)
	if err != nil {
		app.ServerErrorResponse()
		return
	}

	var storeInfos []*proto.StoreInfo
	for _, header := range headers {
		var file multipart.File
		if file, err = header.Open(); err != nil {
			app.LogicErrorResponse("参数错误")
			return
		}

		randInt := rand.Int()
		objectName := fmt.Sprintf("%d/%d_%d%s", randInt, mengerId, time.Now().Unix(), path.Ext(header.Filename))
		if err = UploadFile(c, category, objectName, header.Size, file); err != nil {
			log.Error("上传文件失败 err: ", err)
			app.ServerErrorResponse()
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
	if saveStoreInfoRes, err = s.storeClient.SaveStoreInfo(c, saveStoreInfoRequest); err != nil {
		log.Error("调用SaveStoreInfo接口失败 err: ", err)
		app.ServerErrorResponse()
		return
	}
	app.Response(saveStoreInfoRes.Name2IdMap)
	return
}