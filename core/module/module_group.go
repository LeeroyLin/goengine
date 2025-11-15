package module

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface"
	"sync"
)

type ModuleGroup struct {
	modules         map[string]iface.IModule // 所有模块
	rpcGetter       iface.IRPCGetter
	etcdGetter      iface.IETCDGetter
	dispatcher      iface.IDispatcher
	msgChanCapacity int
	closeChan       chan interface{}
	wg              *sync.WaitGroup // 等待组
}

func NewModuleGroup(dispatcher iface.IDispatcher, rpcGetter iface.IRPCGetter, etcdGetter iface.IETCDGetter, msgChanCapacity int, closeChan chan interface{}) *ModuleGroup {
	g := &ModuleGroup{
		modules:         make(map[string]iface.IModule),
		dispatcher:      dispatcher,
		rpcGetter:       rpcGetter,
		etcdGetter:      etcdGetter,
		msgChanCapacity: msgChanCapacity,
		closeChan:       closeChan,
		wg:              &sync.WaitGroup{},
	}

	return g
}

func (g *ModuleGroup) InitModules(modules []iface.IModule) {
	// 添加模块
	for _, m := range modules {
		n := m.GetName()
		g.modules[n] = m
	}

	// 初始化模块，添加之后统一调用
	for _, m := range g.modules {
		m.DoInit(g.dispatcher, g.rpcGetter, g.etcdGetter, g.msgChanCapacity, g.closeChan)
	}
}

func (g *ModuleGroup) RunModules() {
	for _, m := range g.modules {
		g.wg.Add(1)
		go func() {
			// 运行模块
			err := m.DoRun()
			g.wg.Done()
			if err != nil {
				elog.Error("[App] run module err: ", m.GetName(), err)
				return
			}
		}()
	}
	g.wg.Wait()
}

func (g *ModuleGroup) StopModules() {
	for _, m := range g.modules {
		g.wg.Add(1)
		go func() {
			err := m.DoStop()
			g.wg.Done()
			if err != nil {
				elog.Error("[App] stop module err: ", m.GetName(), err)
				return
			}
		}()
	}
	g.wg.Wait()
}

func (g *ModuleGroup) GetModule(module string) (iface.IModule, bool) {
	m, ok := g.modules[module]
	return m, ok
}
