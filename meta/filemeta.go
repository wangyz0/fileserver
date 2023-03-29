package meta

import (
	"fileserver/db"
	"fmt"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// 将文件信息传给db操作数据库插入
func UpdataFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 查询文件元信息
func GetFileMeta(filehash string) FileMeta {

	DBfilemeta, err := db.GetFileMeta(filehash)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	FileMeta := FileMeta{
		FileSha1: DBfilemeta.FileSha1,
		FileName: DBfilemeta.FileName,
		FileSize: DBfilemeta.FileSize,
		Location: DBfilemeta.Location,
		UploadAt: DBfilemeta.UploadAt,
	}
	return FileMeta

}
