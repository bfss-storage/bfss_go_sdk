//
package main

import (
	"bfss/bfss_sdk"
	"fmt"
	"github.com/oklog/ulid"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const ChipSize = 32 * 1024

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

func testStream(fileName string) {

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
	var tag = "SDK-STREAM"
	var flag int32 = 1
	/// --- obj<-CreateObject
	obj, err := bfss_sdk.CreateObject(oid, size, flag, tag)
	if err != nil {
		fmt.Printf("CreateObject %s   for %s failed.\n", oid, fileName)
		return
	}

	var readSize, totalSize int32 = 0, size
	for readSize < totalSize {
		var chipSize = int(totalSize - readSize)
		if chipSize > ChipSize {
			chipSize = ChipSize
		}

		if readSize == 0 {
			num := rand.Int31n(100) /*随机等待100s，模拟延时*/
			fmt.Printf("Sleeping %d seconds......\n", num)
			time.Sleep(time.Duration(num) * time.Second)
		}

		data := make([]byte, chipSize)
		_, err := f.Read(data)
		if err != nil {
			fmt.Printf("Read %s   from %s failed.\n", oid, fileName)
			return
		}

		/// --- obj.Write
		w, err := obj.Write(data)
		if w != chipSize {
			fmt.Printf("Write %s failed. writed %d\n", oid, w)
			return
		}
		readSize += int32(chipSize)
	}
	/// --- obj.close
	obj.Close()
	f.Close()
	fmt.Printf("============================ CompleteObject %s Writed.\n", oid)

	/// --- 读，需要先打开 类似文件操作
	obj, err = bfss_sdk.OpenObject(oid)
	if err != nil {
		fmt.Printf("OpenObject %s failed. \n", oid)
		return
	}

	fileNew := strings.Replace(fileName, "testdata", "sgenerated", 1)
	f, err = os.Create(fileNew)
	if err != nil {
		fmt.Printf("Create %s failed.\n", fileNew)
		return
	}

	var writedSize int32 = 0
	for writedSize < totalSize {
		var chipSize = int(totalSize - writedSize)
		if chipSize > ChipSize {
			chipSize = ChipSize
		}

		if writedSize == 0 {
			num := rand.Int31n(100) /*随机等待100s，模拟延时*/
			fmt.Printf("Sleeping %d seconds......\n", num)
			time.Sleep(time.Duration(num) * time.Second)
		}

		/// --- obj.Read
		rData := make([]byte, chipSize)
		rData, err = obj.Read(chipSize)
		if err != nil {
			fmt.Printf("Read %s   to %s failed.\n", oid, fileNew)
			return
		}

		_, err = f.Write(rData)
		if err != nil {
			fmt.Printf("file Write  %s failed.\n", fileNew)
			return
		}
		writedSize += int32(chipSize)
	}
	/// --- obj.close 读完需要关闭 类似文件操作
	obj.Close()
	f.Close()
	fmt.Printf("Done  file-copy-via-bfss  from %s  to %s\n", fileName, fileNew)
}

func TestStream() {
	var addr = "10.10.100.52"
	var port int32 = 9092
	var maxConn uint32 = 1024
	var connTimeout uint32 = 20
	var idleTimeout uint32 = 60
	err := bfss_sdk.InitAPI(addr, port, maxConn, connTimeout, idleTimeout)
	if err != nil {
		fmt.Printf("InitAPI failed.\n")
		return
	}

	addr = "10.10.100.52"
	port = 9090
	err = bfss_sdk.InitREGM(addr, port, maxConn, connTimeout, idleTimeout)
	if err != nil {
		fmt.Printf("InitREGM failed.\n")
		return
	}

	os.RemoveAll("../sgenerated")
	os.Mkdir("../sgenerated", 0755)
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
			testStream(f)
		}(<-chFiles)
	}
	wg.Wait()

	fmt.Printf("-All- Done.\n")
	close(chFiles)
}

func main() {
	TestStream()
}
