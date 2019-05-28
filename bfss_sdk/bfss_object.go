package bfss_sdk

import (
	"bfss/bfss_api"
	"bfss/bfss_regm"
	"bfss/bfss_sn"
	"bfss/utils"
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type storeInfo struct {
	_vid int32
	_vdx int32
	_vsf int32
}

type mode int

const (
	_ mode = iota
	unknown
	open
	crate
)

type whence int

const (
	_       whence = iota
	SeekSet        = 0
	SeekCur        = 1
	SeekEnd        = 2
)

type object struct {
	_sm   mode
	_size int32
	_time int64
	_flag int32
	_sof  int32
	_bks  int32
	_oid  string
	_tag  string
	_hash string
	_cxt  []byte
	_si   storeInfo
	_vi   *bfss_regm.VOLUME_INFO
	_pn   *poolClientMap
	_sc   *bfss_sn.BFSS_SNDClient
}

//TransferBytesMax
const TransferBytesMax int = dataSizeMax

func (o *object) Read(bytes int) ([]byte, error) {
	if o._sm != open {
		return nil, AccessReadException
	}
	if bytes < 0 || bytes > TransferBytesMax {
		return nil, InvalidSizeException
	}
	if o._sof == o._size {
		return nil, FileEofException
	}
	index := o._si._vdx
	pos := (int64(index) * int64(o._bks)) + int64(o._si._vsf+o._sof)
	readBytes := min32(o._size-o._sof, int32(bytes))
	if readBytes == 0 {
		return make([]byte, 0), nil
	}
	off := int32(pos - pos_align_x64(pos, int64(o._bks)))
	pos = pos_align_x64(pos, int64(o._bks))
	index = int32(pos / int64(o._bks))
	var r *bfss_sn.READ_RESULT
	err := o.callSn(func() error {
		_r, err := o._sc.ReadData(context.Background(), index, off, readBytes, 1)
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return nil, err
	}
	if r.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
		return nil, createException(r.Result_.String())
	}
	o._sof += readBytes
	return r.Data, nil
}

func (o *object) Write(bytes []byte) (int, error) {
	if o._sm != crate {
		return -1, AccessWriteException
	}
	if o._sof == o._size {
		return -1, AlreadyWrittenException
	}
	if nil == bytes || len(bytes) > TransferBytesMax {
		return -1, InvalidDataException
	}
	index := o._si._vdx
	pos := (int64(index) * int64(o._bks)) + int64(o._si._vsf+o._sof)
	writeBytes := min32(o._size-o._sof, int32(len(bytes)))
	bytes = bytes[0:writeBytes]
	if len(bytes) == 0 {
		return 0, nil
	}
	off := int32(pos - pos_align_x64(pos, int64(o._bks)))
	pos = pos_align_x64(pos, int64(o._bks))
	index = int32(pos / int64(o._bks))
	var r utils.BFSS_RESULT
	err := o.callSn(func() error {
		_r, err := o._sc.WriteData(context.Background(), index, off, bytes, o._cxt, 0)
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return -1, err
	}
	if r != utils.BFSS_RESULT_BFSS_SUCCESS {
		return -1, createException(r.String())
	}
	o._sof += writeBytes
	return int(writeBytes), nil
}

func (o *object) Seek(offset int, whence int) (int, error) {
	if o._sm == unknown {
		return -1, NotOpenException
	}
	switch whence {
	case SeekSet:
		break
	case SeekCur:
		offset += int(o._sof)
		break
	case SeekEnd:
		offset = int(o._size) - offset
		break
	default:
		return -1, InvalidWhenceException
	}
	if int(o._size) < offset || offset < 0 {
		return -1, InvalidOffsetException
	}
	r := o._sof
	o._sof = int32(offset)
	return int(r), nil
}

func (o *object) Resize(newSize int) error {
	if o._sm != crate {
		return AccessResizeException
	}
	var r utils.BFSS_RESULT
	err := CallAPI(func(c *bfss_api.BFSS_APIDClient) error {
		_r, err := c.ResizeObject(context.Background(), o._oid, int32(newSize))
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return err
	}
	if r != utils.BFSS_RESULT_BFSS_SUCCESS {
		err = createException(r.String())
	}
	o._size = int32(newSize)
	return err
}

func (o *object) Rename(newMame string, newTag string) error {
	if o._sm != open {
		return AccessRenameException
	}
	var r utils.BFSS_RESULT
	err := CallAPI(func(c *bfss_api.BFSS_APIDClient) error {
		_r, err := c.ResetObjectId(context.Background(), o._oid, newMame, newTag)
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return err
	}
	if r != utils.BFSS_RESULT_BFSS_SUCCESS {
		err = createException(r.String())
	}
	o._tag = newTag
	o._oid = newMame
	return err
}

func (o *object) GetSize() (int, error) {
	if o._sm == unknown {
		return -1, NotOpenException
	}
	return int(o._size), nil
}

func (o *object) Eof() (bool, error) {
	if o._sm == unknown {
		return false, NotOpenException
	}
	return (o._size - o._sof) == 0, nil
}

func (o *object) GetTime() (r time.Time, err error) {
	if o._sm == unknown {
		err = NotOpenException
		return
	}
	r = time.Unix(o._time/1000, 0)
	return
}

func (o *object) GetFlag() (int, error) {
	if o._sm == unknown {
		return 0, NotOpenException
	}
	return int(o._flag), nil
}

func (o *object) GetTag() (r string, err error) {
	if o._sm == unknown {
		err = NotOpenException
		return
	}
	r = o._tag
	return
}

func (o *object) Close() (err error) {
	if o._sm != unknown {
		defer func() {
			o._pn = nil
			o._sc = nil
		}()
		if o._sm == crate {
			var r utils.BFSS_RESULT
			err = CallAPI(func(c *bfss_api.BFSS_APIDClient) error {
				_r, err := c.CompleteObject(context.Background(), o._oid)
				r = _r
				return err
			})
			if err != nil {
				return err
			}
			if r != utils.BFSS_RESULT_BFSS_SUCCESS {
				err = createException(r.String())
			}
		}
		o._sm = unknown
	}
	return
}

func (o *object) getSnCUri(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return toBfssException(err)
	}
	if strings.ToLower(u.Scheme) != "snd" {
		return InvalidSchemeException
	}
	var port int32 = 0
	if len(u.Port()) == 0 {
		port = 9092
	} else {
		_port, err := strconv.Atoi(u.Port())
		if err != nil {
			return toBfssException(err)
		}
		port = int32(_port)
	}
	sc, pn, err := getSNClient(u.Hostname(), port)
	if err != nil {
		return err
	}
	o._pn = pn
	o._sc = sc
	return nil
}

func (o *object) getSnC() (err error) {
reGetPool:
	if o._pn == nil {
		err = o.getSnCUri(o._vi.URI)
		if err == nil {
			return
		}
		if o._vi.Type != utils.BFSS_SN_NODE_TYPES_Master || o._vi.Slave == nil {
			return
		}
		if o._sm != open {
			return
		}
		err = o.getSnCUri(o._vi.Slave.URI)
		return
	}
	if o._sc, err = getSNClientByPool(o._pn); err == nil {
		return
	}
	o._pn = nil
	goto reGetPool
}

func (o *object) callSn(f func() error) error {
	recount := reconMax
	var err error
	for ; recount > 0; recount-- {
		err = o.getSnC()
		if err != nil {
			return err
		}
		err = f()
		_ = o._pn.release(o._sc, err != nil)
		if err == nil {
			o._sc = nil
			return nil
		}
		o._pn = nil
		o._sc = nil
	}
	return err
}

func CreateObject(oid string, size int32, flag int32, tag string) (*object, error) {
	var r *bfss_api.OBJECT_INFO_EX_RESULT
	cm := crate
	err := CallAPI(func(c *bfss_api.BFSS_APIDClient) error {
		_r, err := c.CreateObjectEx(context.Background(), oid, size, flag, tag)
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return nil, err
	}
	if r.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
		return nil, createException(r.Result_.String())
	}

	var rr *bfss_regm.VOLUME_RESULT
	err = callRgm(func(c *bfss_regm.BFSS_REGMDClient) error {
		_r, err := c.GetVolumeInfo(context.Background(), r.ObjectVolInfo.VolumeId)
		if err != nil {
			return err
		}
		rr = _r
		return nil
	})
	if err != nil {
		return nil, err
	}
	if rr.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
		return nil, createException(rr.Result_.String())
	}
	return &object{_sm: cm,
		_si: storeInfo{
			_vid: r.ObjectVolInfo.VolumeId,
			_vdx: r.ObjectVolInfo.BeginIndex,
			_vsf: r.ObjectVolInfo.BeginOffset,
		},
		_size: r.ObjectInfoEx.ObjectSize,
		_tag:  r.ObjectInfoEx.ObjectTag,
		_flag: r.ObjectInfoEx.ObjectFlag,
		_bks:  r.ObjectInfoEx.ObjectBlkSize,
		_time: r.ObjectInfoEx.CreateTime,
		_cxt:  r.ObjectInfoEx.CreateCtx,
		_hash: r.ObjectInfoEx.Hash,
		_oid:  oid,
		_sof:  0,
		_vi:   rr.Volume,
		_sc:   nil,
		_pn:   nil,
	}, nil
}

func OpenObject(oid string) (*object, error) {
	var r *bfss_api.OBJECT_INFO_EX_RESULT
	om := open
	err := CallAPI(func(c *bfss_api.BFSS_APIDClient) error {
		_r, err := c.GetObjectInfoEx(context.Background(), oid)
		if err != nil {
			return err
		}
		r = _r
		return nil
	})
	if err != nil {
		return nil, toBfssException(err)
	}
	if r.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
		return nil, createException(r.Result_.String())
	}
	if !r.ObjectInfoEx.Complete {
		om = crate
	}
	var rr *bfss_regm.VOLUME_RESULT
	err = callRgm(func(c *bfss_regm.BFSS_REGMDClient) error {
		_r, err := c.GetVolumeInfo(context.Background(), r.ObjectVolInfo.VolumeId)
		if err != nil {
			return err
		}
		rr = _r
		return nil
	})
	if err != nil {
		return nil, toBfssException(err)
	}
	if rr.Result_ != utils.BFSS_RESULT_BFSS_SUCCESS {
		return nil, createException(rr.Result_.String())
	}
	return &object{_sm: om,
		_si: storeInfo{
			_vid: r.ObjectVolInfo.VolumeId,
			_vdx: r.ObjectVolInfo.BeginIndex,
			_vsf: r.ObjectVolInfo.BeginOffset,
		},
		_size: r.ObjectInfoEx.ObjectSize,
		_tag:  r.ObjectInfoEx.ObjectTag,
		_flag: r.ObjectInfoEx.ObjectFlag,
		_bks:  r.ObjectInfoEx.ObjectBlkSize,
		_time: r.ObjectInfoEx.CreateTime,
		_cxt:  r.ObjectInfoEx.CreateCtx,
		_hash: r.ObjectInfoEx.Hash,
		_oid:  oid,
		_sof:  0,
		_vi:   rr.Volume,
		_sc:   nil,
		_pn:   nil,
	}, nil
}
