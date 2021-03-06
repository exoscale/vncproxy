package main

import (
	"flag"
	"os"

	"github.com/exoscale/vncproxy/logger"
	"github.com/exoscale/vncproxy/proxy"
)

func main() {
	//create default session if required
	var tcpPort = flag.String("tcpPort", "", "tcp port")
	var wsPort = flag.String("wsPort", "", "websocket port")
	var vncPass = flag.String("vncPass", "", "password on incoming vnc connections to the proxy, defaults to no password")
	var recordDir = flag.String("recDir", "", "path to save FBS recordings WILL NOT RECORD if not defined.")
	var targetVnc = flag.String("target", "", "target vnc server (host:port or /path/to/unix.socket)")
	var targetVncPort = flag.String("targPort", "", "target vnc server port (deprecated, use -target)")
	var targetVncHost = flag.String("targHost", "", "target vnc server host (deprecated, use -target)")
	var targetVncPass = flag.String("targPass", "", "target vnc password")
	var dynamicLookup = flag.Bool("dynamicLookup", false, "lookup target UNIX socket path based on WebSocket URI")
	var logLevel = flag.String("logLevel", "info", "change logging level")

	flag.Parse()
	logger.SetLogLevel(*logLevel)

	if *tcpPort == "" && *wsPort == "" {
		logger.Error("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if *targetVnc == "" && *targetVncPort == "" {
		logger.Error("no target vnc server host/port or socket defined")
		flag.Usage()
		os.Exit(1)
	}

	if *dynamicLookup && *wsPort == "" {
		logger.Error("dynamic lookup requires specifying -wsPort")
		os.Exit(1)
	}

	if *vncPass == "" {
		logger.Warn("proxy will have no password")
	}

	tcpUrl := ""
	if *tcpPort != "" {
		tcpUrl = ":" + string(*tcpPort)
	}

	vncProxy := &proxy.VncProxy{
		WsListeningUrl:   "http://0.0.0.0:" + string(*wsPort) + "/", // empty = not listening on ws
		TcpListeningUrl:  tcpUrl,
		ProxyVncPassword: *vncPass, //empty = no auth
		SingleSession: &proxy.VncSession{
			Target:         *targetVnc,
			TargetHostname: *targetVncHost,
			TargetPort:     *targetVncPort,
			TargetPassword: *targetVncPass, //"vncPass",
			ID:             "",
			Status:         proxy.SessionStatusInit,
			Type:           proxy.SessionTypeProxyPass,
		}, // to be used when not using sessions
		DynamicLookup: *dynamicLookup,
		UsingSessions: false, //false = single session - defined in the var above
	}

	if *recordDir != "" {
		logger.Warn("FBS recording is turned on")
		vncProxy.RecordingDir = *recordDir
		vncProxy.SingleSession.Type = proxy.SessionTypeRecordingProxy
	}

	vncProxy.StartListening()
}
