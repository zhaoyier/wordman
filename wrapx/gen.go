package wrapx

import (
	"sync"
)

var _mu sync.Mutex // protects the serviceMap

func getInfo() map[string][]genRouterInfo {
	_mu.Lock()
	defer _mu.Unlock()

	genInfo := _genInfo
	if _genInfoCnf.Tm > genInfo.Tm { // config to update more than coding
		genInfo = _genInfoCnf
	}

	mp := make(map[string][]genRouterInfo, len(genInfo.List))
	for _, v := range genInfo.List {
		tmp := v
		mp[tmp.HandFunName] = append(mp[tmp.HandFunName], tmp)
	}
	return mp
}
