package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	. "wanbu_stat_zzmm/src/logs"
	"wanbu_stat_zzmm/src/protocol"
)

type Walkday struct {
	Walkdate  int64  `json:"walkdate"`
	Walkhour  string `json:"walkhour"`
	Walktotal int    `json:"walktotal"`
	Recipe    string `json:"recipe"`
}

type Walkdata struct {
	Userid    int       `json:"userid"`
	Timestamp int64     `json:"timestamp"`
	Walkdays  []Walkday `json:"walkdays"`
}

func send(conn net.Conn) {

	var total int

	n, err := conn.Write(protocol.Enpack(&protocol.Message{"javaserver@xxxooo", 0}))
	//n, err := conn.Write([]byte{114, 101, 103, 105, 115, 116, 101, 114, 0, 0, 0, 8, 106, 107, 39, 115, 32, 109, 97, 99})

	if err != nil {
		total += n
		fmt.Printf("write %d bytes, error:%s\n", n, err)
		os.Exit(1)
	}
	total += n
	fmt.Printf("write regist %d bytes this time, %d bytes in total\n", n, total)

	var total0 int

	n0, err0 := conn.Write(protocol.Enpack(&protocol.Message{"heartbeat", 1}))

	if err0 != nil {
		total0 += n0
		fmt.Printf("write %d bytes, error:%s\n", n0, err0)
		os.Exit(1)
	}
	total0 += n0
	fmt.Printf("write heartbeat %d bytes this time, %d bytes in total\n", n0, total0)

	//os.Exit(1)

	//defer conn.Close()
}

func HandleRead(conn net.Conn) {

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan *protocol.Message)
	//fmt.Printf("%d connection connected into server\n", index)
	go reader(conn, readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {

			if err != io.EOF {

				Logger.Debug(conn.RemoteAddr().String(), " connection error: ", err)
				return
			}
		}

		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
	defer conn.Close()
}

func reader(conn net.Conn, readerChannel chan *protocol.Message) {
	for {
		select {

		case data := <-readerChannel:

			switch data.MsgType {
			/*
				新来的注册client，需要Server先发送心跳包，开始双方之间的aliveCheck，同时启动SetDeadline，
				如超时未收到消息，则关闭链接
			*/
			case 0:
				fmt.Println("zero")
			//收到心跳包，重启计时；否则，短连接处理到时，会销毁conn
			case 1:
				//Logger.Debug(conn.RemoteAddr().String(), "receive data string:", data.MsgContent)
				//Logger.Debug("record time is ", time.Now())
				n, err := conn.Write(protocol.Enpack(&protocol.Message{"heartbeat", 1}))

				if err != nil {
					fmt.Printf("write %d bytes, error:%s\n", n, err)
					os.Exit(1)
				}
				time.Sleep(2 * time.Second)
				//conn.Write([]byte("abcdefghijk"))
				//conn.SetDeadline(time.Now().Add(time.Duration(protocol.HBTimeOut) * time.Second))

			case 2:
				fmt.Println("two")
			default:
				fmt.Println("weird happens")
			}

		}
	}
}

func main() {

	server := "localhost:6081"
	//server := "192.168.60.33:65432"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	go HandleRead(conn)
	send(conn)

	var s Walkdata

	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1457193600,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1457280000,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})

	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1457366400,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	/*s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1453132800,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1455897600,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1455984000,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1456070400,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})*/

	s.Timestamp = 1455724804

	for i := 0; i < 100000; i++ {

		total := 0

		for i := 402080; i < 402083; i++ {

			s.Userid = i

			b, err := json.Marshal(s)
			if err != nil {
				fmt.Println("json err:", err)
			}

			_, err1 := conn.Write(protocol.Enpack(&protocol.Message{string(b), 2}))
			if err1 != nil {
				fmt.Println("in for run ", err1)
				os.Exit(1)
			}
			time.Sleep(time.Duration(1000*2) * time.Millisecond)

			total += 1
			fmt.Println("total send msg is ", total)

		}

		//time.Sleep(time.Duration(60) * time.Second)
	}
}
