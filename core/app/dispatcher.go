package app

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/def"
)

// Call 同步调用模块
func (a *App) Call(module string, req def.ICommReq) (interface{}, error) {
	mc, err := a.getModuleMsgCenter(module)
	if err != nil {
		elog.Error("[App] call err.", a.Name, module, err)
		return nil, err
	}

	return mc.Call(req)
}

// CallAsync 异步调用模块
func (a *App) CallAsync(module string, req def.ICommReq, cb def.MsgRespHandler) error {
	mc, err := a.getModuleMsgCenter(module)
	if err != nil {
		elog.Error("[App] call async err.", a.Name, module, err)
		return err
	}

	mc.CallAsync(req, cb)
	return nil
}

// Cast 向模块投递消息
func (a *App) Cast(module string, req def.ICommReq) {
	mc, err := a.getModuleMsgCenter(module)
	if err != nil {
		elog.Error("[App] cast err.", a.Name, module, err)
		return
	}

	mc.Cast(req)
}
