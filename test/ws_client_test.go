package test

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/network"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestWSClient(t *testing.T) {
	// 服务器地址
	url := "ws://0.0.0.0:8999/ws"

	dp := network.NewDataPack(1048576)

	// 建立连接
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	closeChan := make(chan interface{})

	// 定期发送 ping 消息，保持连接活跃
	ticker := time.NewTicker(5 * time.Second)

	var idCnt = 1

	go func() {
		for {
			select {
			case <-closeChan:
				return
			case <-ticker.C:
				var content string = "123" + strconv.Itoa(idCnt)
				msg := network.NewMsgPackage(1, []byte(content))
				idCnt++
				bytes, err := dp.Pack(msg)
				if err != nil {
					log.Println("send msg pack err:", err)
					return
				}

				err = c.WriteMessage(websocket.BinaryMessage, bytes)
				if err != nil {
					log.Println("send msg write err:", err)
					return
				}

				log.Println("send")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-closeChan:
				return
			default:
				_, reader, err := c.NextReader()
				if err != nil {
					elog.Error("[Conn] read msg get reader err:", err)
					return
				}

				// 读取客户端的msg head
				headData := make([]byte, dp.GetHeadLen())
				if _, err := io.ReadFull(reader, headData); err != nil {
					elog.Error("[Conn] read msg head err:", err)
					return
				}

				// 拆包，得到 MsgId 和 DataLen
				msg, err := dp.Unpack(headData)
				if err != nil {
					elog.Error("[Conn] unpack msg err:", err)
					return
				}

				// 根据 DataLen 读取数据
				var data []byte
				if msg.GetDataLen() > 0 {
					data = make([]byte, msg.GetDataLen())
					if _, err := io.ReadFull(reader, data); err != nil {
						elog.Error("[Conn] read msg data err:", err)
						return
					}
				}
				msg.SetData(data)

				log.Println("recv ", string(data))
			}
		}
	}()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
