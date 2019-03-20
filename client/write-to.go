package client

import (
	"io"

	"github.com/exoscale/vncproxy/common"
	"github.com/exoscale/vncproxy/logger"
)

type WriteTo struct {
	Writer io.Writer
	Name   string
}

func (p *WriteTo) Consume(seg *common.RfbSegment) error {
	switch seg.SegmentType {
	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		if err != nil {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentBytes): problem writing to port: %s", err)
		}
		return err

	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		err := clientMsg.Write(p.Writer)
		if err != nil {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err

	default:
	}
	return nil
}
