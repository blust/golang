package blog

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

//Session 结构体
type Session struct {
	sessionId    string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (s *Session) Set(key, value interface{}) {
	s.value[key] = value
	globalSession.Update(s.sessionId)
}

func (s *Session) Get(key interface{}) interface{} {
	if v, ok := s.value[key]; ok {
		globalSession.Update(s.sessionId)
		return v
	}
	return nil
}

func (s *Session) Delete(key interface{}) {
	delete(s.value, key)
}

func (s *Session) SessionId() string {
	return s.sessionId
}

type BlogSessionManager struct {
	cookieName  string
	lock        sync.Mutex
	maxLifeTime int64
	sessions    map[string]*list.Element
	list        *list.List
}

func (m *BlogSessionManager) Start(w http.ResponseWriter, r *http.Request) *Session {
	cookie, err := r.Cookie(m.cookieName)

	var session *Session
	if err != nil || cookie.Value == "" {
		sessionId := m.GetSessionId()
		session = m.Init(sessionId)
		cookie := http.Cookie{Name: m.cookieName, Value: url.QueryEscape(sessionId), Path: "/", HttpOnly: true, MaxAge: int(m.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sessionId, _ := url.QueryUnescape(cookie.Value)
		session = m.Read(sessionId)
	}

	return session
}

//session初始化
func (m *BlogSessionManager) Init(sessionId string) *Session {
	m.lock.Lock()
	defer m.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	session := &Session{sessionId: sessionId, timeAccessed: time.Now(), value: v}
	element := m.list.PushBack(session)
	m.sessions[sessionId] = element
	return session
}

//读取session
func (m *BlogSessionManager) Read(sessionId string) *Session {
	if element, ok := m.sessions[sessionId]; ok {
		return element.Value.(*Session)
	} else {
		return m.Init(sessionId)
	}
}

//销毁session
func (m *BlogSessionManager) Destroy(sessionId string) error {
	if element, ok := m.sessions[sessionId]; ok {
		delete(m.sessions, sessionId)
		m.list.Remove(element)
		return nil
	}
	return nil
}

func (m *BlogSessionManager) Update(sessionId string) {
	if element, ok := m.sessions[sessionId]; ok {
		element.Value.(*Session).timeAccessed = time.Now()
		m.list.MoveToFront(element)
	}
}

//GC你懂得
func (m *BlogSessionManager) GC() {
	for {
		element := m.list.Back()

		if element == nil {
			break
		}

		if (element.Value.(*Session).timeAccessed.Unix() + m.maxLifeTime) < time.Now().Unix() {
			m.Destroy(element.Value.(*Session).sessionId)
		} else {
			break
		}
	}

	time.AfterFunc(time.Duration(m.maxLifeTime)*time.Second, func() {
		fmt.Printf("%s", "a")
		m.GC()
	})
}

//生成唯一SessionId
func (m *BlogSessionManager) GetSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewManager(cookieName string, maxLifeTime int64) *BlogSessionManager {
	return &BlogSessionManager{cookieName: cookieName, maxLifeTime: maxLifeTime, list: list.New(), sessions: make(map[string]*list.Element, 0)}
}

var globalSession *BlogSessionManager

func init() {
	globalSession = NewManager("PHPSESSIONID", 1800)
	go globalSession.GC()
}
