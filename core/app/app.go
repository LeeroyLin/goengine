package app

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/closer"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/etcd"
	"github.com/LeeroyLin/goengine/core/module"
	"github.com/LeeroyLin/goengine/core/rpc"
	"github.com/LeeroyLin/goengine/iface"
	"sync"
)

// App 运行在节点上的应用
type App struct {
	Name              string
	Desc              string
	preModuleGroup    *module.ModuleGroup
	middleModuleGroup *module.ModuleGroup
	lateModuleGroup   *module.ModuleGroup
	AppHandler        iface.IAppHandler
	msgChanCapacity   int // 模块间消息通道容量
	closer            *closer.SigCloser
	sync.RWMutex
	RPC  iface.IRPC
	ETCD iface.IETCD
}

// NewApp 返回一个初始化的App
func NewApp(name, desc string) *App {
	a := &App{
		Name:            name,
		Desc:            desc,
		msgChanCapacity: 1024,
		RPC:             rpc.NewRPC(),
		ETCD:            etcd.NewETCD(),
		closer:          closer.NewSigCloser(),
	}

	return a
}

func (a *App) Init(preModules, modules, lateModules []iface.IModule) {
	elog.Info("[App] Init.", a.Name)

	if a.AppHandler != nil {
		a.AppHandler.OnBeforeInit()
	}

	// 先运行PreModule
	if preModules != nil && len(preModules) > 0 {
		a.preModuleGroup = module.NewModuleGroup(a, a, a, a.msgChanCapacity, a.closer.CloseChan)
		a.preModuleGroup.InitModules(preModules)
		a.preModuleGroup.RunModules()
	}

	if modules != nil && len(modules) > 0 {
		a.middleModuleGroup = module.NewModuleGroup(a, a, a, a.msgChanCapacity, a.closer.CloseChan)

		// 添加并初始化模块
		a.middleModuleGroup.InitModules(modules)
	}

	if lateModules != nil && len(lateModules) > 0 {
		a.lateModuleGroup = module.NewModuleGroup(a, a, a, a.msgChanCapacity, a.closer.CloseChan)

		// 添加并初始化模块
		a.lateModuleGroup.InitModules(lateModules)
	}

	if a.AppHandler != nil {
		a.AppHandler.OnAfterInit()
	}
}

// Run 运行应用
//
// param successCb 启动成功时的回调
// param finalCb 结束后的回调
func (a *App) Run(successCb func(), finalCb func()) {
	elog.Info("[App] Run.", a.Name, a.Desc)

	if a.AppHandler != nil {
		a.AppHandler.OnBeforeRun()
	}

	go func() {
		// 运行模块
		a.middleModuleGroup.RunModules()

		a.RPC.StartServe()

		if a.AppHandler != nil {
			a.AppHandler.OnAfterRun()
		}

		elog.Info("[App] app is running.")

		if successCb != nil {
			successCb()
		}

		// 运行后置模块
		a.lateModuleGroup.RunModules()

		select {
		case <-a.closer.CloseChan:
			a.doStop()
			break
		}
	}()

	a.closer.Listen(finalCb)
}

// Stop 停止应用
func (a *App) Stop() {
	select {
	case <-a.closer.CloseChan:
		return
	default:
		a.closer.Close()
	}
}

func (a *App) doStop() {
	elog.Info("[App] start stop app. ", a.Name)

	if a.AppHandler != nil {
		a.AppHandler.OnBeforeStop()
	}

	// 停止前处理
	a.lateModuleGroup.BeforeStopModules()
	a.middleModuleGroup.BeforeStopModules()

	if a.preModuleGroup != nil {
		a.preModuleGroup.BeforeStopModules()
	}

	go func() {
		// 停止模块
		a.lateModuleGroup.StopModules()
		a.middleModuleGroup.StopModules()

		if a.preModuleGroup != nil {
			a.preModuleGroup.StopModules()
		}

		a.RPC.ClearAll()
		a.ETCD.Stop()

		if a.AppHandler != nil {
			a.AppHandler.OnAfterStop()
		}

		elog.Info("[App] app stoped. ", a.Name)
	}()
}

// 获得模块
func (a *App) getModule(module string) (iface.IModule, error) {
	m, ok := a.middleModuleGroup.GetModule(module)

	if ok {
		return m, nil
	}

	if a.preModuleGroup != nil {
		m, ok = a.preModuleGroup.GetModule(module)

		if ok {
			return m, nil
		}
	}

	err := errors.New(fmt.Sprintln("[App] get msg center failed. can not find module.", a.Name, module))
	return nil, err
}

// 获得模块的消息中心
func (a *App) getModuleMsgCenter(module string) (iface.IMsgCenter, error) {
	a.RLock()
	defer a.RUnlock()

	m, err := a.getModule(module)

	if err != nil {
		return nil, err
	}

	return m.GetMsgCenter(), nil
}
