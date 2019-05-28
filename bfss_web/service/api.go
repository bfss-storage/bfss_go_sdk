package api

import (
	"bfss/bfss_sdk_api"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ReadBlkData struct {
	OId    string `json:"oid"`    //对象id
	Size   int32  `json:"size"`   //预读数据大小，-1全部
	Offset int32  `json:"offset"` //对象的数据偏移（16对齐）
}

// 错误响应
type errorResponse struct {
	OK        bool        `json:"ok"`         // 是否成功
	ErrorCode int         `json:"error_code"` // 错误代码
	Result    interface{} `json:"result"`     // 返回结果
}

// 默认错误响应
var defaultErrorResponse = []byte(`{"Result":false,"Data":"internal error"}`)
var parameterErrorResponse = []byte(`{"Result":false,"Data":"wrong parameter in request"}`)
var resultErrorResponse = `{"Result":"false","Data":"%v"}`

/*
函数名： ReadBlk 读取加密对象
参数：
oid 对象id，同上
size 预读数据大小，-1全部
offset 对象的数据偏移（16对齐）
返回值：
     参考READ_RESULT结构体说明
READ_RESULT ReadBlk(1:string oid, 2:i32 size, 3:i32 offset)
*/
func ReadBlk(w http.ResponseWriter, req *http.Request) {
	// 读取请求内容(url传参)
	req.ParseForm()
	OId, ok := req.Form["OId"]
	if !ok {
		log.Print("From ReadBlk Func,parameter OId nil:", ok)
		w.Write(parameterErrorResponse)
		return
	}
	size, ok := req.Form["size"]
	if !ok {
		log.Print("From ReadBlk Func,parameter size nil:", ok)
		w.Write(parameterErrorResponse)
		return
	}
	offset, ok := req.Form["offset"]
	if !ok {
		log.Print("From ReadBlk Func,parameter offset nil:", ok)
		w.Write(parameterErrorResponse)
		return
	}

	sizeInt, err := strconv.ParseInt(size[0], 10, 32)
	if err != nil {
		log.Print("From ReadBlk Func,size err:", err)
		w.Write(defaultErrorResponse)
		return
	}
	offsetInt, err := strconv.ParseInt(offset[0], 10, 32)
	if err != nil {
		log.Print("From ReadBlk Func,offset err:", err)
		w.Write(defaultErrorResponse)
		return
	}

	readBlk := &ReadBlkData{
		OId[0],
		int32(sizeInt),
		int32(offsetInt),
	}
	log.Print("readBlk : ", readBlk)

	/*// 读取请求内容(json格式)
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("From ReadBlk Func,ReadAll err:", err)
		w.Write(defaultErrorResponse)
		return
	}
	defer req.Body.close()

	readBlk := &ReadBlkData{}
	err = json.Unmarshal(data, readBlk)
	if err != nil {
		log.Print("From ReadBlk Func,Unmarshal err:", err)
		w.Write(defaultErrorResponse)
		return
	}
	log.Print("readBlk : ", readBlk)*/

	d, err := bfss_sdk_api.ReadBlk(readBlk.OId, readBlk.Size, readBlk.Offset)
	if err != nil {
		log.Print("From ReadBlk Func,ApiClient.ReadBlk err:", err)
		w.Write([]byte(fmt.Sprintf(resultErrorResponse, err.Error())))
		return
	}
	log.Print("readBlk data : ", d.String())
	w.Write(d.Data)
}
