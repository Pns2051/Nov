package native

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/Pns2051/adblock-system/proxy/internal/config"
	"github.com/Pns2051/adblock-system/proxy/internal/proxy"
)

type NativeMessage struct {
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type NativeResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func RunNativeHost(adBlockerProxy *proxy.AdBlockerProxy) {
	for {
		var length uint32
		err := binary.Read(os.Stdin, binary.LittleEndian, &length)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading length: %v", err)
		}

		msgBytes := make([]byte, length)
		_, err = io.ReadFull(os.Stdin, msgBytes)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		var msg NativeMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		var resp NativeResponse
		switch msg.Command {
		case "ping":
			resp = NativeResponse{Status: "ok", Message: "pong"}
		case "setEnabled":
			if val, ok := msg.Payload["value"].(bool); ok {
				adBlockerProxy.SetEnabled(val)
				resp = NativeResponse{Status: "ok"}
			} else {
				resp = NativeResponse{Status: "error", Message: "invalid payload"}
			}
		case "updateBlocklist":
			go adBlockerProxy.UpdateBlocklist([]string{config.PrimaryBlocklistURL, config.FallbackBlocklistURL})
			resp = NativeResponse{Status: "ok", Message: "update started"}
		case "getStatus":
			resp = NativeResponse{
				Status: "ok",
				Data: map[string]interface{}{
					"enabled":       adBlockerProxy.Enabled(),
					"blocklistSize": adBlockerProxy.Blocklist.Size(),
				},
			}
		default:
			resp = NativeResponse{Status: "error", Message: "unknown command"}
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}

		binary.Write(os.Stdout, binary.LittleEndian, uint32(len(respBytes)))
		os.Stdout.Write(respBytes)
	}
}
