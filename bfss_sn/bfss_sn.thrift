
// BFSS  Block File Storage System 块文件存储系统
// BFSS（存储系统微服务） BFSS内部调用 外部透明


include "../utils/bfss_result.thrift"

namespace cpp BFSS_SN
namespace go bfss.bfss_sn

struct READ_RESULT{
    1:required bfss_result.BFSS_RESULT Result=bfss_result.BFSS_RESULT.BFSS_UNKNOWN_ERROR;      // 错误代码
    2:optional binary Data           // 读取的数据
}

//返回 hash
struct HASH_RESULT{
    1:required bfss_result.BFSS_RESULT Result=bfss_result.BFSS_RESULT.BFSS_UNKNOWN_ERROR;      // 错误代码
    2:required string hash
}

struct BLKKEY_RESULT{
    1:required bfss_result.BFSS_RESULT Result=bfss_result.BFSS_RESULT.BFSS_UNKNOWN_ERROR;      // 错误代码
    2:optional binary Data           // 读取的数据
}


//---------------------------------------------------------------------------------------------
service BFSS_SND{
	// 函数名：GetVersion 获取服务版本信息(系统内部调用)
	// 参数：
	string GetVersion(),

    // 函数名：ManageMessage 消息管理  用于微服务之前的消息通知(系统内部调用)
	// 参数：
	// CmdId    自定义消息ID
	// Param    自定义消息命令参数
	// Data 	自定义消息数据
	// 返回值：
	//			参数 MESSAGE_RESULT 结构体定义
    bfss_result.MESSAGE_RESULT ManageMessage(1:bfss_result.BFSS_CMD CmdID,2:i32 Param,3:binary Data),

    // 函数名： WiteData 写数据
	// 参数：
	// index    写目标的块ID
	// offset   写目标块的起始偏移
	// data     写的具体数据
	// flag     可选
	// 返回值： 
	//			参考BFSS_RESULT结构体定义
    bfss_result.BFSS_RESULT WriteData(1:i32 index, 2:i32 offset, 3:binary data, 4:binary octx, 5:i32 flag),

    // 函数名： ReadData 读数据
	// 参数：
	// index    需要读取的数据所在块的ID
	// offset   需要读取的数据块的起始地址
	// size     需要读取的数据的大小
	// flag     可选
	// 返回值：
	//			参考READ_RESULT结构体定义
    READ_RESULT ReadData(1:i32 index, 2:i32 offset, 3:i32 size, 4:i32 flag),

    HASH_RESULT CompleteWriteObj(1:binary octx),

    BLKKEY_RESULT GetBlkKey(1:i32 index),
}
