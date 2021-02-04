package wrapx

import (
	"strconv"
)

func (i ErrCode) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _ErrCode_name_0[_ErrCode_index_0[i]:_ErrCode_index_0[i+1]]
	case 101 <= i && i <= 102:
		i -= 101
		return _ErrCode_name_1[_ErrCode_index_1[i]:_ErrCode_index_1[i+1]]
	case i == 200:
		return _ErrCode_name_2
	case i == 405:
		return _ErrCode_name_3
	case 1001 <= i && i <= 1025:
		i -= 1001
		return _ErrCode_name_4[_ErrCode_index_4[i]:_ErrCode_index_4[i+1]]
	case 1027 <= i && i <= 1036:
		i -= 1027
		return _ErrCode_name_5[_ErrCode_index_5[i]:_ErrCode_index_5[i+1]]
	case 1043 <= i && i <= 1087:
		i -= 1043
		return _ErrCode_name_6[_ErrCode_index_6[i]:_ErrCode_index_6[i+1]]
	case 2001 <= i && i <= 2002:
		i -= 2001
		return _ErrCode_name_7[_ErrCode_index_7[i]:_ErrCode_index_7[i+1]]
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func _tryRegisteryCode(code ErrCode) {
	_mp.LoadOrStore(code.String(), int(code))
}
