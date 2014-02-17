package mailboss

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	LINE_END = "\r\n"
)

type Server struct {
	Listener *net.TCPListener
	closed   bool
}

func (s *Server) Listen(addr string) error {
	s.closed = false
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	s.Listener = listener
	for !s.closed {
		listener.SetDeadline(time.Now().Add(time.Second))
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			return err
		}
		go Handle(conn, laddr)
	}
	return listener.Close()
}

func (s *Server) Close() {
	s.closed = true
}

func Handle(conn *net.TCPConn, laddr *net.TCPAddr) error {
	conn.LocalAddr()
	if err := WriteLine(conn, fmt.Sprintf("220 %s SMTP mailboss",
		laddr.String())); err != nil {
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Second))
	recv, err := ReadLine(conn)
	if err != nil {
		return err
	}
	if err := WriteLine(conn, fmt.Sprintf("250 Hello %s, nice to meet you",
		string(recv))); err != nil {
		return err
	}
	return nil
}

func WriteLine(w io.Writer, s string) error {
	_, err := io.WriteString(w, fmt.Sprintf("%s%s", s, LINE_END))
	return err
}

func ReadLine(r io.Reader) ([]byte, error) {
	finder := &bytes.Buffer{}
	buf := &bytes.Buffer{}
	br := bufio.NewReader(r)
	for {
		c, err := br.ReadByte()
		if err != nil && err != io.EOF {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		finder.WriteByte(c)
		if !strings.HasPrefix(LINE_END, finder.String()) {
			findBytes := finder.Bytes()
			buf.WriteString(string(findBytes[:len(findBytes)-1]))
			finder = bytes.NewBuffer([]byte{c})
		} else if finder.Len() == len(LINE_END) {
			return buf.Bytes(), nil
		}
	}
	return buf.Bytes(), nil
}
