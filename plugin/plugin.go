package plugin

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makmanu/client_for_donatex/config"
	"gopkg.in/yaml.v3"
)

type Hotkey struct {
	ID          int     `yaml:"id"`
	Name        string  `yaml:"name"`
	HotkeyID    string  `yaml:"hotkeyID"`
	Coefficient float64 `yaml:"coefficient"`
}

type HotkeysResponse struct {
	AvailableHotkeys []Hotkey `yaml:"availableHotkeys"`
}

type wsRequest struct {
	msg      map[string]any
	response chan wsResponse
}

type wsResponse struct {
	messageType string
	data        map[string]any
	err         error
}

var (
	wsRequestCh      chan wsRequest
	pendingResponses map[string]chan wsResponse
	pendingMu        sync.Mutex
	workersOnce      sync.Once
)

func startWsWorkers(conn *websocket.Conn) {
	wsRequestCh = make(chan wsRequest, 32)
	pendingResponses = make(map[string]chan wsResponse)
	go wsSender(conn)
	go wsReader(conn)
}

func ensureWsWorkers(conn *websocket.Conn) {
	workersOnce.Do(func() {
		startWsWorkers(conn)
	})
}

func wsSender(conn *websocket.Conn) {
	for req := range wsRequestCh {
		log.Printf("[REQUEST] %v", req.msg)
		if err := conn.WriteJSON(req.msg); err != nil {
			req.response <- wsResponse{err: fmt.Errorf("failed to send websocket request: %w", err)}
			continue
		}
	}
}

func wsReader(conn *websocket.Conn) {
	for {
		var resp struct {
			RequestID   string         `json:"requestID"`
			MessageType string         `json:"messageType"`
			Data        map[string]any `json:"data"`
		}
		if err := conn.ReadJSON(&resp); err != nil {
			log.Printf("[RESPONSE ERROR] %v", err)
			signalPendingWsError(fmt.Errorf("failed to read websocket response: %w", err))
			return
		}

		log.Printf("[RESPONSE] %v", resp)

		pendingMu.Lock()
		ch, ok := pendingResponses[resp.RequestID]
		if ok {
			delete(pendingResponses, resp.RequestID)
		}
		pendingMu.Unlock()

		if ok {
			ch <- wsResponse{messageType: resp.MessageType, data: resp.Data}
		}
	}
}

func signalPendingWsError(err error) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	for requestID, ch := range pendingResponses {
		ch <- wsResponse{err: err}
		close(ch)
		delete(pendingResponses, requestID)
	}
}

func sendWebsocketRequest(conn *websocket.Conn, msg map[string]any) (wsResponse, error) {
	ensureWsWorkers(conn)

	requestID, ok := msg["requestID"].(string)
	if !ok || requestID == "" {
		requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		msg["requestID"] = requestID
	}

	resultCh := make(chan wsResponse, 1)
	pendingMu.Lock()
	pendingResponses[requestID] = resultCh
	pendingMu.Unlock()

	wsRequestCh <- wsRequest{msg: msg, response: resultCh}

	res := <-resultCh
	if res.err != nil {
		return wsResponse{}, res.err
	}
	return res, nil
}

func ConnectPluginWebsocket(cfg *config.Config) (*websocket.Conn, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}

	url := fmt.Sprintf("%s:%d", cfg.VTubeStudio.URL, cfg.VTubeStudio.Port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to VTubeStudio websocket at %s: %w", url, err)
	}
	return conn, nil
}

func GetAuthToken(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("websocket connection must not be nil")
	}

	pluginName := "DonatexPlugin"

	iconBytes, err := os.ReadFile("plugin/base64_plugin_logo.txt")
	if err != nil {
		return fmt.Errorf("failed to read plugin icon: %w", err)
	}
	pluginIcon := strings.TrimSpace(string(iconBytes))
	if pluginIcon == "" {
		return fmt.Errorf("plugin icon file is empty")
	}

	req := map[string]any{
		"apiName":     "VTubeStudioPublicAPI",
		"apiVersion":  "1.0",
		"requestID":   fmt.Sprintf("auth-%d", time.Now().UnixNano()),
		"messageType": "AuthenticationTokenRequest",
		"data": map[string]any{
			"pluginName":      pluginName,
			"pluginDeveloper": "Makmanu",
			"pluginIcon":      pluginIcon,
		},
	}

	res, err := sendWebsocketRequest(conn, req)
	if err != nil {
		return err
	}

	switch res.messageType {
	case "AuthenticationTokenResponse":
		authToken, ok := "", false
		if res.data != nil {
			switch v := res.data["authenticationToken"].(type) {
			case string:
				authToken, ok = v, true
			}
		}
		if !ok || authToken == "" {
			return fmt.Errorf("authentication response missing authenticationToken")
		}
		if err := os.WriteFile("plugin/token", []byte(authToken), 0o600); err != nil {
			return fmt.Errorf("failed to save auth token: %w", err)
		}
		return SessionAuthPlugin(conn)

	case "APIError":
		if res.data != nil {
			errorID := 0
			if v, ok := res.data["errorID"].(float64); ok {
				errorID = int(v)
			}
			message := ""
			if v, ok := res.data["message"].(string); ok {
				message = v
			}
			if errorID == 50 {
				return fmt.Errorf("user rejected plugin: %s", message)
			}
			return fmt.Errorf("authentication failed with APIError %d: %s", errorID, message)
		}
		return fmt.Errorf("authentication failed with APIError")

	default:
		return fmt.Errorf("unexpected auth response messageType: %s", res.messageType)
	}
}
func SessionAuthPlugin(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("websocket connection must not be nil")
	}

	tokenBytes, err := os.ReadFile("plugin/token")
	if err != nil {
		return fmt.Errorf("failed to read auth token: %w", err)
	}
	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return fmt.Errorf("auth token file is empty")
	}

	pluginName := "DonatexPlugin"

	req := map[string]any{
		"apiName":     "VTubeStudioPublicAPI",
		"apiVersion":  "1.0",
		"requestID":   fmt.Sprintf("session-auth-%d", time.Now().UnixNano()),
		"messageType": "AuthenticationRequest",
		"data": map[string]any{
			"pluginName":          pluginName,
			"pluginDeveloper":     "Makmanu",
			"authenticationToken": token,
		},
	}
	res, err := sendWebsocketRequest(conn, req)
	if err != nil {
		return err
	}

	switch res.messageType {
	case "AuthenticationResponse":
		if res.data != nil {
			if authenticated, ok := res.data["authenticated"]; !ok || !authenticated.(bool) {
				fmt.Print(ok, authenticated, "3\n")
				fmt.Printf("Session authentication failed: %v\n Trying to get new auth token...\n", res.data)
				return GetAuthToken(conn)
			}
			return nil
		}
		return fmt.Errorf("session auth response missing data")

	}
	return fmt.Errorf("unexpected session auth response messageType: %s", res.messageType)
}

