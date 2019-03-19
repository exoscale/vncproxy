package server

import (
	"bytes"
	"crypto/des"
	"crypto/rand"
	"errors"
	"log"

	"github.com/exoscale/vncproxy/common"
)

type SecurityType uint8

const (
	SecTypeUnknown  = SecurityType(0)
	SecTypeNone     = SecurityType(1)
	SecTypeVNC      = SecurityType(2)
	SecTypeVeNCrypt = SecurityType(19)
)

type SecuritySubType uint32

const (
	SecSubTypeUnknown = SecuritySubType(0)
)

const (
	SecSubTypeVeNCrypt01Unknown   = SecuritySubType(0)
	SecSubTypeVeNCrypt01Plain     = SecuritySubType(19)
	SecSubTypeVeNCrypt01TLSNone   = SecuritySubType(20)
	SecSubTypeVeNCrypt01TLSVNC    = SecuritySubType(21)
	SecSubTypeVeNCrypt01TLSPlain  = SecuritySubType(22)
	SecSubTypeVeNCrypt01X509None  = SecuritySubType(23)
	SecSubTypeVeNCrypt01X509VNC   = SecuritySubType(24)
	SecSubTypeVeNCrypt01X509Plain = SecuritySubType(25)
)

const (
	SecSubTypeVeNCrypt02Unknown   = SecuritySubType(0)
	SecSubTypeVeNCrypt02Plain     = SecuritySubType(256)
	SecSubTypeVeNCrypt02TLSNone   = SecuritySubType(257)
	SecSubTypeVeNCrypt02TLSVNC    = SecuritySubType(258)
	SecSubTypeVeNCrypt02TLSPlain  = SecuritySubType(259)
	SecSubTypeVeNCrypt02X509None  = SecuritySubType(260)
	SecSubTypeVeNCrypt02X509VNC   = SecuritySubType(261)
	SecSubTypeVeNCrypt02X509Plain = SecuritySubType(262)
)

type SecurityHandler interface {
	Type() SecurityType
	SubType() SecuritySubType
	Auth(common.IServerConn) error
}

// ServerAuthNone is the "none" authentication. See 7.2.1.
type ServerAuthNone struct{}

func (*ServerAuthNone) Type() SecurityType {
	return SecTypeNone
}

func (*ServerAuthNone) Auth(c common.IServerConn) error {
	return nil
}

func (*ServerAuthNone) SubType() SecuritySubType {
	return SecSubTypeUnknown
}

// ServerAuthVNC is the standard password authentication. See 7.2.2.
type ServerAuthVNC struct {
	Pass string
}

func (*ServerAuthVNC) Type() SecurityType {
	return SecTypeVNC
}

func (*ServerAuthVNC) SubType() SecuritySubType {
	return SecSubTypeUnknown
}

const AUTH_FAIL = "Authentication Failure"

func (auth *ServerAuthVNC) Auth(c common.IServerConn) error {
	buf := make([]byte, 8+len([]byte(AUTH_FAIL)))
	rand.Read(buf[:16]) // Random 16 bytes in buf
	sndsz, err := c.Write(buf[:16])
	if err != nil {
		log.Printf("Error sending challenge to client: %s\n", err.Error())
		return errors.New("Error sending challenge to client:" + err.Error())
	}
	if sndsz != 16 {
		log.Printf("The full 16 byte challenge was not sent!\n")
		return errors.New("The full 16 byte challenge was not sent")
	}
	buf2 := make([]byte, 16)
	_, err = c.Read(buf2)
	if err != nil {
		log.Printf("The authentication result was not read: %s\n", err.Error())
		return errors.New("The authentication result was not read" + err.Error())
	}
	AuthText := auth.Pass
	bk, err := des.NewCipher([]byte(fixDesKey(AuthText)))
	if err != nil {
		log.Printf("Error generating authentication cipher: %s\n", err.Error())
		return errors.New("Error generating authentication cipher")
	}
	buf3 := make([]byte, 16)
	bk.Encrypt(buf3, buf)               //Encrypt first 8 bytes
	bk.Encrypt(buf3[8:], buf[8:])       // Encrypt second 8 bytes
	if bytes.Compare(buf2, buf3) != 0 { // If the result does not decrypt correctly to what we sent then a problem
		SetUint32(buf, 0, 1)
		SetUint32(buf, 4, uint32(len([]byte(AUTH_FAIL))))
		copy(buf[8:], []byte(AUTH_FAIL))
		c.Write(buf)
		return errors.New("Authentication failed")
	}
	return nil
}

// SetUint32 set 4 bytes at pos in buf to the val (in big endian format)
// A test is done to ensure there are 4 bytes available at pos in the buffer
func SetUint32(buf []byte, pos int, val uint32) {
	if pos+4 > len(buf) {
		return
	}
	for i := 0; i < 4; i++ {
		buf[3-i+pos] = byte(val)
		val >>= 8
	}
}

// fixDesKeyByte is used to mirror a byte's bits
// This is not clearly indicated by the document, but is in actual fact used
func fixDesKeyByte(val byte) byte {
	var newval byte = 0
	for i := 0; i < 8; i++ {
		newval <<= 1
		newval += (val & 1)
		val >>= 1
	}
	return newval
}

// fixDesKey will make sure that exactly 8 bytes is used either by truncating or padding with nulls
// The bytes are then bit mirrored and returned
func fixDesKey(key string) []byte {
	tmp := []byte(key)
	buf := make([]byte, 8)
	if len(tmp) <= 8 {
		copy(buf, tmp)
	} else {
		copy(buf, tmp[:8])
	}
	for i := 0; i < 8; i++ {
		buf[i] = fixDesKeyByte(buf[i])
	}
	return buf
}
