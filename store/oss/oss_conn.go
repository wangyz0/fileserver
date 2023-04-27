package oss

import (
	"fileserver/config"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

// 创建oss链接
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(config.OSSEndpoint, config.OSSAccesskey, config.OSSAccessKeySecret)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	return ossCli
}

// 获取bucket对象
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSSBucket)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		return bucket
	}
	return nil
}

// 从阿里云下载  获取临时下载url
func DownloadURL(objName string) (signedUrl string) {
	signedUrl, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return signedUrl
}
