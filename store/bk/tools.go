package bk

import (
	"github.com/kr/beanstalk"
	"strconv"
	"time"
)

//获取队列ready的消息个数
func GetQueueLen(addr, tube string) (int, error) {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	t := &beanstalk.Tube{Conn: conn, Name: tube}
	ret, err := t.Stats()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(ret["current-jobs-ready"])
}

func MultiGetQueueLen(addr string, tubes []string) (map[string]int, error) {
	ret := map[string]int{}
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return ret, err
	}

	defer conn.Close()
	t := &beanstalk.Tube{Conn: conn}
	for _, tube := range tubes {
		t.Name = tube

		m, err := t.Stats()
		if err != nil {
			continue
		}
		num, _ := strconv.Atoi(m["current-jobs-ready"])
		ret[tube] = num
	}

	return ret, nil
}

func ReserveOne(addr, tube string) (uint64, []byte, error) {
	return MultiReserveOne(addr, []string{tube})
}

func MultiReserveOne(addr string, tubes []string) (id uint64, body []byte, err error) {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return
	}

	connTube := beanstalk.NewTubeSet(conn, tubes...)

	defer connTube.Conn.Close()

	id, body, err = connTube.Reserve(time.Second)
	if err != nil {
		return
	}
	if err = connTube.Conn.Delete(id); err != nil {
		return
	}
	return
}
