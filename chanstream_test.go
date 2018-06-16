// Copyright 2014 Garrett D'Amore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package chanstream provides an API that is similar to that used for TCP
// and Unix Domain sockets (see net.TCP), for use in intra-process
// communication on top of Go channels.  This makes it easy to swap it for
// another net.Conn interface.
//
// By using channels, we avoid exposing any
// interface to other processors, or involving the kernel to perform data
// copying.
package chanstream

import "testing"
import "bytes"

func TestListenAndAccept(t *testing.T) {
	name := "test1"
	t.Logf("Establishing listener.")
	listener, err := ListenChan(name)
	if err != nil {
		t.Errorf("ListenChan failed: %v", err)
		return
	}

	go func() {
		t.Logf("Connecting")
		client, err := DialChan(name)
		if err != nil {
			t.Errorf("DialChan failed: %v", err)
			return
		}
		t.Logf("Connected client: %s->%s", client.LocalAddr(), client.RemoteAddr())
	}()

	server, err := listener.Accept()
	if err != nil {
		t.Errorf("Accept failed: %v", err)
	}
	t.Logf("Connected server: %s (client %s)", server.LocalAddr(), server.RemoteAddr())
}

func TestDuplicateListen(t *testing.T) {
	name := "test2"
	var err error
	listener, err := ListenChan(name)
	if err != nil {
		t.Errorf("ListenChan failed: %v", err)
		return
	}
	t.Logf("listener: %v", listener)
	_, err = ListenChan(name)
	if err != ErrAddrInUse {
		t.Errorf("Duplicate listen should not be permitted!")
		return
	}
	t.Logf("Got expected error: %v", err)
}

func TestConnRefused(t *testing.T) {
	name := "test3"
	var err error
	conn, err := DialChan(name)
	if err != ErrConnRefused {
		t.Errorf("Connection not refused (%s)!", conn.LocalAddr())
		return
	}
	t.Logf("Got expected error: %v", err)
}

func TestEcho(t *testing.T) {
	name := "test4"

	master := make([]byte, 1024)
	for i := range master {
		master[i] = uint8(i & 0xff)
	}

	t.Logf("Establishing listener.")
	listener, err := ListenChan(name)
	if err != nil {
		t.Errorf("ListenChan failed: %v", err)
		return
	}

	go func() {
		// Client side
		req := make([]byte, len(master))
		copy(req, master)
		t.Logf("Connecting")
		client, err := DialChan(name)
		if err != nil {
			t.Errorf("DialChan failed: %v", err)
			return
		}
		t.Logf("Connected client: %s", client.LocalAddr())

		// Now try to send data
		t.Logf("Sending %d bytes", len(req))

		n, err := client.Write(req)
		t.Logf("Client sent %d bytes, err %v", n, err)

		// Now zero out our path
		for i := range req {
			req[i] = 1
		}

		rep := make([]byte, len(req))
		n, err = client.Read(rep)
		if n != len(rep) {
			t.Errorf("Client receive error: %n, %v", n, err)
			return
		}

		if !bytes.Equal(rep, master) {
			t.Errorf("Reply mismatch: %v, %v", rep, master)
		}
	}()

	server, err := listener.Accept()
	if err != nil {
		t.Errorf("Accept failed: %v", err)
	}
	t.Logf("Connected server: %s (client %s)", server.LocalAddr(), server.RemoteAddr())

	rcv := make([]byte, len(master))
	rep := make([]byte, len(master))
	// Now we can try to send and receive
	n, err := server.Read(rcv)
	if n != len(master) {
		t.Errorf("Server received too few bytes: %n, %v", n, err)
		return
	}
	t.Logf("Server received %d bytes, err %v", n, err)
	if !bytes.Equal(rcv, master) {
		t.Errorf("Request mismatch: %v, %v", rep, master)
		return
	}

	// Now reply
	copy(rep, rcv)
	n, err = server.Write(rep)
	if n != len(rep) {
		t.Errorf("Server sent too few bytes: %n, %v", n, err)
		return
	}
	t.Logf("Server replied with %d bytes", len(rep))

}
