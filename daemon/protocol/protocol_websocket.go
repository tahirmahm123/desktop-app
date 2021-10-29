//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for IVPN Client Desktop.
//
//  The Daemon for IVPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for IVPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for IVPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package protocol

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func (p *Protocol) wsStart() error {
	log.Info("WebSocket server started")
	defer log.Info("WebSocket server stopped")

	ln, err := net.Listen("tcp", ":7294")
	if err != nil {
		return fmt.Errorf("failed to listen socket: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("PANIC during WebSocket communication!: ", r)
			if err, ok := r.(error); ok {
				log.ErrorTrace(err)
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		_, err = ws.Upgrade(conn)
		if err != nil {
			return err
		}

		{
			// save connection info
			oldconn := p._webSocketUnsafePluginConnection
			p._webSocketUnsafePluginConnection = conn
			if oldconn != nil {
				oldconn.Close()
			}
			log.Info("(WebSocket) New connection")
		}

		for {
			reader := wsutil.NewReader(conn, ws.StateServerSide)
			_, err := reader.NextFrame()
			if err != nil {
				log.Warning("(WebSocket) Connection error: ", err)
				break
			}

			data, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Warning("(WebSocket) Connection error: ", err)
				break
			}

			log.Debug("(WebSocket) data received: ", string(data))

			p.processRequest(conn, string(data))
		}

		{
			// reset connection info
			oldconn := p._webSocketUnsafePluginConnection
			p._webSocketUnsafePluginConnection = nil
			if oldconn != nil {
				oldconn.Close()
			}
		}
	}
}
