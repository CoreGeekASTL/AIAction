// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"GIDS/models/browsergateway"

	"github.com/beego/beego/v2/client/orm"

	"GIDS/common/constants"
	"GIDS/common/https"
	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/req"
)

const (
	progressComplete     = 100
	percentageMultiplier = 100
	defaultRetryCount    = 2
)

type PluginService interface {
	UploadPluginPackage(req *req.UploadPluginPackageReq) error
	DeletePluginPackage(request *req.PluginPackageReq) error
	GetPluginPackages() ([]db.PluginPackage, error)
	GetCurrentPlugins() ([]db.PluginPackage, error)
	LoadPlugin(r *req.PluginPackageReq, browserGWs []browsergateway.ServiceInstance) error
}

var _ PluginService = &PluginServiceImpl{}

func NewPluginService() *PluginServiceImpl {
	return &PluginServiceImpl{ppd: dao.NewPluginPackageDao(), fd: dao.NewFileDao(), httpClient: https.Instance()}
}

type PluginServiceImpl struct {
	ppd        *dao.PluginPackageDao
	fd         *dao.FileDao
	httpClient https.HTTPDoer
}

func (p *PluginServiceImpl) LoadPlugin(r *req.PluginPackageReq, browserGWs []browsergateway.ServiceInstance) error {
	// 查询对应插件信息
	// 查询BrowserGW实例信息和在线状态
	// 异步任务调用BrowserGW加载插件并更新进度
	// 返回结果
	var pluginPackage = &db.PluginPackage{
		Field: r.GetKey(),
	}
	err := p.ppd.Get(pluginPackage)
	if err != nil {
		logger.Errorf("get pluginPackage %s failed, err is %v", pluginPackage.GetKey(), err)
		return err
	}
	logger.Infof("pluginActive %s is %v", pluginPackage.GetField(), pluginPackage)
	err = p.switchActivePlugin(pluginPackage)
	if err != nil {
		logger.Errorf("set pluginActive %s failed, err is %v", pluginPackage.GetKey(), err)
		return nil
	}

	var progressChan = make(chan db.PluginPackage, len(browserGWs))
	go p.loadPlugin(pluginPackage, browserGWs, progressChan)
	go p.recordLoadPluginProgress(progressChan)

	return nil
}

func (p *PluginServiceImpl) switchActivePlugin(pluginPackage *db.PluginPackage) error {
	// 不考虑任务并发执行
	pluginPackage.Status = db.Doing
	pluginPackage.Progress = 0
	pluginPackage.IfActive = true
	err := p.ppd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		if _, err := p.ppd.ExecWithOrm(ctx, txOrm, "update t_plugin_package set if_active=false where plugin_type=?",
			constants.ChromeExtendType); err != nil {
			return err
		}
		if err := p.ppd.UpdateWithOrm(ctx, txOrm, pluginPackage); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (p *PluginServiceImpl) GetCurrentPlugins() ([]db.PluginPackage, error) {
	var pluginActive = db.PluginPackage{
		Type: constants.ChromeExtendType,
	}
	var ppl []db.PluginPackage
	err := p.ppd.List(&ppl, *dao.NewQueryOption().Filter("Type", constants.ChromeExtendType).Filter("IfActive", true))
	if err != nil {
		logger.Errorf("get pluginActive %s failed, err is %v", pluginActive.Type, err)
		return nil, err
	}
	return ppl, nil
}

func (p *PluginServiceImpl) UploadPluginPackage(req *req.UploadPluginPackageReq) error {
	defer func(File multipart.File) {
		err := File.Close()
		if err != nil {
			logger.Warnf("file close error: %s", err)
		}
	}(req.File)

	pkg, err := p.readPackageMeta(req)
	if err != nil {
		return err
	}
	logger.Infof("upload package, meta is %+v", pkg)
	// 检查数据是否存在
	oldPP := &db.PluginPackage{
		Field: pkg.GetField(),
	}
	err = p.ppd.Get(oldPP)
	if err != nil && err != orm.ErrNoRows {
		logger.Errorf("check key %s exist failed, err is %v", pkg.GetKey(), err)
		return err
	}
	if err == nil {
		err = fmt.Errorf("key %s is exist, forbidden upload", pkg.GetKey())
		return err
	}
	// 上传软件包
	_, err = req.File.Seek(0, io.SeekStart)
	if err != nil {
		logger.Errorf("seek file to start error: %s", err)
		return err
	}
	content, err := io.ReadAll(req.File)
	if err != nil {
		logger.Errorf("read request file failed: %v", err)
		return err
	}
	return p.ppd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		pkg.Status = db.NotStart
		pkg.IfActive = false
		f := &db.File{
			Bucket:    pkg.PackageBucket,
			Name:      pkg.PackageName,
			Content:   content,
			Size:      req.Size,
			CreatedAt: time.Now().Format(time.DateTime),
		}
		if err := p.fd.InsertWithOrm(ctx, txOrm, f); err != nil {
			logger.Errorf("put file %s to storage failed, err is %v", req.Filename, err)
			return err
		}
		pkg.Field = pkg.GetField()
		pkg.CreatedAt = time.Now().Format(time.DateTime)
		if err := p.ppd.InsertWithOrm(ctx, txOrm, pkg); err != nil {
			logger.Errorf("Set data to redis failed, data: %v, err %v", pkg, err)
			return err
		}
		return nil
	})
}

