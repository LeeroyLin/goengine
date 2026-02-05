package iface

type IModuleLife interface {
	OnInit()                              // 初始化
	OnRun() error                         // 运行
	OnBeforeStop() error                  // 停止前处理
	OnStop() error                        // 停止
	OnSetMgrs() []IMgr                    // 获得模块的管理器
	OnRegMsgHandler(msgCenter IMsgCenter) // 注册消息中心处理函数
}
