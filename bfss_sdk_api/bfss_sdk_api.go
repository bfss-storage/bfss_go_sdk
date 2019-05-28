package bfss_sdk_api

import (
	"bfss/bfss_api"
	"bfss/bfss_sdk"
	"bfss/utils"
	"context"
)


func GetVersion() (r string, err error) {
	var _args1 bfss_api.BFSS_APIDGetVersionArgs
	var _result2 bfss_api.BFSS_APIDGetVersionResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "GetVersion", &_args1, &_result2)
	})
	if err != nil{
		return
	}
	return _result2.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Size
//  - Flag
//  - Tag
func CreateObject( oid string, size int32, flag int32, tag string) (r utils.BFSS_RESULT, err error) {
	var _args3 bfss_api.BFSS_APIDCreateObjectArgs
	_args3.Oid = oid
	_args3.Size = size
	_args3.Flag = flag
	_args3.Tag = tag
	var _result4 bfss_api.BFSS_APIDCreateObjectResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "CreateObject", &_args3, &_result4)
	})
	if err != nil{
		return
	}
	return _result4.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Size
//  - Flag
//  - Tag
func CreateObjectEx( oid string, size int32, flag int32, tag string) (r *bfss_api.OBJECT_INFO_EX_RESULT, err error) {
	var _args5 bfss_api.BFSS_APIDCreateObjectExArgs
	_args5.Oid = oid
	_args5.Size = size
	_args5.Flag = flag
	_args5.Tag = tag
	var _result6 bfss_api.BFSS_APIDCreateObjectExResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "CreateObjectEx", &_args5, &_result6)
	})
	if err != nil{
		return
	}
	return _result6.GetSuccess(), nil
}

// Parameters:
//  - Oid
func DeleteObject( oid string) (r utils.BFSS_RESULT, err error) {
	var _args7 bfss_api.BFSS_APIDDeleteObjectArgs
	_args7.Oid = oid
	var _result8 bfss_api.BFSS_APIDDeleteObjectResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "DeleteObject", &_args7, &_result8)
	})
	if err != nil{
		return
	}
	return _result8.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Offset
//  - Data
func Write( oid string, offset int32, data []byte) (r utils.BFSS_RESULT, err error) {
	var _args9 bfss_api.BFSS_APIDWriteArgs
	_args9.Oid = oid
	_args9.Offset = offset
	_args9.Data = data
	var _result10 bfss_api.BFSS_APIDWriteResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "Write", &_args9, &_result10)
	})
	if err != nil{
		return
	}
	return _result10.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Newsize_
func ResizeObject( oid string, newsize int32) (r utils.BFSS_RESULT, err error) {
	var _args11 bfss_api.BFSS_APIDResizeObjectArgs
	_args11.Oid = oid
	_args11.Newsize_ = newsize
	var _result12 bfss_api.BFSS_APIDResizeObjectResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "ResizeObject", &_args11, &_result12)
	})
	if err != nil{
		return
	}
	return _result12.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Newoid_
//  - Newtag_
func ResetObjectId( oid string, newoid string, newtag string) (r utils.BFSS_RESULT, err error) {
	var _args13 bfss_api.BFSS_APIDResetObjectIdArgs
	_args13.Oid = oid
	_args13.Newoid_ = newoid
	_args13.Newtag_ = newtag
	var _result14 bfss_api.BFSS_APIDResetObjectIdResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "ResetObjectId", &_args13, &_result14)
	})
	if err != nil{
		return
	}
	return _result14.GetSuccess(), nil
}

// Parameters:
//  - Oid
func CompleteObject( oid string) (r utils.BFSS_RESULT, err error) {
	var _args15 bfss_api.BFSS_APIDCompleteObjectArgs
	_args15.Oid = oid
	var _result16 bfss_api.BFSS_APIDCompleteObjectResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "CompleteObject", &_args15, &_result16)
	})
	if err != nil{
		return
	}
	return _result16.GetSuccess(), nil
}

// Parameters:
//  - Oid
func GetObjectInfo( oid string) (r *bfss_api.OBJECT_INFO_RESULT, err error) {
	var _args17 bfss_api.BFSS_APIDGetObjectInfoArgs
	_args17.Oid = oid
	var _result18 bfss_api.BFSS_APIDGetObjectInfoResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "GetObjectInfo", &_args17, &_result18)
	})
	if err != nil{
		return
	}
	return _result18.GetSuccess(), nil
}

