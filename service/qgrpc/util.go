package qgrpc

import (
	"encoding/json"
	"github.com/qsock/qf/net/ipaddr"
	"github.com/qsock/qf/service/qetcd"
	"path"
	"strings"
	"time"
)

func GetEtcdClient() *qetcd.Op {
	return O.etcdClient
}

func GetPrefixAddrsRegisterModel(serviceName string) []*RegisterModel {
	prefixKey := GetListenerKey(serviceName)
	kvs := GetEtcdClient().GetWatch(prefixKey)
	ms := make([]*RegisterModel, 0, len(kvs))
	for k, v := range kvs {
		m := new(RegisterModel)
		if err := json.Unmarshal([]byte(v), m); err != nil {
			continue
		}
		m.Name = k
		ms = append(ms, m)
	}
	return ms
}

func GetRegisterServerNames() []string {
	c := O.C
	arr := strings.Split(c.ServerName, "|")
	return arr
}

func GetRegisterKey(key string) string {
	c := O.C
	k := path.Join("/", c.Schema, c.Prefix, key, ipaddr.GetHostname())
	return k
}

func GetRegisterValue() string {
	c := O.C
	m := new(RegisterModel)
	m.Name = c.ServerName
	m.Addr = c.Addr
	m.Ext = c.Ext
	m.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	m.Hostname = ipaddr.GetHostname()
	v, _ := json.Marshal(m)
	return string(v)
}

func GetListenerKey(serverName string) string {
	c := O.C
	k := path.Join("/", c.Schema, c.Prefix, serverName)
	return k
}

func ExtractRegisterModel(val string) *RegisterModel {
	m := new(RegisterModel)
	if err := json.Unmarshal([]byte(val), m); err != nil {
		return nil
	}
	return m
}
