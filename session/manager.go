package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func init() {

}

// Manager 封装Provider
type Manager struct {
	mutex      sync.Mutex //互斥锁
	provider   IProvider  //Session存储方式
	cookieName string     //Cookie名称
	maxAge     int64      //过期时间
}

// NewManager 实例化Session管理器
func NewManager(providerName, cookieName string, maxAge int64) *Manager {
	provider, ok := providers[providerName]
	if !ok {
		log.Printf("session 初始化失败")
		return nil
	}
	return &Manager{provider: provider, cookieName: cookieName, maxAge: maxAge}
}

// SessionID 生成全局唯一Session标识用于识别每个用户
func (m *Manager) SessionID() string {
	buf := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(buf)
}

// Start 根据当前请求中的COOKIE判断是否存在有效的Session，不存在则创建。
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) ISession {
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//获取Cookie
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		//创建SessionID
		sid := m.SessionID()
		//Session初始化
		session, _ := m.provider.Init(sid)
		c1 := http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(m.maxAge),
		}
		c2 := http.Cookie{
			Name:     "JSESSIONID",
			Value:    sid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(m.maxAge),
		}
		w.Header().Add("Set-Cookie", c1.String())
		w.Header().Add("Set-Cookie", c2.String())
		return session
	} else {
		//从Cookie获取SessionID
		sid, _ := url.QueryUnescape(cookie.Value)
		//获取Session
		session, _ := m.provider.Read(sid)
		return session
	}
}

// Destroy 注销Session
func (m *Manager) Destroy(w http.ResponseWriter, r *http.Request) {
	//从请求中读取Cookie值
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//销毁Session内容
	m.provider.Destroy(cookie.Value)
	//设置客户端Cookie立即过期
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now(),
	})
}

// GC 销毁Session
func (m *Manager) GC() {
	//添加互斥锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//设置过期时间销毁Seesion
	m.provider.GC(m.maxAge)
	//添加计时器当Session超时自动销毁
	time.AfterFunc(time.Duration(m.maxAge), func() {
		m.GC()
	})
}
