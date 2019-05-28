package bfss_sdk

import (
	"bfss/bfss_api"
	"bfss/bfss_regm"
	"bfss/bfss_sn"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

func dialRgm(addr string, port int32, connTimeout time.Duration) (*idleClient, error) {
	return dial(addr, port, connTimeout, func(c thrift.TClient) interface{} {
		return bfss_regm.NewBFSS_REGMDClient(c)
	})
}

func dialSn(addr string, port int32, connTimeout time.Duration) (*idleClient, error) {
	return dial(addr, port, connTimeout, func(c thrift.TClient) interface{} {
		return bfss_sn.NewBFSS_SNDClient(c)
	})
}

func dialApi(addr string, port int32, connTimeout time.Duration) (*idleClient, error) {
	return dial(addr, port, connTimeout, func(c thrift.TClient) interface{} {
		return bfss_api.NewBFSS_APIDClient(c)
	})
}

func dial(addr string, port int32, connTimeout time.Duration, newfun func(c thrift.TClient) interface{}) (*idleClient, error) {
	socket, err := thrift.NewTSocketTimeout(fmt.Sprintf("%s:%d", addr, port), connTimeout)
	if err != nil {
		return nil, toBfssException(err)
	}
	trans := thrift.NewTFramedTransport(socket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	client := newfun(thrift.NewTStandardClient(iprot, oprot))

	return &idleClient{
		c: client,
		s: socket,
	}, nil
}

func close(c *idleClient) error {
	return c.s.Close()
}

func InitAPI(addr string, port int32, maxConn, connTimeout, idleTimeout uint32) error {
	if poolApi == nil {
		poolApi = newPoolMap(newThriftPool(addr, port, maxConn, connTimeout, idleTimeout, dialApi, close))
	}
	return nil
}

func InitREGM(addr string, port int32, maxConn, connTimeout, idleTimeout uint32) error {
	if poolRgm == nil {
		poolRgm = newPoolMap(newThriftPool(addr, port, maxConn, connTimeout, idleTimeout, dialRgm, close))
	}
	return nil
}

var snMaxConn, snConnTimeout, snIdleTimeoutint uint32 = 10000, 100, 1800

func InitSN(maxConn, connTimeout, idleTimeout uint32) error {
	snMaxConn, snConnTimeout, snIdleTimeoutint = maxConn, connTimeout, idleTimeout
	return nil
}