func (p *PluginServiceImpl) readPackageMeta(request *req.UploadPluginPackageReq) (*db.PluginPackage, error) {
	var handlers = map[string]func(*req.UploadPluginPackageReq) (*db.PluginPackage, error){
		".zip": p.readMetaFromZip,
	}
	for key, handler := range handlers {
		if strings.HasSuffix(request.Filename, key) {
			pluginPackage, err := handler(request)
			if err != nil {
				return nil, err
			}
			return pluginPackage, nil
		}
	}
	err := fmt.Errorf("file %s is not a zip file", request.Filename)
	return nil, err
}

func (p *PluginServiceImpl) readMetaFromZip(req *req.UploadPluginPackageReq) (*db.PluginPackage, error) {
	zipReader, err := zip.NewReader(req.File, req.Size)
	if err != nil {
		logger.Errorf("open zip file error: %v", err)
		return nil, err
	}
	var pkg *db.PluginPackage
	for _, file := range zipReader.File {
		// 非描述文件跳过
		if file.Name != constants.ExtensionDescribeFile {
			continue
		}
		// 超大失败
		if file.FileInfo().Size() > constants.MaxFileSize {
			err = fmt.Errorf("open package.json error, because file size is %d", file.FileInfo().Size())
			logger.Errorf("%v", err)
			return nil, err
		}
		fileReader, err := file.Open()
		// 读取package.json
		err = func() error {
			defer fileReader.Close()
			if err != nil {
				logger.Errorf("open package.json error: %s", err)
				return err
			}
			fileContent, err := io.ReadAll(fileReader)
			if err != nil {
				logger.Errorf("read package.json error: %s", err)
				return err
			}
			pkg = new(db.PluginPackage)
			err = json.Unmarshal(fileContent, pkg)
			if err != nil {
				logger.Errorf("load package.json error: %s", err)
				return err
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
		pkg.PackageBucket = constants.PluginPackageBucket
		pkg.PackageName = req.Filename
		if pkg.Type != constants.ChromeExtendType || pkg.Name == "" || pkg.Version == "" {
			err = errors.New("load package.json error, content exception")
			logger.Errorf("type %s, name %s, version %s, err is %v",
				pkg.Type, pkg.Name, pkg.Version, err)
			return nil, err
		}
		return pkg, nil
	}
	err = errors.New("file package.json not found")
	logger.Errorf("%v", err)
	return nil, err
}

func (p *PluginServiceImpl) DeletePluginPackage(request *req.PluginPackageReq) error {
	var pkg = db.PluginPackage{
		Name:    request.Name,
		Type:    request.Type,
		Version: request.Version,
	}
	oldPP := &db.PluginPackage{
		Field: pkg.GetField(),
	}
	err := p.ppd.Get(oldPP)
	if err != nil && err == orm.ErrNoRows {
		logger.Infof("plugin:%v is already deleted", pkg.GetField())
		return nil
	}
	if err != nil {
		logger.Errorf("get plugin package %s failed, err is %v", pkg.GetKey(), err)
		return err
	}
	if oldPP.IfActive {
		logger.Errorf("plugin is in used, cannot delete it!")
		return errors.New("plugin is in used, cannot delete it")
	}

	return p.ppd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		if _, err := txOrm.Delete(oldPP); err != nil {
			logger.Errorf("delete plugin package %s failed, err is %v", pkg.GetKey(), err)
			return err
		}
		oldFile := &db.File{
			Bucket: oldPP.PackageBucket,
			Name:   oldPP.PackageName,
		}
		if _, err := txOrm.Delete(oldFile); err != nil {
			logger.Errorf("delete plugin package file %s failed, err is %v", pkg.GetKey(), err)
			return err
		}
		return nil
	})
}

