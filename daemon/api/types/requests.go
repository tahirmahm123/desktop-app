//
//  Daemon for VPN Client Desktop
//  https://github.com/tahirmahm123/vpn-desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for VPN Desktop.
//
//  The Daemon for VPN Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for VPN Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for VPN Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package types

// SessionNewRequest request to create new session
type SessionNewRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SessionDeleteRequest request to delete session
type SessionDeleteRequest struct {
	Session string `json:"session_token"`
}

// SessionStatusRequest request to get session status
type SessionStatusRequest struct {
	Session string `json:"session_token"`
}

// SessionStatusRequest request to get session status
type ServersListRequest struct {
	Session string `json:"session_token"`
}

// SessionWireGuardKeySetRequest request to set new WK key for a session
type SessionWireGuardKeySetRequest struct {
	Session            string `json:"session_token"`
	PublicKey          string `json:"public_key"`
	ConnectedPublicKey string `json:"connected_public_key"`
}
