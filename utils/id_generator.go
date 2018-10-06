package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type idGenerator struct {
	mu       sync.Mutex
	count    uint32
	hardware []byte
	pid      int
}

var idGen idGenerator

func init() {
	idGen = idGenerator{
		mu:       sync.Mutex{},
		count:    0,
		hardware: make([]byte, 6),
		pid:      os.Getpid(),
	}
}

func InitGenerator(interfaceName string) {
	netInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		panic(err)
	}
	copy(idGen.hardware, []byte(netInterface.HardwareAddr))
}

// 版本:hardware:pid:nanoSecond:count
func Generate() (string, error) {
	version := make([]byte, 2)
	binary.LittleEndian.PutUint16(version, uint16(1))

	pid := make([]byte, 4)
	binary.LittleEndian.PutUint32(pid, uint32(idGen.pid))

	nanoSecond := make([]byte, 8)
	binary.LittleEndian.PutUint64(nanoSecond, uint64(time.Now().UnixNano()))

	var countValue uint32
	idGen.mu.Lock()
	countValue = idGen.count
	idGen.count += 1
	idGen.mu.Unlock()
	count := make([]byte, 4)
	binary.LittleEndian.PutUint32(count, countValue)

	buffer := bytes.NewBuffer(version)
	buffer.Write(idGen.hardware[:6])
	buffer.Write(pid[:4])
	buffer.Write(nanoSecond[:8])
	buffer.Write(count[:4])

	str := base64.URLEncoding.EncodeToString(buffer.Bytes())
	return str, nil
}

func GenID() (string, error) {
	return Generate()
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// TODO: 做成单独的服务