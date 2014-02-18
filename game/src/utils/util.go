package utils

import (
	"fmt"
	"runtime"
	"stat"
	"strconv"
	"net"
	" "encoding/binary""
)


type Stat struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Cpunum  int
	Memused float64
	Clients int
}

func SysStat() *Stat {
	load := stat.GetLoadAvgSample()
	mem := stat.GetMemSample()

	memused := float64(mem.MemTotal - mem.MemFree - mem.Cached - mem.Buffers)
	memtotal := float64(mem.MemTotal)
	usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memused/memtotal), 64)

	return &Stat{
		Load1:   load.One,
		Load5:   load.Five,
		Load15:  load.Fifteen,
		Cpunum:  runtime.NumCPU(),
		Memused: usedPercent,
	}
}

func Ip2Uint32(ipstr string) uint32 {
        ip := net.ParseIP(ipstr)
        if ip == nil {
            return 0
        }
        ip = ip.To4()
        return binary.BigEndian.Uint32(ip)
}

func Ip2String(ipLong uint32) string {
        ipByte := make([]byte, 4)
        binary.BigEndian.PutUint32(ipByte, ipLong)
        ip := net.IP(ipByte)
        return ip.String()
}

func PrintPanicStack() {
	if x := recover(); x != nil {
		log.Printf("%v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}