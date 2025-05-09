// Copyright 2025 novnc4svc Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package novnc4svc

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(r *http.Request) bool { return true },
	Subprotocols:      []string{"binary", "chat"},
	EnableCompression: true,
}

func wsVNC(idQueryName string, transformer func(id string) (wsUrl string)) func(ctx *gin.Context) {
	if transformer == nil {
		panic("the ws VNC transformer is nil")
	}
	return func(ctx *gin.Context) {
		clientConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "websocket upgrade error:" + err.Error()})
			return
		}
		defer clientConn.Close()
		targetHeaders := http.Header{"Sec-WebSocket-Protocol": []string{"binary"}}
		url0 := ctx.Request.URL
		id0 := url0.Query().Get(idQueryName)
		wsUrl := transformer(id0)
		targetURL := wsUrl + "?" + url0.RawQuery
		dialer := *websocket.DefaultDialer
		dialer.HandshakeTimeout = 10 * time.Second
		backendConn, _, err := dialer.Dial(targetURL, targetHeaders)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer backendConn.Close()
		errChan := make(chan error, 2)
		go func() {
			for {
				messageType, message, err0 := clientConn.ReadMessage()
				if err0 != nil {
					errChan <- errors.New("")
					return
				}
				if err0 = backendConn.WriteMessage(messageType, message); err0 != nil {
					errChan <- errors.New("")
					return
				}
			}
		}()
		go func() {
			for {
				messageType, message, err0 := backendConn.ReadMessage()
				if err0 != nil {
					errChan <- errors.New("")
					return
				}
				if err0 = clientConn.WriteMessage(messageType, message); err0 != nil {
					errChan <- errors.New("")
					return
				}
			}
		}()
		<-errChan
	}
}
