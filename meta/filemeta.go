package meta

import "fileserver/db"

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

//将文件信息传给db操作数据库插入
func UpdataFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//查询文件元信息
func GetFileMeta(filehash string) FileMeta {

}
