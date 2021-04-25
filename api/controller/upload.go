package controller

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/tech-share/config"
	"github.com/mylxsw/tech-share/internal/service"
)

type UploadController struct {
	cc   infra.Resolver
	conf *config.Config
}

func NewUploadController(cc infra.Resolver, conf *config.Config) web.Controller {
	return &UploadController{cc: cc, conf: conf}
}

func (ctl UploadController) Register(router web.Router) {
	router.Group("upload/", func(router web.Router) {
		router.Post("/", ctl.Upload)
	})
}

// Upload upload a file
func (ctl UploadController) Upload(ctx web.Context, req web.Request, attaSrv service.AttachmentService) (*IDRes, error) {
	uploadFile, err := req.File("file")
	if err != nil {
		return nil, err
	}

	savePath := fmt.Sprintf("%s-%s.%s",
		time.Now().Format("2006/01-02/150405"),
		hash(fmt.Sprintf("%s %d", uploadFile.Name(), time.Now().Nanosecond())),
		uploadFile.Extension(),
	)

	absPath := filepath.Join(ctl.conf.StoragePath, savePath)
	_ = os.MkdirAll(filepath.Dir(absPath), os.ModePerm)

	if err := uploadFile.Store(absPath); err != nil {
		return nil, err
	}

	atta := service.Attachment{}
	atta.Name = req.InputWithDefault("name", uploadFile.Name())
	atta.AttaType = uploadFile.Extension()
	atta.AttaPath = savePath

	id, err := attaSrv.CreateAttachment(context.TODO(), atta)
	if err != nil {
		return nil, err
	}

	return &IDRes{ID: id}, nil
}

func hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
