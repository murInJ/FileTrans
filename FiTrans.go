package FiTrans

import (
	"errors"
	"github.com/murInJ/FiTrans/TCP"
)

type FiTrans interface {
	StartServer()
	CloseServer() error
	Send(serverAddr string, filePath string, perSize int) error
	GetFileChannel() chan TCP.File
}

func NewFiTrans(tp string, port int, debug bool) (FiTrans, error) {
	var ft FiTrans
	if tp == "tcp" {
		tcp := newTcpTrans(port, debug)
		ft = tcp
	} else {
		return nil, errors.New("wrong param tp")
	}
	return ft, nil
}