// Parameters:
//  - Oid
func GetObjectInfoEx( oid string) (r *bfss_api.OBJECT_INFO_EX_RESULT, err error) {
	var _args19 bfss_api.BFSS_APIDGetObjectInfoExArgs
	_args19.Oid = oid
	var _result20 bfss_api.BFSS_APIDGetObjectInfoExResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "GetObjectInfoEx", &_args19, &_result20)
	})
	if err != nil{
		return
	}
	return _result20.GetSuccess(), nil
}

// Parameters:
//  - Hash
//  - Size
//  - Head
func ObjectLockHasHash( hash string, size int32, head []byte) (r utils.BFSS_RESULT, err error) {
	var _args21 bfss_api.BFSS_APIDObjectLockHasHashArgs
	_args21.Hash = hash
	_args21.Size = size
	_args21.Head = head
	var _result22 bfss_api.BFSS_APIDObjectLockHasHashResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "ObjectLockHasHash", &_args21, &_result22)
	})
	if err != nil{
		return
	}
	return _result22.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Hash
//  - Size
//  - Head
//  - Flag
//  - Tag
func CreateObjectLink( oid string, hash string, size int32, head []byte, flag int32, tag string) (r utils.BFSS_RESULT, err error) {
	var _args23 bfss_api.BFSS_APIDCreateObjectLinkArgs
	_args23.Oid = oid
	_args23.Hash = hash
	_args23.Size = size
	_args23.Head = head
	_args23.Flag = flag
	_args23.Tag = tag
	var _result24 bfss_api.BFSS_APIDCreateObjectLinkResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "CreateObjectLink", &_args23, &_result24)
	})
	if err != nil{
		return
	}
	return _result24.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Size
//  - Offset
func Read( oid string, size int32, offset int32) (r *bfss_api.READ_RESULT, err error) {
	var _args25 bfss_api.BFSS_APIDReadArgs
	_args25.Oid = oid
	_args25.Size = size
	_args25.Offset = offset
	var _result26 bfss_api.BFSS_APIDReadResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "Read", &_args25, &_result26)
	})
	if err != nil{
		return
	}
	return _result26.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Size
//  - Offset
func ReadBlk( oid string, size int32, offset int32) (r *bfss_api.READ_RESULT, err error) {
	var _args27 bfss_api.BFSS_APIDReadBlkArgs
	_args27.Oid = oid
	_args27.Size = size
	_args27.Offset = offset
	var _result28 bfss_api.BFSS_APIDReadBlkResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "ReadBlk", &_args27, &_result28)
	})
	if err != nil{
		return
	}
	return _result28.GetSuccess(), nil
}

// Parameters:
//  - Oid
func GetObjectBlksInfo( oid string) (r *bfss_api.OBJECT_BLKS_RESULT, err error) {
	var _args29 bfss_api.BFSS_APIDGetObjectBlksInfoArgs
	_args29.Oid = oid
	var _result30 bfss_api.BFSS_APIDGetObjectBlksInfoResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "GetObjectBlksInfo", &_args29, &_result30)
	})
	if err != nil{
		return
	}
	return _result30.GetSuccess(), nil
}

// Parameters:
//  - Oid
//  - Offset
func GetObjectBlkKey( oid string, offset int32) (r *bfss_api.BLK_KEY_RESULT, err error) {
	var _args31 bfss_api.BFSS_APIDGetObjectBlkKeyArgs
	_args31.Oid = oid
	_args31.Offset = offset
	var _result32 bfss_api.BFSS_APIDGetObjectBlkKeyResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "GetObjectBlkKey", &_args31, &_result32)
	})
	if err != nil{
		return
	}
	return _result32.GetSuccess(), nil
}

// Parameters:
//  - CmdId
//  - Param
//  - Data
func ManageMessage( CmdId utils.BFSS_CMD, Param int32, Data []byte) (r *utils.MESSAGE_RESULT, err error) {
	var _args33 bfss_api.BFSS_APIDManageMessageArgs
	_args33.CmdId = CmdId
	_args33.Param = Param
	_args33.Data = Data
	var _result34 bfss_api.BFSS_APIDManageMessageResult
	err = bfss_sdk.CallAPI(func(c *bfss_api.BFSS_APIDClient) error{
		return c.Client_().Call(context.Background(), "ManageMessage", &_args33, &_result34)
	})
	if err != nil{
		return
	}
	return _result34.GetSuccess(), nil
}