package main

import (
	"bufio"
	"file-server/pkg/rpc"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	)
var download = flag.String("download", "default", "download")
var upload = flag.String("upload", "default", "upload")
var list = flag.Bool("list", false, "list")

func main() {

	flag.Parse()
	var cmd, fileName string
	if *download != "default" {
		fileName = *download
		cmd = "download"
	} else if *upload != "default" {
		fileName = *upload
		cmd = "upload"
	} else if *list != false {
		cmd = "list"
	} else {
		return
	}

	operations(cmd, fileName)

}

func operations(cmd ,fileName string)  {
	addr := "localhost:9999"
	log.Print("client connecting")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("can't connect to %s: %v", addr, err)
	}
	defer conn.Close()
	log.Print("client connected")

	writer := bufio.NewWriter(conn)
	line :=cmd +":"+fileName
	log.Print("command sending")
	err = rpc.WriteLine(line, writer)
	if err != nil {
		log.Fatalf("can't send command %s to server: %v", line, err)
	}


	switch cmd {
	case "download":

		log.Print("command sent")
		downloadFromServer(conn,fileName)
		log.Print("downloaded from server successfully")

	case "upload":

		log.Print("command sent")
		uploadToServer(conn,fileName)
		log.Print("uploaded to server successfully ")

	case "list":

		log.Print("command sent")
		listFile(conn)
		log.Print("got list of files")

	default:
		fmt.Printf("Entered the wrong command: %s\n", cmd)
	}
	
}

func uploadToServer(conn net.Conn, fileName string) {

	options := strings.TrimSuffix(fileName, "\n")
	file, _:= os.Open("files/"+options)
	writer := bufio.NewWriter(conn)
	log.Print("file opened")
	err := rpc.WriteLine("result: ok", writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return
	}
	_, err = io.Copy(writer, file)
	log.Print("file sent")

}

func downloadFromServer(conn net.Conn, fileName string) {

	reader := bufio.NewReader(conn)
	for {
		line, err := rpc.ReadLine(reader)
		if err != nil {
			log.Printf("can't read: %v", err)
			return
		}
		log.Print(line)
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
			}
		}
		log.Print(len(bytes))
		err = ioutil.WriteFile("rpc/client/files/"+fileName, bytes, 0666)
		if err != nil {
			log.Printf("can't write file: %v", err)
		}
		fmt.Printf("File with name %s downloaded\n", fileName)
	}
}


func listFile(conn net.Conn) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	fmt.Println("List of files")
	var list string
	for i := 0; i < len(line); i++{
		if string(line[i]) == " " || string(line[i]) == "\n"{
			fmt.Println(list)
			list = ""
		} else {
			list = list + string(line[i])
		}
	}
	_, err = ioutil.ReadAll(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
}