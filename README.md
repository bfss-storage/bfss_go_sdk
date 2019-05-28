# 检出代码 
<strong>务必检出到 bfss，根目录为$GOPATH/src/bfss</strong>

```bash
git clone https://github.com/bfss-storage/bfss_go_sdk.git bfss
```

# SDK接口文件
```text
#$GOPATH/src/bfss
./bfss_sdk/     - bfss_sdk提供类似文件流的读写方式，特点是
                - 创建CreateObject()、写Write()、关闭Close()
                - 打开OpenObject()、读Read()、关闭Close()
./bfss_sdk_api/ - 提供bfss_api接口的文件读写方式，特点是
                - 创建CreateObject()、写Write()、完成CompleteObject()
                - 读Read()
```

# 生成thrift接口文件

./bfss_api/  
./bfss_sn/  
./bfss_regm/  
./utils/
这四个目录下的代码是基于*.thrift接口定义文件自动生成的，命令如下：

```bash
#pwd=$GOPATH/src
cd bfss/bfss_sdk
thrift --gen go --out ../ -r ./bfss_api/bfss_api.thrift 
thrift --gen go --out ../ -r ./bfss_sn/bfss_sn.thrift 
thrift --gen go --out ../ -r ./bfss_regm/bfss_regm.thrift

```

# 运行测试代码

```bash
# bfss_api的 读写操作
cd bfss/bfss_test/test_api
go build
./test_api

# 类似流/文件的 读写操作
cd bfss/bfss_test/test_stream
go build
./test_stream

```


# 内部测试环境信息
| IP：Port          | 说明                       | 备注                              |
| ----------------- |-------------------------- | --------------------------------- |
| 10.0.1.185:9090   | BFSS_REGM 注册管理服务      | .185 上部署了mongod 服务  .185也是bfss的CentOS编译环境         |
| 10.0.1.182:9091   | BFSS_sn 存储节点服务        |                                |
| 10.0.1.183:9092   | BFSS_sn 存储节点服务        |                                |
| 10.0.1.184:9091   | BFSS_sn 存储节点服务        |                                  |
| 10.0.1.185:9092   | BFSS_API BFSS对外接口服务   | .185 上部署了memcached服务          |
| 10.0.1.186:9092   | BFSS_API BFSS对外接口服务   | .186 上部署了memcached服务          |
| 10.0.1.119:30000  | BFSS_API 负载均衡服务器     |                                    |


# 正式环境信息 - for 朋友圈
| IP：Port          | 说明                       | 备注                              |
| ----------------- |-------------------------- | --------------------------------- |
| 10.6.55.10:9090   | BFSS_REGM 注册管理服务      |                                  |
| 10.6.55.11:9091   | BFSS_sn 存储节点服务        |  .11主  .13备                     |
| 10.6.55.13:9091   | BFSS_sn 存储节点服务        |                                 | 
| 10.6.55.12:9091   | BFSS_sn 存储节点服务        |  .12主  .14备                    |
| 10.6.55.14:9091   | BFSS_sn 存储节点服务        |                                  |
| 10.6.55.15:9092   | BFSS_API BFSS对外接口服务   | .15 上部署了memcached服务          |
|                   |                           |  MONGO=10.6.10.81:27017,10.6.10.82:27017,10.6.10.83:27017/?replicaSet=replSet |


