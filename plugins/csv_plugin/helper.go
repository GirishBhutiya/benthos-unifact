package csv_plugin

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jlaffaye/ftp"
)

var (
	ST_FTP          = 1
	ST_SHAREDFOLDER = 0
)

func getLastLineFromAbsPath(filepath string) (map[string]string, int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, 0, err
	}
	defer file.Close()

	return getLastLineFromFile(file)
}
func getLastLineFromFile(ir io.Reader) (map[string]string, int, error) {
	csvReader := csv.NewReader(ir)
	all, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err)
	}
	log.Println(len(all))
	if len(all) > 0 {
		headers := all[0]
		lastLine := all[len(all)-1]
		if len(headers) != len(lastLine) {
			return nil, 0, errors.New("headers and records are not the same length")
		}
		recordMap := make(map[string]string)
		for i, hdr := range headers {
			recordMap[hdr] = lastLine[i]
		}
		if f, ok := ir.(*os.File); ok {
			f.Close()
		}

		return recordMap, len(all), nil
	}
	return nil, 0, errors.New("no records found")
}
func getLastLineFromFTPFile(fr *ftp.Response) (map[string]string, int, error) {

	return getLastLineFromFile(fr)
}

func ParseNodeDef(subscription string) nodeDef {

	var nodeMap map[string][]map[string]string
	var node nodeDef
	err := json.Unmarshal([]byte(subscription), &nodeMap)
	if err != nil {
		log.Println(err)
		return node
	}
	for key, values := range nodeMap {
		for _, obj := range values {

			node.ID, _ = strconv.Atoi(key)
			node.Group = obj["group"]
			node.DB = obj["db"]
			node.Historian = obj["historian"]
			node.SqlSp = obj["sqlSp"]
			node.Name = obj["name"]

		}
	}

	return node

}
func convertSourceType(sourceType string) int {
	if sourceType == "ftp" {
		return ST_FTP
	} else {
		return ST_SHAREDFOLDER
	}

}