func (p *PluginServiceImpl) GetPluginPackages() ([]db.PluginPackage, error) {
	var ppl []db.PluginPackage
	if err := p.ppd.List(&ppl); err != nil {
		logger.Errorf("list plugin package failed, err is %v ", err)
		return nil, err
	}
	return ppl, nil
}

func (p *PluginServiceImpl) recordLoadPluginProgress(progressChan <-chan db.PluginPackage) {
	for pluginActive := range progressChan {
		// 未考虑重试，未考虑服务重启中断，未记录失败节点
		logger.Infof("record progress, pluginActive %s is %v", pluginActive.GetField(), pluginActive)
		err := p.ppd.Update(&pluginActive)
		if err != nil {
			logger.Errorf("set redis data %s:%v failed, err is %v", pluginActive.GetKey(), pluginActive, err)
			continue
		}
	}
}

func (p *PluginServiceImpl) loadPlugin(
	pluginPackage *db.PluginPackage,
	browserGWs []browsergateway.ServiceInstance,
	progressChan chan db.PluginPackage) {
	var completeCount = 0
	for _, browserGW := range browserGWs {
		extensionLoadResponse, err := p.loadPluginToBrowserGW(browserGW, pluginPackage)
		if err != nil {
			continue
		}
		if extensionLoadResponse.Code == http.StatusOK {
			completeCount++
			pluginPackage.Progress = completeCount * percentageMultiplier / len(browserGWs)
			if pluginPackage.Progress == progressComplete {
				pluginPackage.Status = db.Complete
			}
			progressChan <- *pluginPackage
		} else {
			logger.Errorf("load plugin to %s failed, resp is %v", browserGW.BrowserInnerEndpoint, extensionLoadResponse)
		}
	}
	if pluginPackage.Status != db.Complete {
		pluginPackage.Status = db.Failed
	}
	progressChan <- *pluginPackage
	close(progressChan)
}

func (p *PluginServiceImpl) loadPluginToBrowserGW(browserGW browsergateway.ServiceInstance,
	pluginPackage *db.PluginPackage) (browsergateway.ExtensionLoadResponse, error) {
	response := https.NewRequest(p.httpClient).
		Method("POST").WithRetry(defaultRetryCount).
		URL(fmt.Sprintf("http://%s/browsergw/extension/load", browserGW.BrowserInnerEndpoint)).
		ParamFromInterface(
			browsergateway.ExtensionLoadRequest{
				BucketName:        pluginPackage.PackageBucket,
				ExtensionFilePath: pluginPackage.PackageName,
				Name:              pluginPackage.Name,
				Version:           pluginPackage.Version,
				Type:              pluginPackage.Type,
			}).
		Complete().
		Do()
	if !response.IsSuccessCode() || response.Error() != nil {
		logger.Errorf("load plugin to %s failed, status is %d, err is %v",
			browserGW.BrowserInnerEndpoint, response.StatusCode(), response.Error())
		return browsergateway.ExtensionLoadResponse{}, response.Error()
	}
	var extensionLoadResponse browsergateway.ExtensionLoadResponse
	err := response.ResponseToStruct(&extensionLoadResponse)
	if err != nil {
		logger.Errorf("ResponseToStruct failed, err is %v", err)
		return browsergateway.ExtensionLoadResponse{}, err
	}
	logger.Infof("BrowserGW %s load plugin success, return %v", browserGW.BrowserInnerEndpoint, extensionLoadResponse)
	return extensionLoadResponse, nil
}
