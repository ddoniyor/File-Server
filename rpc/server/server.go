package main

import (
	"bufio"
	"file-server/pkg/rpc"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)


func main() {
	const addr = "0.0.0.0:9999"
	log.Print("server starting")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can't listen on %s: %v", addr, err)
	}
	defer listener.Close()
	log.Print("server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v", err)
			continue
		}
		go handleConn(conn)
	}
}


func handleConn(conn net.Conn) error {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("error while reading: %v", err)
		return nil
	}
	index := strings.IndexByte(line, ':')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		log.Printf("invalid line received %s", line)
		err := rpc.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return nil
		}
		return nil
	}


	cmd, options := line[:index], line[index+1:]
	log.Printf("command received: %s", cmd)
	log.Printf("options received: %s", options)

	switch cmd {

	case "upload":
		reader := bufio.NewReader(conn)

			line, err := rpc.ReadLine(reader)
			if err != nil {
				log.Printf("can't read: %v", err)
				return nil
			}
			log.Print(line)
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				if err != io.EOF {
					log.Printf("can't read data: %v", err)
				}
			}
			log.Print(len(bytes))
			err = ioutil.WriteFile("rpc/server/files/"+options, bytes, 0666)
			if err != nil {
				log.Printf("can't write file: %v", err)
			}
			fmt.Printf("File with name %s uploaded to server\n", options)


	case "download":
		options = strings.TrimSuffix(options, "\n")
		file, _:= os.Open("files/"+options)

		log.Print("file opened")
		err = rpc.WriteLine("result: ok", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return nil
		}

		_, err = io.Copy(writer, file)
		log.Print("file sent")

	case "list":
		options = strings.TrimSuffix(options, "\n")
		fileName:= rpc.GetListOfFiles("rpc/server/files")
		err := rpc.WriteLine(fileName, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}

	default:
		err := rpc.WriteLine("result: error", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return nil
		}
	}
	return nil
}
