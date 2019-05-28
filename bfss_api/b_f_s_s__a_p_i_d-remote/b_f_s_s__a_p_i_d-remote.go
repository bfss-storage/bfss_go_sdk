// Autogenerated by Thrift Compiler (0.12.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
        "context"
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        "github.com/apache/thrift/lib/go/thrift"
	"bfss/utils"
        "bfss/bfss_api"
)

var _ = utils.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  string GetVersion()")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT CreateObject(string oid, i32 size, i32 flag, string tag)")
  fmt.Fprintln(os.Stderr, "  OBJECT_INFO_EX_RESULT CreateObjectEx(string oid, i32 size, i32 flag, string tag)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT DeleteObject(string oid)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT Write(string oid, i32 offset, string data)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT ResizeObject(string oid, i32 newsize)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT ResetObjectId(string oid, string newoid, string newtag)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT CompleteObject(string oid)")
  fmt.Fprintln(os.Stderr, "  OBJECT_INFO_RESULT GetObjectInfo(string oid)")
  fmt.Fprintln(os.Stderr, "  OBJECT_INFO_EX_RESULT GetObjectInfoEx(string oid)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT ObjectLockHasHash(string hash, i32 size, string head)")
  fmt.Fprintln(os.Stderr, "  BFSS_RESULT CreateObjectLink(string oid, string hash, i32 size, string head, i32 flag, string tag)")
  fmt.Fprintln(os.Stderr, "  READ_RESULT Read(string oid, i32 size, i32 offset)")
  fmt.Fprintln(os.Stderr, "  READ_RESULT ReadBlk(string oid, i32 size, i32 offset)")
  fmt.Fprintln(os.Stderr, "  OBJECT_BLKS_RESULT GetObjectBlksInfo(string oid)")
  fmt.Fprintln(os.Stderr, "  BLK_KEY_RESULT GetObjectBlkKey(string oid, i32 offset)")
  fmt.Fprintln(os.Stderr, "  MESSAGE_RESULT ManageMessage(BFSS_CMD CmdId, i32 Param, string Data)")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
  var m map[string]string = h
  return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
  parts := strings.Split(value, ": ")
  if len(parts) != 2 {
    return fmt.Errorf("header should be of format 'Key: Value'")
  }
  h[parts[0]] = parts[1]
  return nil
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  headers := make(httpHeaders)
  var parsedUrl *url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
  flag.Parse()
  
  if len(urlString) > 0 {
    var err error
    parsedUrl, err = url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
    if len(headers) > 0 {
      httptrans := trans.(*thrift.THttpClient)
      for key, value := range headers {
        httptrans.SetHeader(key, value)
      }
    }
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  iprot := protocolFactory.GetProtocol(trans)
  oprot := protocolFactory.GetProtocol(trans)
  client := bfss_api.NewBFSS_APIDClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "GetVersion":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetVersion requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetVersion(context.Background()))
    fmt.Print("\n")
    break
  case "CreateObject":
    if flag.NArg() - 1 != 4 {
      fmt.Fprintln(os.Stderr, "CreateObject requires 4 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err38 := (strconv.Atoi(flag.Arg(2)))
    if err38 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err39 := (strconv.Atoi(flag.Arg(3)))
    if err39 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    argvalue3 := flag.Arg(4)
    value3 := argvalue3
    fmt.Print(client.CreateObject(context.Background(), value0, value1, value2, value3))
    fmt.Print("\n")
    break
  case "CreateObjectEx":
    if flag.NArg() - 1 != 4 {
      fmt.Fprintln(os.Stderr, "CreateObjectEx requires 4 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err42 := (strconv.Atoi(flag.Arg(2)))
    if err42 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err43 := (strconv.Atoi(flag.Arg(3)))
    if err43 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    argvalue3 := flag.Arg(4)
    value3 := argvalue3
    fmt.Print(client.CreateObjectEx(context.Background(), value0, value1, value2, value3))
    fmt.Print("\n")
    break
  case "DeleteObject":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "DeleteObject requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.DeleteObject(context.Background(), value0))
    fmt.Print("\n")
    break
  case "Write":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "Write requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err47 := (strconv.Atoi(flag.Arg(2)))
    if err47 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    argvalue2 := []byte(flag.Arg(3))
    value2 := argvalue2
    fmt.Print(client.Write(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "ResizeObject":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "ResizeObject requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err50 := (strconv.Atoi(flag.Arg(2)))
    if err50 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    fmt.Print(client.ResizeObject(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "ResetObjectId":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "ResetObjectId requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    argvalue1 := flag.Arg(2)
    value1 := argvalue1
    argvalue2 := flag.Arg(3)
    value2 := argvalue2
    fmt.Print(client.ResetObjectId(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "CompleteObject":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "CompleteObject requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.CompleteObject(context.Background(), value0))
    fmt.Print("\n")
    break
  case "GetObjectInfo":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetObjectInfo requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.GetObjectInfo(context.Background(), value0))
    fmt.Print("\n")
    break
  case "GetObjectInfoEx":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetObjectInfoEx requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.GetObjectInfoEx(context.Background(), value0))
    fmt.Print("\n")
    break
  case "ObjectLockHasHash":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "ObjectLockHasHash requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err58 := (strconv.Atoi(flag.Arg(2)))
    if err58 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    argvalue2 := []byte(flag.Arg(3))
    value2 := argvalue2
    fmt.Print(client.ObjectLockHasHash(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "CreateObjectLink":
    if flag.NArg() - 1 != 6 {
      fmt.Fprintln(os.Stderr, "CreateObjectLink requires 6 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    argvalue1 := flag.Arg(2)
    value1 := argvalue1
    tmp2, err62 := (strconv.Atoi(flag.Arg(3)))
    if err62 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    argvalue3 := []byte(flag.Arg(4))
    value3 := argvalue3
    tmp4, err64 := (strconv.Atoi(flag.Arg(5)))
    if err64 != nil {
      Usage()
      return
    }
    argvalue4 := int32(tmp4)
    value4 := argvalue4
    argvalue5 := flag.Arg(6)
    value5 := argvalue5
    fmt.Print(client.CreateObjectLink(context.Background(), value0, value1, value2, value3, value4, value5))
    fmt.Print("\n")
    break
  case "Read":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "Read requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err67 := (strconv.Atoi(flag.Arg(2)))
    if err67 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err68 := (strconv.Atoi(flag.Arg(3)))
    if err68 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.Read(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "ReadBlk":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "ReadBlk requires 3 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err70 := (strconv.Atoi(flag.Arg(2)))
    if err70 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    tmp2, err71 := (strconv.Atoi(flag.Arg(3)))
    if err71 != nil {
      Usage()
      return
    }
    argvalue2 := int32(tmp2)
    value2 := argvalue2
    fmt.Print(client.ReadBlk(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "GetObjectBlksInfo":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetObjectBlksInfo requires 1 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    fmt.Print(client.GetObjectBlksInfo(context.Background(), value0))
    fmt.Print("\n")
    break
  case "GetObjectBlkKey":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "GetObjectBlkKey requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    tmp1, err74 := (strconv.Atoi(flag.Arg(2)))
    if err74 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    fmt.Print(client.GetObjectBlkKey(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "ManageMessage":
    if flag.NArg() - 1 != 3 {
      fmt.Fprintln(os.Stderr, "ManageMessage requires 3 args")
      flag.Usage()
    }
    tmp0, err := (strconv.Atoi(flag.Arg(1)))
    if err != nil {
      Usage()
     return
    }
    argvalue0 := utils.BFSS_CMD(tmp0)
    value0 := argvalue0
    tmp1, err75 := (strconv.Atoi(flag.Arg(2)))
    if err75 != nil {
      Usage()
      return
    }
    argvalue1 := int32(tmp1)
    value1 := argvalue1
    argvalue2 := []byte(flag.Arg(3))
    value2 := argvalue2
    fmt.Print(client.ManageMessage(context.Background(), value0, value1, value2))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
