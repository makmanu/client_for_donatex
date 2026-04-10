package plugin

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makmanu/client_for_donatex/config"
)

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

	if err := conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)

	}

	var resp struct {
		MessageType string                 `json:"messageType"`
		Data        map[string]interface{} `json:"data"`
	}
	if err := conn.ReadJSON(&resp); err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	switch resp.MessageType {
	case "AuthenticationTokenResponse":
		authToken, ok := "", false
		if resp.Data != nil {
			switch v := resp.Data["authenticationToken"].(type) {
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
		if resp.Data != nil {
			errorID := 0
			if v, ok := resp.Data["errorID"].(float64); ok {
				errorID = int(v)
			}
			message := ""
			if v, ok := resp.Data["message"].(string); ok {
				message = v
			}
			if errorID == 50 {
				return fmt.Errorf("user rejected plugin: %s", message)
			}
			return fmt.Errorf("authentication failed with APIError %d: %s", errorID, message)
		}
		return fmt.Errorf("authentication failed with APIError")

	default:
		return fmt.Errorf("unexpected auth response messageType: %s", resp.MessageType)
	}
}
 func SessionAuthPlugin(conn *websocket.Conn) error {
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
			"pluginName":      pluginName,
			"pluginDeveloper": "Makmanu",
			"authenticationToken": token,
		},
	}
	if err := conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to send session auth request: %w", err)
	}

	var resp struct {
		MessageType string                 `json:"messageType"`
		Data        map[string]interface{} `json:"data"`
	}
	if err := conn.ReadJSON(&resp); err != nil {
		return fmt.Errorf("failed to read session auth response: %w", err)
	}
	
	switch resp.MessageType {
	case "AuthenticationResponse":
		if resp.Data != nil {
			if authenticated, ok := resp.Data["authenticated"]; !ok || !authenticated.(bool) {
				fmt.Print(ok, authenticated, "3\n")
				fmt.Printf("Session authentication failed: %v\n Trying to get new auth token...\n", resp.Data)
				return GetAuthToken(conn)
			}
			return nil
		}
		return fmt.Errorf("session auth response missing data")

	}
	return fmt.Errorf("unexpected session auth response messageType: %s", resp.MessageType)
}