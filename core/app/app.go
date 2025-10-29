package app

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/module"
	"github.com/LeeroyLin/goengine/core/rpc"
	"github.com/LeeroyLin/goengine/iface"
	"sync"
)

// App 运行在节点上的应用
type App struct {
	Name              string
	Desc              string
	closeChan         chan interface{} // 用于关闭的通道
	preModuleGroup    *module.ModuleGroup
	middleModuleGroup *module.ModuleGroup
	markStop          bool // 标记停止
	preparing         bool // 准备中
	running           bool // 运行中
	appHandler        iface.IAppHandler
	msgChanCapacity   int // 模块间消息通道容量
	sync.RWMutex
	RPC iface.IRPC
}

// NewApp 返回一个初始化的App
func NewApp(name, desc string, appHandler iface.IAppHandler) *App {
	a := &App{
		Name:            name,
		Desc:            desc,
		markStop:        false,
		preparing:       false,
		running:         false,
		msgChanCapacity: 1024,
		appHandler:      appHandler,
		RPC:             rpc.NewRPC(),
	}

	return a
}

func (a *App) Init(preModules, modules []iface.IModule) {
	elog.Info("[App] Init.", a.Name)

	a.appHandler.OnBeforeInit()

	if a.closeChan == nil {
		a.closeChan = make(chan interface{})
	}

	// 先运行PreModule
	if preModules != nil && len(preModules) > 0 {
		a.preModuleGroup = module.NewModuleGroup(a, a, a.msgChanCapacity, a.closeChan)
		a.preModuleGroup.InitModules(preModules)
		a.preModuleGroup.RunModules()
	}

	if modules != nil && len(modules) > 0 {
		a.middleModuleGroup = module.NewModuleGroup(a, a, a.msgChanCapacity, a.closeChan)

		// 添加并初始化模块
		a.middleModuleGroup.InitModules(modules)
	}

	a.appHandler.OnAfterInit()
}

// Run 运行应用
//
// param successCb 启动成功时的回调
func (a *App) Run(successCb func()) {
	a.Lock()
	defer a.Unlock()

	if a.running {
		elog.Error("[App] app is running. ", a.Name)
		return
	}

	if a.preparing {
		elog.Error("[App] app is preparing. ", a.Name)
		return
	}

	if a.markStop {
		elog.Error("[App] app already mark stop. ", a.Name)
		return
	}

	elog.Info("[App] Run.", a.Name, a.Desc)

	a.preparing = true

	a.appHandler.OnBeforeRun()

	go func() {
		// 运行模块
		a.middleModuleGroup.RunModules()

		a.Lock()
		a.preparing = false
		a.running = true
		a.Unlock()

		a.RPC.StartServe()

		a.appHandler.OnAfterRun()

		elog.Info("[App] app is running.")

		if successCb != nil {
			successCb()
		}

		select {
		case <-a.closeChan:
			a.doStop()
			break
		}
	}()
}

// Stop 停止应用
func (a *App) Stop() {
	select {
	case <-a.closeChan:
		return
	default:
		close(a.closeChan)
	}
}

func (a *App) doStop() {
	a.Lock()

	if a.preparing {
		a.markStop = true
		a.Unlock()
		return
	}

	if !a.running {
		elog.Error("[App] can not stop app without running. ", a.Name)
		a.Unlock()
		return
	}

	a.markStop = true
	a.running = false

	a.Unlock()

	elog.Info("[App] start stop app. ", a.Name)

	a.appHandler.OnBeforeStop()

	go func() {
		// 停止模块
		a.middleModuleGroup.StopModules()

		if a.preModuleGroup != nil {
			a.preModuleGroup.StopModules()
		}

		a.RPC.ClearAll()

		a.Lock()
		a.markStop = false
		a.Unlock()

		a.appHandler.OnAfterStop()

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
