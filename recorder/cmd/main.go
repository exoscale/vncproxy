package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/exoscale/vncproxy/client"
	"github.com/exoscale/vncproxy/common"
	"github.com/exoscale/vncproxy/encodings"
	"github.com/exoscale/vncproxy/logger"
	"github.com/exoscale/vncproxy/recorder"
)

func main() {
	var recordDir = flag.String("recFile", "", "FBS file to create, recordings WILL NOT RECORD IF EMPTY.")
	var targetVncPort = flag.String("targPort", "", "target vnc server port")
	var targetVncPass = flag.String("targPass", "", "target vnc password")
	var targetVncHost = flag.String("targHost", "localhost", "target vnc hostname")
	var logLevel = flag.String("logLevel", "info", "change logging level")

	flag.Parse()
	logger.SetLogLevel(*logLevel)

	if *targetVncHost == "" {
		logger.Error("no target vnc server host defined")
		flag.Usage()
		os.Exit(1)
	}

	if *targetVncPort == "" {
		logger.Error("no target vnc server port defined")
		flag.Usage()
		os.Exit(1)
	}

	if *targetVncPass == "" {
		logger.Warn("no password defined, trying to connect with null authentication")
	}
	if *recordDir == "" {
		logger.Warn("FBS recording is turned off")
	} else {
		logger.Infof("Recording rfb stream into file: '%s'", *recordDir)
	}

	nc, err := net.Dial("tcp", *targetVncHost+":"+*targetVncPort)

	if err != nil {
		logger.Errorf("error connecting to vnc server: %s", err)
	}
	var noauth client.ClientAuthNone
	authArr := []client.ClientAuth{&client.PasswordAuth{Password: *targetVncPass}, &noauth}

	rec, err := recorder.NewRecorder(*recordDir) //"/Users/amitbet/vncRec/recording.rbs")
	if err != nil {
		logger.Errorf("error creating recorder: %s", err)
		return
	}

	clientConn, err := client.NewClientConn(nc,
		&client.ClientConfig{
			Auth:      authArr,
			Exclusive: true,
		})

	clientConn.Listeners.AddListener(rec)
	clientConn.Listeners.AddListener(&recorder.RfbRequester{Conn: clientConn, Name: "Rfb Requester"})
	clientConn.Connect()

	if err != nil {
		logger.Errorf("error creating client: %s", err)
		return
	}
	encs := []common.IEncoding{
		&encodings.TightEncoding{},
		&encodings.PseudoEncoding{int32(common.EncJPEGQualityLevelPseudo8)},
	}

	clientConn.SetEncodings(encs)

	for {
		time.Sleep(time.Minute)
	}
}

func getNowMillisec() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
