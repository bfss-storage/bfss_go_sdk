package bfss_sdk

func max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func min32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func pos_align_x(__pos int32, __align int32) int32 {
	__append := __pos % __align
	if __append < 0 {
		return __pos + -(__align + __append)
	}
	return __pos - __append
}

func pos_align_x64(__pos int64, __align int64) int64 {
	__append := __pos % __align
	if __append < 0 {
		return __pos + -(__align + __append)
	}
	return __pos - __append
}
