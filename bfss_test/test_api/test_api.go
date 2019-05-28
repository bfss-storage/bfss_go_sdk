//
package main

import (
	"bfss/bfss_sdk"
	"bfss/bfss_sdk_api"
	"bfss/utils"
	"fmt"
	"github.com/oklog/ulid"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const ChipSize = 32*1024


func getOId() (string, error) {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	entropy := rand.New(source)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		fmt.Printf("ulid.New  failed.\n")
		return "", err
	}
	//fmt.Printf("New oid=%s", id.String())
	return id.String(), err
}


func testAPI(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Open file failed.\n")
		return
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Stat file failed.\n")
		return
	}

	oid, err := getOId()
	fmt.Printf("Process %s  with oid=%s\n", fileName, oid)
	var size = int32(fi.Size())
	var tag = "SDK-API"
	var flag int32 = 1
	/// --- obj<-CreateObject
	result, err := bfss_sdk_api.CreateObject(oid, size, flag, tag)
	if err != nil {
		fmt.Printf("CreateObject %s   for %s failed.\n", oid, fileName)
		return
	}
	if result != utils.BFSS_RESULT_BFSS_SUCCESS {
		fmt.Printf("CreateObject %s failed. %s\n", oid, result)
		return
	}

	var readSize, totalSize int32 = 0, size
	for readSize < totalSize {
		var chipSize = int(totalSize - readSize)
		if chipSize > ChipSize {
			chipSize = ChipSize
		}

		if readSize == 0 {
			num := rand.Int31n(100)	/*随机等待100s，模拟延时*/
			fmt.Printf("Sleeping %d seconds......\n", num)
			time.Sleep(  time.Duration(num) * time.Second)
		}

		data := make([]byte, chipSize)
		_, err := f.Read(data)
		if err != nil {
			fmt.Printf("Read %s   from %s failed.\n", oid, fileName)
			return
		}

		/// --- bfss_sdk_api.Write
		result, err = bfss_sdk_api.Write(oid, readSize, data)
		if result != utils.BFSS_RESULT_BFSS_SUCCESS {
			fmt.Printf("Write %s failed. %s\n", oid, result)
			return
		}
		readSize += int32(chipSize)
	}
	/// --- bfss_sdk_api.CompleteObject
	bfss_sdk_api.CompleteObject(oid)
	f.Close()
	fmt.Printf("============================ CompleteObject %s Writed.\n", oid)

	fileNew := strings.Replace(fileName, "testdata", "agenerated", 1)
	f, err = os.Create(fileNew)
	if err != nil {
		fmt.Printf("Create %s failed.\n", fileNew)
		return
	}

	var writedSize int32 = 0
	for writedSize < totalSize {
		var chipSize = int32(totalSize - writedSize)
		if chipSize > ChipSize {
			chipSize = ChipSize
		}

		if writedSize == 0 {
			num := rand.Int31n(100)	/*随机等待100s，模拟延时*/
			fmt.Printf("Sleeping %d seconds......\n", num)
			time.Sleep(  time.Duration(num) * time.Second)
		}

		/// --- bfss_sdk_api.Read
		result, err := bfss_sdk_api.Read(oid, chipSize, writedSize)
		if err != nil {
			fmt.Printf("Read %s   to %s failed.\n", oid, fileNew)
			return
		}
		if result.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
			fmt.Printf("Read %s  failed %s.\n", oid, result.Result_)
			return
		}

		_, err = f.Write(result.Data)
		if err != nil {
			fmt.Printf("file Write  %s failed.\n", fileNew)
			return
		}
		writedSize += int32(chipSize)
	}
	f.Close()
	fmt.Printf("Done  file-copy-via-bfss  from %s  to %s\n", fileName, fileNew)
}


func TestAPI() {
	var addr= "10.0.1.185"
	var port int32 = 9092
	var maxConn uint32 = 1024
	var connTimeout uint32 = 20
	var idleTimeout uint32 = 60
	err := bfss_sdk.InitAPI(addr, port, maxConn, connTimeout, idleTimeout)
	if err != nil {
		fmt.Printf("InitAPI failed.\n")
		return
	}

	addr = "10.0.1.185"
	port = 9090
	err = bfss_sdk.InitREGM(addr, port, maxConn, connTimeout, idleTimeout)
	if err != nil {
		fmt.Printf("InitREGM failed.\n")
		return
	}

	os.RemoveAll("../agenerated")
	os.Mkdir("../agenerated", 0755)
	files, err := filepath.Glob("../testdata/*.*")
	if err != nil {
		fmt.Printf("ReadDir failed.\n")
		return
	}

	fmt.Printf("start upload-download-via-bfss  files=%d\n", len(files))
	rand.Seed(rand.Int63())

	chFiles := make(chan string, len(files))

	var wg sync.WaitGroup
	for _, f := range files {
		chFiles <- f
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			testAPI(f)
		}(<-chFiles)
	}
	wg.Wait()

	fmt.Printf("-All- Done.\n")
	close(chFiles)
}

func main() {
	TestAPI()
}




