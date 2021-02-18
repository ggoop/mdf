package services

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ggoop/mdf/bootstrap/errors"
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

//上传
func (s *ossSvImpl) UploadObject(ossConfig *model.Oss, fileItem model.OssObject, file multipart.File, header *multipart.FileHeader) (*model.OssObject, error) {
	if ossConfig != nil {
		fileItem.OssID = ossConfig.ID
		fileItem.OssType = ossConfig.Type
		fileItem.OssBucket = ossConfig.Bucket
	} else {
		fileItem.OssType = "local"
		fileItem.OssBucket = "storage/uploads"
	}
	if fileItem.OssBucket == "" {
		fileItem.OssBucket = "storage/uploads"
	}
	fileItem.ID = utils.GUID()
	fileItem.Code = utils.GUID()
	fileItem.Type = "obj"
	fileItem.OriginalName = header.Filename
	fileItem.Name = header.Filename
	fileItem.Size = header.Size
	fileItem.Ext = path.Ext(header.Filename)
	if err := s.uploadObjectByLocal(ossConfig, &fileItem, file); err != nil {
		return nil, err
	}
	return s.SaveObject(&fileItem)
}

func (s *ossSvImpl) uploadObjectByLocal(item *model.Oss, fileItem *model.OssObject, file multipart.File) error {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fileItem.MimeType = http.DetectContentType(fileBytes)
	if s := strings.Split(fileItem.MimeType, ";"); len(s) > 1 {
		fileItem.MimeType = s[0]
	}
	fileItem.Path = fmt.Sprintf("%s%s%s", s.GetUploadObjectStoreDir(fileItem), fileItem.Code, fileItem.Ext)
	fileItem.Url = fmt.Sprintf("/api/oss/%s/show", fileItem.ID)
	outFilePathFull := utils.JoinCurrentPath(path.Join(fileItem.OssBucket, fileItem.Path))
	if err := utils.CreatePath(filepath.Dir(outFilePathFull)); err != nil {
		return err
	}
	outFile, err := os.OpenFile(outFilePathFull, os.O_WRONLY|os.O_CREATE, os.FileMode(0777))
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	defer outFile.Close()
	outFile.Write(fileBytes)
	return nil
}

// 获取文件对象存储目录,如 temp/，kdirl/ ，kdnfd/dkfdsf/,以/ 结尾
func (s *ossSvImpl) GetUploadObjectStoreDir(fileItem *model.OssObject) string {
	fileKey := ""
	tag := false
	if fileItem.DirectoryID != "" { //指定文件夹时，按指定的文件夹存放
		if dir, _ := s.GetObjectBy(fileItem.DirectoryID); dir != nil {
			fileKey = filepath.Join(dir.Code, fileKey)
			tag = true
		}
	}
	if !tag && fileItem.Folder != "" { //按指定的目录存在
		fileKey = filepath.Join(fileItem.Folder, fileKey)
		tag = true
	}
	if !tag { //存在到临时文件夹中
		fileKey = filepath.Join(utils.TimeNow().Format("200601"), fileKey)
		tag = true
	}
	return fileKey + "/"
}
func (s *ossSvImpl) DeleteObject(entID, id string) error {
	item, err := s.GetObjectBy(id)
	if err != nil {
		return err
	}
	count := 0
	if item.Type == "dir" {
		if db.Default().Model(model.OssObject{}).Where("directory_id=?", item.ID).Count(&count); count > 0 {
			return fmt.Errorf("目录 %v 下存在 %v 个文件，不能被删除!", item.Name, count)
		}
	}
	//引用校验

	if count > 0 {
		return errors.IsQuoted(item.ID)
	}
	if err := db.Default().Delete(&item).Error; err != nil {
		return err
	}
	if item.OssType == "local" {
		os.Remove(utils.JoinCurrentPath(path.Join(item.OssBucket, item.Path)))
	}
	return nil
}
func (s *ossSvImpl) ObjectMove(objectIDs []string, directoryID string) error {
	if obj, err := s.GetObjectBy(directoryID); err != nil {
		return err
	} else {
		db.Default().Model(&model.OssObject{}).Where("id in (?)", objectIDs).Update("DirectoryID", obj.ID)
	}
	return nil
}
func (s *ossSvImpl) ObjectUpdates(objectID string, updates map[string]interface{}) error {
	if obj, err := s.GetObjectBy(objectID); err != nil {
		return err
	} else if len(updates) > 0 {
		db.Default().Model(&model.OssObject{}).Where("id =?", obj.ID).Updates(updates)
	}
	return nil
}
