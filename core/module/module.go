package module

import (
	"github.com/LeeroyLin/goengine/core/msgcenter"
	"github.com/LeeroyLin/goengine/def"
	"github.com/LeeroyLin/goengine/iface"
)

type Module struct {
	name       string
	msgCenter  iface.IMsgCenter // 模块内的消息中心
	rpcGetter  iface.IRPCGetter
	etcdGetter iface.IETCDGetter
	dispatcher iface.IDispatcher // 模块间消息分发器
	mgrs       []iface.IMgr      // 管理器
	life       iface.IModuleLife
}

func NewModule(name string) Module {
	m := Module{
		name: name,
		mgrs: make([]iface.IMgr, 0),
	}

	return m
}

func (m *Module) SetLife(life iface.IModuleLife) {
	m.life = life
}

func (m *Module) GetName() string {
	return m.name
}

// GetMsgCenter 获得模块内消息中心
func (m *Module) GetMsgCenter() iface.IMsgCenter {
	return m.msgCenter
}

// GetDispatcher 获得模块间消息分发器
func (m *Module) GetDispatcher() iface.IDispatcher {
	return m.dispatcher
}

func (m *Module) DoInit(dispatcher iface.IDispatcher, rpcGetter iface.IRPCGetter, etcdGetter iface.IETCDGetter, msgChanCapacity int, closeChan chan interface{}) {
	m.dispatcher = dispatcher
	m.rpcGetter = rpcGetter
	m.etcdGetter = etcdGetter
	m.msgCenter = msgcenter.NewMsgCenter(m.name, msgChanCapacity, closeChan)

	m.life.OnInit()

	// 注册消息处理函数
	m.life.OnRegMsgHandler(m.GetMsgCenter())

	// 获得管理器
	mgrs := m.life.OnSetMgrs()

	// 添加管理器
	m.addMgrs(mgrs)
}

func (m *Module) DoRun() error {
	// 运行模块内消息中心
	m.msgCenter.Run()

	// 运行管理器
	m.runMgrs()

	return m.life.OnRun()
}

func (m *Module) DoStop() error {
	err := m.life.OnStop()

	m.msgCenter.Close()

	// 停止管理器
	m.stopMgrs()

	return err
}

func (m *Module) GetRPC() iface.IRPC {
	return m.rpcGetter.GetRPC()
}

func (m *Module) GetETCD() iface.IETCD { return m.etcdGetter.GetETCD() }

// addMgrs 添加管理器
func (m *Module) addMgrs(mgrs []iface.IMgr) {
	for _, mgr := range mgrs {
		m.mgrs = append(m.mgrs, mgr)
	}

	for _, mgr := range m.mgrs {
		mgr.OnInit()
	}
}

// 运行管理器
func (m *Module) runMgrs() {
	// 正序运行
	for _, mgr := range m.mgrs {
		mgr.OnRun()
	}
}

// 停止管理器
func (m *Module) stopMgrs() {
	// 反序关闭
	mgrsLen := len(m.mgrs)
	for i := mgrsLen - 1; i >= 0; i-- {
		mgr := m.mgrs[i]

		mgr.OnStop()
	}
}

func AddMsgCommHandler[Q def.ICommReq, S interface{}](msgCenter iface.IMsgCenter, commId uint32, handler func(sync bool, req Q) (S, error)) {
	msgCenter.AddHandler(commId, func(sync bool, mr def.ICommReq) (interface{}, error) {
		s, err := handler(sync, mr.(Q))
		if err != nil {
			return nil, err
		}

		return s, nil
	})
}
