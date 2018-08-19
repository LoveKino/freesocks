package freesocks

import (
	"testing"
)

func TestServer(t *testing.T) {
	ListenTcpServer("127.0.0.1", 1234)
}
