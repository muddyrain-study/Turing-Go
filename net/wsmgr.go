package net

import "sync"

var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
}

type WsMgr struct {
	uc        sync.RWMutex
	userCache map[int]WSConn
}

func (m *WsMgr) UserLogin(conn WSConn, uid int, token string) {
	m.uc.Lock()
	defer m.uc.Unlock()
	oldConn := m.userCache[uid]
	if oldConn != nil {
		if oldConn != conn {
			oldConn.Push("robLogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
