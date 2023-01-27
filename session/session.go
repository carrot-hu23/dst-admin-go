package session

//ISession 操作接口
//Session数据结构为散列表kv
type ISession interface {
	//Set 设置
	Set(key, value interface{}) error
	//Get 获取
	Get(key interface{}) interface{}
	//Delete 删除
	Delete(key interface{}) error
	//SessionID
	SessionID() string
}
