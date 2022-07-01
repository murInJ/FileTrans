package TCP

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"strings"
)

type fileInfo struct {
	Path          string `json:"path,omitempty"`
	ServerAddress string `json:"server_address,omitempty"`
	ClientAddress string `json:"client_address,omitempty"`
	Size          int    `json:"size,omitempty"`
	FileName      string `json:"file_name"`
	PerSize       int    `json:"per_size,omitempty"`
	TaskNum       int    `json:"task_num"`
	SegNum        int    `json:"seg_num,omitempty"`
	Mode          string `json:"mode"`
}

type File struct {
	Path          string `json:"path,omitempty"`
	ServerAddress string `json:"server_address,omitempty"`
	ClientAddress string `json:"client_address,omitempty"`
	Size          int    `json:"size,omitempty"`
	FileName      string `json:"file_name"`
	Data          []byte `json:"data"`
}

func GetIP_Local() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.ParseIP("127.0.0.1")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return net.ParseIP(localAddr[0:idx])
}

func errMsg(err error) {
	fmt.Printf("%s %s\n",
		color.New(color.FgHiRed).Sprintf("FiTransError:"),
		err.Error())
}
