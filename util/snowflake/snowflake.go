package snowflake

import (
	"errors"
	"github.com/qsock/qf/net/ipaddr"
	"net"
	"sync"
	"time"
)

// 公有的const
const (
	// 总共64位的snowflake算法
	// BIT_LEN_TIME 时间的长度
	BIT_LEN_TIME = 33
	// BIT_LEN_SEQUENCE 自增序列长度
	BIT_LEN_SEQUENCE = 10
	// BIT_LEN_MACHINE_ID 机器码长度
	BIT_LEN_MACHINE_ID = 10
)

// 可以设置的私有变量
var (
	// startTime 开始时间33位毫秒数
	// 2020-09-13 20:26:39 -> 1599999999
	startTime = int64(1599999999999)
	// machineID 机器码,默认值为0
	machineID uint16
	// elapsedTime 上一次的时间序列
	elapsedTime int64
	// 自增序列
	sequence uint16

	mutex *sync.Mutex
)

// 私有的const
const (
	snowflakeTimeUnit = 1e6
	maskSequence      = uint16(1<<BIT_LEN_SEQUENCE - 1)
)

func init() {
	SetMachineID(0)
	mutex = new(sync.Mutex)
}

// SetStartTime 设置开始时间,返回设置成功或者失败
func SetStartTime(st time.Time) bool {
	if st.After(time.Now()) {
		return false
	}

	if st.IsZero() {
		// 不变
		return true
	}
	// 这里设置snowflake时间
	startTime = toSnowFlakeTime(st)
	return true
}

// GetStartTime 得到开始时间序列
func GetStartTime() int64 {
	return startTime
}

// toSnowFlakeTime 转换为snowflake时间
func toSnowFlakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / snowflakeTimeUnit
}

// SetMachineID 设置机器码
func SetMachineID(mid uint16) bool {
	// 机器id
	maxMachineID := uint16((1 << BIT_LEN_MACHINE_ID) - 1)
	if maxMachineID < mid {
		return false
	}
	if mid == 0 {
		var err error
		// 取机器ip的低16位
		machineID, err = getLower16BitPrivateIP()
		if err != nil {
			// 这里panic,因为都是初始化操作
			panic(err)
		}
	} else {
		machineID = mid
	}

	return true

}

// SetMachineID 得到机器码
func GetMachineID() uint16 {
	return machineID
}

// 得到机器码 v4版本
func getLower16BitPrivateIP() (uint16, error) {
	localIp := ipaddr.GetLocalIp()
	pIPv4 := net.ParseIP(localIp)

	if pIPv4 == nil {
		return 0, errors.New("no private ip address")
	}

	// 交换192.168.0.123,交换123和0的高低位,保证局域网内唯一
	return uint16(pIPv4[2]>>(BIT_LEN_MACHINE_ID-8) + pIPv4[3]<<(BIT_LEN_MACHINE_ID-8)), nil
}

// 生成id
func NextId() int64 {
	mutex.Lock()
	defer mutex.Unlock()
	// 当前已经过去的时间
	current := currentElapsedTime()

	if elapsedTime < current {
		// 重新开始序列
		elapsedTime = current
		sequence = 0
	} else {
		// 取自增序列
		sequence = (sequence + 1) & maskSequence
		// 达到最大的值了
		if sequence == 0 {
			elapsedTime++
			overTime := elapsedTime - current // 睡眠
			time.Sleep(sleepTime(overTime))
		}
	}
	return toID()
}

func toID() int64 {
	// 如果超过时间序列, 到了2286年？那就只能认命了
	if elapsedTime >= 1<<BIT_LEN_TIME {
		return 0
	}
	seq1 := int64(elapsedTime << (BIT_LEN_SEQUENCE + BIT_LEN_MACHINE_ID))
	seq2 := int64(machineID << BIT_LEN_SEQUENCE)
	seq3 := int64(sequence)
	// 第一部分时间序列,第二部分机器序列,第三部分自增序列
	return seq1 | seq2 | seq3
}

//  提取trace的时间
func ToTimeUnix(id int64) int64 {
	return ((id >> (BIT_LEN_SEQUENCE + BIT_LEN_MACHINE_ID)) + startTime) / 1000
}

// 已经过去的时间
func currentElapsedTime() int64 {
	return toSnowFlakeTime(time.Now().UTC()) - startTime
}

// 获取睡眠时间
func sleepTime(overtime int64) time.Duration {
	// 计算准确的睡眠时间,度过下一次的时间序列
	return time.Duration(overtime)*time.Millisecond - time.Duration(time.Now().UTC().UnixNano()%snowflakeTimeUnit)*time.Nanosecond
}
