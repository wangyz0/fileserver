package mq

type TransferData struct {
	FileHash      string
	CurLocation   string //在本地的临时地址
	DestLocation  string // 要转移到的地址
	DestStoreType string
}
