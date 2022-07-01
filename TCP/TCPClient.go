package TCP

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net"
	"os"
	"strconv"
)

type TcpClient struct {
	debug         bool
	port          int
	clientAddress string
	perSize       int
}

func NewTcpClient(port int, perSize int, debug bool) *TcpClient {
	return &TcpClient{
		port:          port,
		clientAddress: GetIP_Local().String() + ":" + strconv.Itoa(port),
		debug:         debug,
		perSize:       perSize,
	}
}

func (c *TcpClient) Send(serverAddress string, filePath string) error {
	var (
		err error
	)

	//segment file first
	taskNum, fileSlice, err := segmentation(filePath, c.perSize)
	if err != nil {
		return err
	}

	//get connect
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}

	//get file infomation
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	fileSize := info.Size()
	fileName := info.Name()

	infoStruct := fileInfo{
		Path:          filePath,
		ServerAddress: serverAddress,
		ClientAddress: c.clientAddress,
		Size:          int(fileSize),
		FileName:      fileName,
		PerSize:       c.perSize,
		TaskNum:       taskNum,
		Mode:          "info",
	}
	serialize, err := json.Marshal(infoStruct)
	if err != nil {
		return err
	}

	//send infoMsg
	conn.Write(serialize)
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	revData := string(buf[:n])
	if revData != "ok" {
		return errors.New("send fileInfo fail")
	}
	if c.debug {
		fmt.Printf("%s send fileInfo success\n",
			color.New(color.FgHiYellow).Sprintf("FiTrans:"))
	}

	//create sending task
	count := 1
	feedbackChannel := make(chan int, taskNum+5)
	for _, data := range fileSlice {
		go sendFile(serverAddress, fileName, count, data, feedbackChannel, taskNum)
		count++
	}

	//waiting
	for i := 1; i <= count-1; i++ {
		percent := strconv.Itoa(100 * i / (count - 1))
		cnt := <-feedbackChannel
		if c.debug {
			fmt.Printf("%s send fileSeg %d success %s %% \n",
				color.New(color.FgHiYellow).Sprintf("FiTrans:"),
				cnt,
				color.New(color.FgGreen).Sprintf(percent),
			)
		}
	}
	if c.debug {
		fmt.Printf("%s send file %s success\n",
			color.New(color.FgHiYellow).Sprintf("FiTrans:"),
			color.New(color.FgYellow).Sprintf(fileName),
		)
	}

	buf = make([]byte, 2048)
	n, err = conn.Read(buf)
	if err != nil {
		return err
	}
	revData = string(buf[:n])
	if revData != "ok" {
		conn.Close()
	}
	conn.Close()
	return nil
}

func segmentation(filePath string, perSize int) (int, [][]byte, error) {
	var (
		fileSlice [][]byte
		segLen    int
	)
	//read file
	f, err := os.Open(filePath)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()

	for {
		buf := make([]byte, perSize)
		n, err := f.Read(buf)
		if err != nil && io.EOF == err {
			return segLen, fileSlice, nil
		}
		fileSlice = append(fileSlice, buf[:n])
		segLen++
	}
}

func sendFile(serverAddress string, fileName string, count int, data []byte, feedBackChannel chan int, taskNum int) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		errMsg(err)
		return
	}
	defer conn.Close()

	info := fileInfo{
		FileName: fileName,
		SegNum:   count,
		Mode:     "seg",
		TaskNum:  taskNum,
	}
	serialize, _ := json.Marshal(info)
	conn.Write(serialize)

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		errMsg(err)
		return
	}
	revData := string(buf[:n])
	if revData != "ok" {
		errMsg(errors.New("send fileInfo fail"))
		return
	}

	conn.Write(data)
	buf = make([]byte, 2048)
	n, err = conn.Read(buf)
	if err != nil {
		errMsg(err)
		return
	}
	revData = string(buf[:n])
	if revData != "ok" {
		errMsg(errors.New("send fileSeg fail"))
		return
	}

	feedBackChannel <- count
}
