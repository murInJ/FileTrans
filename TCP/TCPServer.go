package TCP

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net"
	"strconv"
	"sync"
)

type TcpServer struct {
	port          int
	listen        net.Listener
	debug         bool
	serverAddress string
	receiverMap   map[string]*receiver
	channelMap    map[string]chan int
	fileChannel   chan File
}

func NewTcpServer(port int, fileChannel chan File, debug bool) *TcpServer {
	address := GetIP_Local().String() + ":" + strconv.Itoa(port)
	return &TcpServer{
		port:          port,
		serverAddress: address,
		debug:         debug,
		fileChannel:   fileChannel,
		receiverMap:   make(map[string]*receiver),
		channelMap:    make(map[string]chan int),
	}
}

func (s *TcpServer) Start() {
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.port))
	if err != nil {
		errMsg(err)
		return
	}
	s.listen = listen

	for {
		conn, err := listen.Accept()
		if err != nil {
			errMsg(err)
			return
		}
		//创建协程
		go s.handleConn(conn)
	}
}

func (s *TcpServer) Close() error {
	err := s.listen.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *TcpServer) handleConn(conn net.Conn) {
	buf := make([]byte, 2048)
	n, _ := conn.Read(buf)
	var t interface{}
	json.Unmarshal(buf[:n], &t)
	m := t.(map[string]interface{})
	mode := m["mode"].(string)
	fileName := m["file_name"].(string)
	taskNum := int(m["task_num"].(float64))

	if mode == "info" {
		s.channelMap[fileName] = make(chan int, taskNum+5)
		s.receiverMap[fileName] = newReceiver(conn, m, s.channelMap[fileName], s.debug)
		conn.Write([]byte("ok"))

		for i := 1; i <= taskNum; i++ {
			segNum := <-s.channelMap[fileName]
			percent := strconv.Itoa(100 * i / taskNum)

			if s.debug {
				fmt.Printf("%s recv seg%d success %s%% \n",
					color.New(color.FgHiYellow).Sprintf("FiTrans:"),
					segNum,
					color.New(color.FgGreen).Sprintf(percent))
			}
		}

		var data []byte
		for i := 1; i <= s.receiverMap[fileName].TaskNum; i++ {
			data = append(data, s.receiverMap[fileName].SegMp[i]...)
		}
		f := File{
			Path:          s.receiverMap[fileName].Path,
			ServerAddress: s.receiverMap[fileName].ServerAddress,
			ClientAddress: s.receiverMap[fileName].ClientAddress,
			Size:          s.receiverMap[fileName].Size,
			FileName:      s.receiverMap[fileName].FileName,
			Data:          data,
		}
		s.fileChannel <- f
		if s.debug {
			fmt.Printf("%s recv file %s success \n",
				color.New(color.FgHiYellow).Sprintf("FiTrans:"),
				color.New(color.FgYellow).Sprintf(fileName))
		}

		conn.Write([]byte("ok"))
		conn.Close()

	} else {
		conn.Write([]byte("ok"))
		segNum := int(m["seg_num"].(float64))
		s.receiverMap[fileName].recvFile(conn, segNum)
		conn.Write([]byte("ok"))
		conn.Close()
	}
}

type receiver struct {
	conn          net.Conn
	Path          string
	ServerAddress string
	ClientAddress string
	Size          int
	FileName      string
	PerSize       int
	TaskNum       int
	Count         int
	SegMp         map[int][]byte
	debug         bool
	writeLock     sync.Mutex
	msgChannel    chan int
}

func newReceiver(conn net.Conn, mp map[string]interface{}, channel chan int, debug bool) *receiver {
	return &receiver{
		conn:          conn,
		Path:          mp["path"].(string),
		ServerAddress: mp["server_address"].(string),
		ClientAddress: mp["client_address"].(string),
		Size:          int(mp["size"].(float64)),
		FileName:      mp["file_name"].(string),
		PerSize:       int(mp["per_size"].(float64)),
		TaskNum:       int(mp["task_num"].(float64)),
		Count:         0,
		SegMp:         make(map[int][]byte),
		debug:         debug,
		msgChannel:    channel,
	}
}

func (r *receiver) recvFile(conn net.Conn, segNum int) {
	buf := make([]byte, r.PerSize)
	n, err := conn.Read(buf)
	if err != nil {
		errMsg(err)
		return
	}

	r.writeLock.Lock()
	r.SegMp[segNum] = buf[:n]
	r.writeLock.Unlock()

	r.msgChannel <- segNum
}
