package server
import (
	"net"
	"log"
	"strconv"
	"bufio"
	"../../../protocol"
	"../service"
)


func StartListening(port int) nil {
	log.Printf("Start listening port %d ...\n", port)
	path := "127.0.0.1:" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", path)
    if err != nil {
        log.Fatal(err)
    }
	log.Printf("Accept connection on port %d.", port)

	for {
		conn, err := ln.Accept()
        if err != nil {
            log.Fatal(err)
        }
		go connectionHandler(conn)
	}

}

func connectionHandler(conn net.Conn) nil {
	defer conn.Close()
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Println(netData)
	result := ParseAndExecuteRequest(netData)
	conn.Write([]byte(result))
}

