package test

import (
	"github.com/LeeroyLin/goengine/core/app"
	"github.com/LeeroyLin/goengine/core/config"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/module"
	"github.com/LeeroyLin/goengine/def"
	"github.com/LeeroyLin/goengine/iface"
	"testing"
	"time"
)

var Config = &config.ConfBase{}

type TestCommReq struct {
	CommId uint32
}

func (r *TestCommReq) GetCommId() uint32 {
	return r.CommId
}

type TestApp1 struct {
	app.App
}

type TestMgr struct {
}

func (mgr *TestMgr) OnInit() {
	elog.Info("[Test] TestMgr OnInit")
}

func (mgr *TestMgr) OnRun() {
	elog.Info("[Test] TestMgr OnRun")
}

func (mgr *TestMgr) OnStop() {
	elog.Info("[Test] TestMgr OnStop")
}

type TestModule1 struct {
	module.Module
}

func NewTestModule1() *TestModule1 {
	m := &TestModule1{
		Module: module.NewModule("module1"),
	}
	m.Module.SetLife(m)

	return m
}

func (tm *TestModule1) OnInit() {
	elog.Info("[Test] TestModule1 OnInit")
}

func (tm *TestModule1) OnRun() error {
	elog.Info("[Test] TestModule1 OnRun")
	return nil
}

func (tm *TestModule1) OnBeforeStop() error {
	elog.Info("[Test] TestModule1 OnBeforeStop")
	return nil
}

func (tm *TestModule1) OnStop() error {
	elog.Info("[Test] TestModule1 OnStop")
	return nil
}

func (tm *TestModule1) OnSetMgrs() []iface.IMgr {
	return []iface.IMgr{
		&TestMgr{},
	}
}

func (tm *TestModule1) OnRegMsgHandler(msgCenter iface.IMsgCenter) {
	msgCenter.AddHandler(11, func(isSync bool, req def.ICommReq) (interface{}, error) {
		elog.Info("[Test] back 11")
		return nil, nil
	})
}

func (tm *TestModule2) OnInit() {
	elog.Info("[Test] TestModule2 OnInit")
}

func (tm *TestModule2) OnRun() error {
	elog.Info("[Test] TestModule2 OnRun")
	return nil
}

func (tm *TestModule2) OnBeforeStop() error {
	elog.Info("[Test] TestModule1 OnBeforeStop")
	return nil
}

func (tm *TestModule2) OnStop() error {
	elog.Info("[Test] TestModule2 OnStop")
	return nil
}

func (tm *TestModule2) OnSetMgrs() []iface.IMgr {
	return []iface.IMgr{}
}

func (tm *TestModule2) OnRegMsgHandler(msgCenter iface.IMsgCenter) {
	msgCenter.AddHandler(22, func(isSync bool, req def.ICommReq) (interface{}, error) {
		elog.Info("[Test] back 22")
		tm.GetDispatcher().Call("module1", &TestCommReq{CommId: 11})
		return nil, nil
	})
}

type TestModule2 struct {
	module.Module
}

func NewTestModule2() *TestModule2 {
	m := &TestModule2{
		Module: module.NewModule("module2"),
	}
	m.Module.SetLife(m)

	return m
}

func TestApp(t *testing.T) {
	// 装载配置
	Config.Setup(Config, "")

	a := &TestApp1{
		App: *app.NewApp(Config.Name, Config.Desc),
	}

	a.Init([]iface.IModule{}, []iface.IModule{
		NewTestModule1(),
		NewTestModule2(),
	}, []iface.IModule{})

	a.Run(func() {
		a.Call("module1", &TestCommReq{
			CommId: 11,
		})

		a.CallAsync("module2", &TestCommReq{
			CommId: 22,
		}, func(resp interface{}, err error) {
			elog.Info("[Test] 22 res:", resp, err)
		})

		a.Cast("module1", &TestCommReq{
			CommId: 11,
		})

		a.Cast("module1", &TestCommReq{
			CommId: 33,
		})
	}, func() {})

	time.Sleep(1 * time.Second)
	a.Stop()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
