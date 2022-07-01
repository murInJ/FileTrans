package FiTrans

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/murInJ/FiTrans/TCP"
	"strconv"
)

type tcpTrans struct {
	port        int
	server      *TCP.TcpServer
	client      *TCP.TcpClient
	fileChannel chan TCP.File
	debug       bool
}

func newTcpTrans(port int, debug bool) *tcpTrans {
	return &tcpTrans{
		port:        port,
		debug:       debug,
		fileChannel: make(chan TCP.File, 50),
	}
}

func (t *tcpTrans) StartServer() {
	t.server = TCP.NewTcpServer(t.port, t.fileChannel, t.debug)
	fmt.Printf("%s server start at %s\n",
		color.New(color.FgHiYellow).Sprintf("FiTrans:"),
		color.New(color.FgYellow).Sprintf(TCP.GetIP_Local().String()+":"+strconv.Itoa(t.port)))
	t.server.Start()
}

func (t *tcpTrans) CloseServer() error {
	err := t.server.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%s server closed\n",
		color.New(color.FgHiYellow).Sprintf("FiTrans:"))
	return nil
}

func (t *tcpTrans) Send(serverAddr string, filePath string, perSize int) error {
	t.client = TCP.NewTcpClient(t.port, perSize, t.debug)
	err := t.client.Send(serverAddr, filePath)
	if err != nil {
		return err
	}
	return nil
}

func (t *tcpTrans) GetFileChannel() chan TCP.File {
	return t.fileChannel
}
