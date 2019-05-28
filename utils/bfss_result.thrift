
namespace cpp bfss
namespace go bfss.utils
// 错误代码
enum BFSS_RESULT {
    //-1 --- -49 BFSS 内部错误使用
    //-50 --- -100 外部接口使用
    //-100 --- -254 预留可选
    BFSS_DATA_READ_INCOMPLETE   = 1,    // 读操作请求的字节数过大，返回的字符数不够请求的数量。
    BFSS_DATA_WRITE_INCOMPLETE  = 2,    // 写操作发送的字符串超长，实际写入的内容不足发送的数据。
    BFSS_SUCCESS                = 0,    // 调用成功
    BFSS_UNKNOWN_ERROR          = -255, // 未知错误

    BFSS_PARAM_ERROR            = -51,  // 参数错误
    BFSS_SCHEME_ERROR           = -52,  // 配置错误

    BFSS_NO_SPACE               = -53,  // 磁盘空间不足，分配失败
    BFSS_NO_MEMORY              = -54,  // 内存不足，分配失败
    BFSS_TIMEOUT                = -55,  // 超时

    BFSS_NOTFOUND               = -60,  // 对象不存在
    BFSS_DUPLICATED             = -61,  // 对象已存在，创建失败
    BFSS_COMPLETED              = -62,  // 对象已写入完成
    BFSS_INCOMPLETED            = -63,  // 对象未写入完成

    BFSS_DATA_WRITE_FAILED      = -70,  // WRITE数据失败
    BFSS_DATA_READ_FAILED       = -71,  // READ数据失败
    BFSS_DATA_COMPLETE_FAILED   = -72,  // 写入完成时计算Hash失败
    BFSS_DATA_UNINITED          = -73,  // 块未被使用过
    BFSS_NO_SND                 = -74,  //
    BFSS_SND_SLAVE_SYNC_FAILED  = -75,  // 从节点写入失败
}


enum BFSS_CMD {
    CMD_SN_MASTER_OFFLINE   = 1,    // APID  -> RegM
    CMD_SN_SLAVE_OFFLINE    = 2,    // 主SN  -> RegM

    CMD_SN_BLK_READ         = 11,   //(om)主/备SN ->  主/备SN
    CMD_SN_BLK_WRITE        = 12,   //主SN   ->  备SN
    CMD_SN_RESTORED         = 13,

    // add cmd above
    CMD_NO_INSTANCE = 32767
}

enum BFSS_SN_NODE_TYPES{
    Standalone,
    Master,
    Slave
}

enum BFSS_SN_NODE_STATUS{
    Writable,
    Readonly,
    Unavailable
}

//消息信息 BFSS内部使用
struct BFSS_MESSAGE {
    1:optional BFSS_CMD Cmd             // 自定义命令
    2:optional i32 Param                // 自定义命令参数
    3:optional binary Data              // 自定义命令所带数据
}
//返回 自定义命令消息
struct MESSAGE_RESULT {
    1:required BFSS_RESULT  Result=bfss_result.BFSS_RESULT.BFSS_UNKNOWN_ERROR;      // 错误代码
    2:optional BFSS_MESSAGE Resp        // 自定义消息响应
}
