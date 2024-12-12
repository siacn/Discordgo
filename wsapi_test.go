package discordgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebsocket_BasicOpenClose(t *testing.T) {
	const discordToken = "test-token"

	// Run a websocket server that acts as Discord.
	server := httptest.NewServer(fakeDiscordWebsocketHandler(t, discordToken))
	defer server.Close()

	session, err := New(discordToken)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Intercept outgoing requests and redirect them to the fake websocket server.
	session.Client.Transport = fakeDiscordGatewayTransport(t, server.URL)

	// Test opening the websocket connection.
	if err := session.Open(); err != nil {
		t.Fatalf("failed to open session: %v", err)
	}

	// Test closing the websocket connection.
	if err := session.Close(); err != nil {
		t.Fatalf("failed to close session: %v", err)
	}
}

func fakeDiscordWebsocketHandler(t *testing.T, expectToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader websocket.Upgrader

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("failed to upgrade websocket: %v", err)
		}

		defer c.Close()

		if err := c.WriteJSON(&Event{Operation: 10, RawData: json.RawMessage(`{"heartbeat_interval": 45000}`)}); err != nil {
			t.Fatalf("failed to write hello: %v", err)
		}

		var event struct {
			Operation int `json:"op"`
			Data      struct {
				Token string `json:"token"`
			} `json:"d"`
		}

		if err := c.ReadJSON(&event); err != nil {
			t.Fatalf("failed to read identity event: %v", err)
		}

		if event.Operation != 2 {
			t.Fatalf("wrong identity op code (got %d; want %d)", event.Operation, 2)
		}

		if event.Data.Token != expectToken {
			t.Errorf("wrong token (got %s; want %s)", event.Data.Token, expectToken)
		}

		if err := c.WriteJSON(&Event{Type: "READY"}); err != nil {
			t.Fatalf("failed to write ready event: %v", err)
		}

		// Discard all messages until a close is received.
		for {
			if _, _, err := c.NextReader(); err != nil {
				var closeErr *websocket.CloseError
				if !errors.As(err, &closeErr) {
					t.Errorf("websocket did not get a close frame (got %v)", err)
				}
				_ = c.Close()
			}
		}
	}
}

func fakeDiscordGatewayTransport(t *testing.T, serverURL string) roundTripperFunc {
	return func(req *http.Request) (*http.Response, error) {
		expectedPath := fmt.Sprintf("/api/v%s/gateway", APIVersion)
		if req.URL.Path != expectedPath {
			t.Fatalf("wrong gateway request path (got %s; want %s)", req.URL.Path, expectedPath)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("wrong gateway request method (got %s; want %s)", req.Method, http.MethodGet)
		}

		url := "ws" + strings.TrimPrefix(serverURL, "http")
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"url":%q}`, url))),
		}
		return resp, nil
	}
}