func GetCurrentHotkeys(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("websocket connection must not be nil")
	}

	req := map[string]any{
		"apiName":     "VTubeStudioPublicAPI",
		"apiVersion":  "1.0",
		"requestID":   fmt.Sprintf("hotkeys-%d", time.Now().UnixNano()),
		"messageType": "HotkeysInCurrentModelRequest",
	}

	res, err := sendWebsocketRequest(conn, req)
	if err != nil {
		return err
	}

	if res.messageType != "HotkeysInCurrentModelResponse" {
		return fmt.Errorf("unexpected response type: %s", res.messageType)
	}

	resp := struct {
		Data map[string]any `json:"data"`
	}{
		Data: res.data,
	}

	// Read coefficients from file
	coeffData, err := os.ReadFile("plugin/coefficient.yaml")
	if err != nil {
		coeffData = []byte{}
	}

	var coefficients map[any]any
	if len(coeffData) > 0 {
		if err := yaml.Unmarshal(coeffData, &coefficients); err != nil {
			coefficients = make(map[any]any)
		}
	} else {
		coefficients = make(map[any]any)
	}

	// Add ID and coefficient to each hotkey
	if hotkeys, ok := resp.Data["availableHotkeys"].([]any); ok {
		for i, hotkey := range hotkeys {
			if hotkeyMap, ok := hotkey.(map[string]any); ok {
				id := i + 1
				hotkeyMap["id"] = id

				// Look up coefficient or use default
				coeff := 2.5
				if val, ok := coefficients["coefficients"].(map[any]any)[id]; ok {
					if f, ok := val.(float64); ok {
						coeff = f
					}
				}
				hotkeyMap["coefficient"] = coeff
			}
		}
	}

	yamlData, err := yaml.Marshal(resp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal response to yaml: %w", err)
	}

	if err := os.WriteFile("plugin/hotkeys.yaml", yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write hotkeys to file: %w", err)
	}

	return nil
}

func ExecuteHotkey(conn *websocket.Conn, identifier string) error {
	if conn == nil {
		return fmt.Errorf("websocket connection must not be nil")
	}

	// Read hotkeys.yaml
	yamlData, err := os.ReadFile("plugin/hotkeys.yaml")
	if err != nil {
		return fmt.Errorf("failed to read hotkeys.yaml: %w", err)
	}

	var hotkeysResp HotkeysResponse
	if err := yaml.Unmarshal(yamlData, &hotkeysResp); err != nil {
		return fmt.Errorf("failed to unmarshal hotkeys.yaml: %w", err)
	}

	// Find the hotkey by id or name
	var targetHotkeyID string
	found := false
	for _, hotkey := range hotkeysResp.AvailableHotkeys {
		if strconv.Itoa(hotkey.ID) == identifier || hotkey.Name == identifier {
			targetHotkeyID = hotkey.HotkeyID
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("hotkey with identifier '%s' not found", identifier)
	}

	req := map[string]any{
		"apiName":     "VTubeStudioPublicAPI",
		"apiVersion":  "1.0",
		"requestID":   fmt.Sprintf("execute-hotkey-%d", time.Now().UnixNano()),
		"messageType": "HotkeyTriggerRequest",
		"data": map[string]any{
			"hotkeyID": targetHotkeyID,
		},
	}

	fmt.Printf("Sending hotkey trigger request: %s, %s, %s, %s, %v\n", req["apiName"], req["apiVersion"], req["requestID"], req["messageType"], req["data"])

	res, err := sendWebsocketRequest(conn, req)
	if err != nil {
		return err
	}

	if res.messageType != "HotkeyTriggerResponse" {
		return fmt.Errorf("unexpected response type: %s", res.messageType)
	}

	return nil
}
