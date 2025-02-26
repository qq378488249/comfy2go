package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	Ws    = "ws"
	Wss   = "wss"
	Http  = "http"
	Https = "https"
)

// NewComfyClient creates a new instance of a Comfy2go client
func NewComfyClientWithProtocolType(server_address string, server_port int, protocolType string, callbacks *ComfyClientCallbacks) *ComfyClient {
	cid := strings.ReplaceAll(uuid.New().String(), "-", "")
	if protocolType != Https {
		// default http
		protocolType = Http
	}
	sbaseaddr := fmt.Sprintf("%s://%s:%d", protocolType, server_address, server_port)

	wsType := Ws
	if protocolType == Https {
		wsType = Wss
	}
	webSocketURL := fmt.Sprintf("%s://%s:%d/ws?clientId=%s", wsType, server_address, server_port, cid)

	retv := &ComfyClient{
		serverBaseAddress: sbaseaddr,
		serverAddress:     server_address,
		serverPort:        server_port,
		clientid:          cid,
		queueditems:       make(map[string]*QueueItem),
		webSocket: &WebSocketConnection{
			WebSocketURL:   webSocketURL,
			ConnectionDone: make(chan bool),
			MaxRetry:       5, // Maximum number of retries
			ManagerStarted: false,
			BaseDelay:      1 * time.Second,
			MaxDelay:       10 * time.Second,
		},
		initialized: false,
		queuecount:  0,
		callbacks:   callbacks,
		timeout:     -1,
		httpclient:  &http.Client{},
	}
	// golang uses mark-sweep GC, so this circular reference should be fine
	retv.webSocket.Callback = retv
	return retv
}
