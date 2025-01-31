package term

import (
	"github.com/gliderlabs/ssh"
	"net"
	"time"
)

// SSHConn is a wrapper around ssh.Session that implements net.Conn
type SSHConn struct {
	ssh.Session
}

// NewSSHConn wraps an ssh.Session into a net.Conn
func NewSSHConn(session ssh.Session) *SSHConn {
	return &SSHConn{Session: session}
}

// Implement the net.Conn methods on SSHConn
func (c *SSHConn) Read(b []byte) (n int, err error) {
	return c.Session.Read(b)
}

func (c *SSHConn) Write(b []byte) (n int, err error) {
	return c.Session.Write(b)
}

func (c *SSHConn) Close() error {
	return c.Session.Close()
}

func (c *SSHConn) LocalAddr() net.Addr {
	return c.Session.LocalAddr()
}

func (c *SSHConn) RemoteAddr() net.Addr {
	return c.Session.RemoteAddr()
}

// Implementing SetDeadline method
func (c *SSHConn) SetDeadline(t time.Time) error {
	// `ssh.Session` does not provide native support for SetDeadline,
	// so we do nothing here or you can handle this using the underlying connection
	return nil
}

// Implementing SetReadDeadline method
func (c *SSHConn) SetReadDeadline(t time.Time) error {
	// `ssh.Session` does not provide native support for SetReadDeadline,
	// so we do nothing here or you can handle this using the underlying connection
	return nil
}

// Implementing SetWriteDeadline method
func (c *SSHConn) SetWriteDeadline(t time.Time) error {
	// `ssh.Session` does not provide native support for SetWriteDeadline,
	// so we do nothing here or you can handle this using the underlying connection
	return nil
}

func SetTitle(conn ssh.Session, message string) {
	conn.Write([]byte("\033]0;" + message + "\007"))
}