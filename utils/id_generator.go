package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	mrand "math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

/*
 uuid utils
 @author Tony Tian
 @date 2018-04-14
 @version 1.0.0
*/

/*
  return eg: 4725f5ae6a350b1c45687c9934456e6f
*/
func GUID() string {
	s := newRFC4122Generator()
	u, _ := s.NewV1()
	buf := make([]byte, 32)
	hex.Encode(buf[0:4], u[6:8])
	hex.Encode(buf[4:8], u[4:6])
	hex.Encode(buf[8:16], u[0:4])
	hex.Encode(buf[16:20], u[8:10])
	hex.Encode(buf[20:], u[10:])
	return string(buf)
}
func UUID() string {
	s := newRFC4122Generator()
	u, _ := s.NewV1()
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])
	return string(buf)
}
func RandomString(length int, bases ...string) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	if bases != nil && len(bases) > 0 {
		str = strings.Join(bases, "")
	}
	bytes := []byte(str)
	result := []byte{}
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const Size = 16

type uuid [Size]byte

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

type rfc4122Generator struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex

	rand io.Reader

	epochFunc     epochFunc
	hwAddrFunc    hwAddrFunc
	lastTime      uint64
	clockSequence uint16
	hardwareAddr  [6]byte
}

const epochStart = 122192928000000000

func newRFC4122Generator() *rfc4122Generator {
	return &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

// NewV1 returns UUID based on current timestamp and MAC address.
func (g *rfc4122Generator) NewV1() (uuid, error) {
	u := uuid{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return uuid{}, err
	}
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	hardwareAddr, err := g.getHardwareAddr()
	if err != nil {
		return uuid{}, err
	}
	copy(u[10:], hardwareAddr)
	return u, nil
}

// Returns epoch and clock sequence.
func (g *rfc4122Generator) getClockSequence() (uint64, uint16, error) {
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = io.ReadFull(g.rand, buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})
	if err != nil {
		return 0, 0, err
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := g.getEpoch()
	// Clock didn't change since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

// Returns hardware address.
func (g *rfc4122Generator) getHardwareAddr() ([]byte, error) {
	var err error
	g.hardwareAddrOnce.Do(func() {
		if hwAddr, err := g.hwAddrFunc(); err == nil {
			copy(g.hardwareAddr[:], hwAddr)
			return
		}

		// Initialize hardwareAddr randomly in case
		// of real network interfaces absence.
		if _, err = io.ReadFull(g.rand, g.hardwareAddr[:]); err != nil {
			return
		}
		// Set multicast bit as recommended by RFC 4122
		g.hardwareAddr[0] |= 0x01
	})
	if err != nil {
		return []byte{}, err
	}
	return g.hardwareAddr[:], nil
}

// Returns difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and current time.
func (g *rfc4122Generator) getEpoch() uint64 {
	return epochStart + uint64(g.epochFunc().UnixNano()/100)
}
func defaultHWAddrFunc() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []byte{}, err
	}
	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}
	return []byte{}, fmt.Errorf("uuid: no HW address found")
}
