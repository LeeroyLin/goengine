package iface

type IApp interface {
	// Init 初始化
	Init(preModules, modules []IModule)

	// Run 运行应用
	//
	// param successCb 启动成功时的回调
	Run(successCb func())

	// Stop 停止应用
	Stop()
}

type IAppHandler interface {
	OnBeforeInit()
	OnAfterInit()
	OnBeforeRun()
	OnAfterRun()
	OnBeforeStop()
	OnAfterStop()
}
