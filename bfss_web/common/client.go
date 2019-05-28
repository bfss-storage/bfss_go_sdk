package common

import (
	"bfss/bfss_api"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"net"
	"sync"
)

var botOnce sync.Once
var ApiClient *bfss_api.BFSS_APIDClient

func CreateApiClient(host string, port string) {
	botOnce.Do(func() {
		tSocket, err := thrift.NewTSocket(net.JoinHostPort(host, port))
		if err != nil {
			log.Print("From ReadBlk Func,tSocket err:", err)
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, err := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

		ApiClient = bfss_api.NewBFSS_APIDClientFactory(transport, protocolFactory)

		if err := transport.Open(); err != nil {
			transport.Close()
			log.Print("From ReadBlk Func,Error opening:", host+":"+port)
		}
		//defer transport.close()
	})
}
