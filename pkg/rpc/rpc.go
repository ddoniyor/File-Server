package rpc

import (
	"bufio"
	"io/ioutil"
	"log"
)


func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(line string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return
}

func GetListOfFiles(line string) (list string) {
	files, err := ioutil.ReadDir(line)
	if err != nil {
		log.Printf("Can't read directory: %v", err)
	}
	for _, file := range files {
		if list == "" {
			list = list + file.Name()
		} else {
			list = list + " " + file.Name()
		}
	}
	list = list + "\n"
	return list
}