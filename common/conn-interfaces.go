package common

import "io"

type IServerConn interface {
	io.ReadWriter
	Protocol() string
	CurrentPixelFormat() *PixelFormat
	SetPixelFormat(*PixelFormat) error
	SetColorMap(*ColorMap)
	Encodings() []IEncoding
	SetEncodings([]EncodingType) error
	Width() uint16
	Height() uint16
	SetWidth(uint16)
	SetHeight(uint16)
	DesktopName() string
	SetDesktopName(string)
	SetProtoVersion(string)
}

type IClientConn interface {
	CurrentPixelFormat() *PixelFormat
	Encodings() []IEncoding
}
