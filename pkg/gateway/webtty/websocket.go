/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package webtty

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const bufferSize = 1024

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebsocket(w http.ResponseWriter, req *http.Request) (*WSConnection, error) {
	wsSocket, err := wsUpgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, err
	}

	wsConn := &WSConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WSMessage, bufferSize),
		outChan:   make(chan *WSMessage, bufferSize),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	go wsConn.ReadLoop()
	go wsConn.WriteLoop()

	return wsConn, nil
}

type WSMessage struct {
	MessageType int
	Data        []byte
}

type WSConnection struct {
	wsSocket  *websocket.Conn
	inChan    chan *WSMessage
	outChan   chan *WSMessage
	mutex     sync.Mutex
	isClosed  bool
	closeChan chan byte
}

func (c *WSConnection) ReadLoop() {
	for {
		msgType, data, err := c.wsSocket.ReadMessage()
		if err != nil {
			break
		}

		c.inChan <- &WSMessage{
			msgType,
			data,
		}
	}
}

func (c *WSConnection) WriteLoop() {
	for {
		select {
		case msg := <-c.outChan:
			if err := c.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				break
			}
		case <-c.closeChan:
			c.Close()
		}
	}
}

func (c *WSConnection) Close() {
	_ = c.wsSocket.Close()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.isClosed {
		c.isClosed = true
		close(c.closeChan)
	}
}

func (c *WSConnection) Write(messageType int, data []byte) (err error) {
	select {
	case c.outChan <- &WSMessage{messageType, data}:
	case <-c.closeChan:
		err = errors.New("write websocket closed")
		break
	}
	return
}

func (c *WSConnection) Read() (msg *WSMessage, err error) {
	select {
	case msg = <-c.inChan:
		return
	case <-c.closeChan:
		err = errors.New("read websocket closed")
		break
	}
	return
}
