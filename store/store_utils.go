package store

import (
	"fmt"
	"mercury/message"
	"mercury/packet"
	"runtime"
	"strconv"
	"strings"
)

var ostype = runtime.GOOS

type Values struct {
	valueMap map[string]interface{}
}

func (v Values) Get(key string) interface{} {
	return v.valueMap[key]
}

func FileSeparator() string {
	if ostype == "windows" {
		return "\\\\"
	} else if ostype == "linux" || ostype == "unix" {
		return "/"
	}
	return ""
}

func FileName2Offset(flileName string) (int64, error) {
	offsetStrSlice := strings.Split(flileName, FileSeparator())
	if len(offsetStrSlice) == 0 {
		return int64(-1), fmt.Errorf("File name %s can not be splited", flileName)
	}
	return strconv.ParseInt(offsetStrSlice[len(offsetStrSlice)-1], 10, 64)
}

func Offset2FileName(offset int64) string {
	return fmt.Sprintf("%020s", strconv.FormatInt(offset, 10))
}

func TranferPacket2Msg(pkt *packet.WebsocketContext) *message.Msg {
	msg := message.NewMsg(pkt.WsPacket.Topic, 0)
	return msg
}
