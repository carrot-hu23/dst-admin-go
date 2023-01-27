package session

//IProvider Session管理接口
//提供 Session 存储，Session存储方式接口
type IProvider interface {
	//初始化：Session初始化以获取Session
	Init(sid string) (ISession, error)
	//读取：根据SessionID获取Session内容
	Read(sid string) (ISession, error)
	//销毁：根据SessionID删除Session内容
	Destroy(sid string) error
	//回收：根据过期时间删除Session
	GC(maxAge int64)
}

//providers Provider管理器集合
var providers = make(map[string]IProvider)

//Register 根绝Provider管理器名称获取Provider管理器
func Register(name string, provider IProvider) {
	if provider == nil {
		panic("provider register: provider is nil")
	}
	if _, ok := providers[name]; ok {
		panic("provider register: provider already exists")
	}
	providers[name] = provider
}
