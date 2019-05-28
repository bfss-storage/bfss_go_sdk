package bfss_sdk

import (
	"bfss/bfss_api"
	"bfss/bfss_regm"
	"bfss/bfss_sn"
	"fmt"
)

var poolApi *poolClientMap = nil
var poolRgm *poolClientMap = nil

func getClient(pool *poolClientMap) (c *idleClient, err error) {
	if pool == nil {
		err = createException("pool uninitialized！")
		return
	}
	c, err = pool.get()
	if err != nil {
		return
	}
	return
}

func getAPIClient() (client *bfss_api.BFSS_APIDClient, err error) {
	var c *idleClient
	c, err = getClient(poolApi)
	if err != nil {
		return
	}
	_client, ok := c.c.(*bfss_api.BFSS_APIDClient)
	if !ok {
		err = createException("conversion failed(BFSS_APIDClient)")
		return
	}
	client = _client
	return
}

func getREGMClient() (client *bfss_regm.BFSS_REGMDClient, err error) {
	var c *idleClient
	c, err = getClient(poolRgm)
	if err != nil {
		return
	}
	_client, ok := c.c.(*bfss_regm.BFSS_REGMDClient)
	if !ok {
		err = createException("conversion failed(BFSS_REGMDClient)")
		return
	}
	client = _client
	return
}

func getSNClientByPool(poolSn *poolClientMap) (client *bfss_sn.BFSS_SNDClient, err error) {
	var c *idleClient
	c, err = getClient(poolSn)
	if err != nil {
		return
	}
	_client, ok := c.c.(*bfss_sn.BFSS_SNDClient)
	if !ok {
		err = createException("conversion failed(BFSS_SNDClient)")
		return
	}
	client = _client
	return
}

func getSNClient(addr string, port int32) (client *bfss_sn.BFSS_SNDClient, pool *poolClientMap, err error) {
	keyAdrr := fmt.Sprintf("%s:%d", addr, port)
	snPoolMapLock.Lock()
	if snPoolMap == nil {
		snPoolMap = make(map[string]*poolClientMap)
	}
	if _, ok := snPoolMap[keyAdrr]; !ok {
		snPoolMap[keyAdrr] = newPoolMap(newThriftPool(addr, port, snMaxConn, snConnTimeout, snIdleTimeoutint, dialSn, close))
	}
	poolSn := snPoolMap[keyAdrr]
	snPoolMapLock.Unlock()
	client, err = getSNClientByPool(poolSn)
	if err != nil {
		return
	}
	pool = poolSn
	return
}

func CallAPI(f func(*bfss_api.BFSS_APIDClient) error) (err error) {
	recount := reconMax
	for ; recount > 0; recount-- {
		var c *bfss_api.BFSS_APIDClient
		c, err = getAPIClient()
		if err != nil {
			return
		}
		err = f(c)
		_ = poolApi.release(c, err != nil)
		if err == nil {
			return
		}
	}
	return
}

func callRgm(f func(client *bfss_regm.BFSS_REGMDClient) error) (err error) {
	recount := reconMax
	for ; recount > 0; recount-- {
		var c *bfss_regm.BFSS_REGMDClient
		c, err = getREGMClient()
		if err != nil {
			return
		}
		err = f(c)
		_ = poolRgm.release(c, err != nil)
		if err == nil {
			return err
		}
	}
	return
}
