package mailboss

import (
	"net"
	"net/textproto"
	"time"
)

const (
	LINE_END = "\r\n"
)

type Server struct {
	Addr        string
	Handler     Handler
	ReadTimeout time.Duration
}

func Listen(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.Listen()
}

func (s *Server) Listen() error {
	addr := s.Addr
	if addr == "" {
		addr = ":smtp"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(ln.(*net.TCPListener))
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration
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
		go Handle(textproto.NewConn(conn), laddr)
	}
	return listener.Close()
}

func (s *Server) Close() {
	s.closed = true
}

func Handle(conn *textproto.Conn, laddr *net.TCPAddr) error {
	if err := conn.Writer.PrintfLine("220 %s SMTP mailboss",
		laddr.String()); err != nil {
		return err
	}
	line, err := conn.ReadLine()
	if err != nil {
		return err
	}
	if err := conn.Writer.PrintfLine("250 Hello %s, nice to meet you",
		line); err != nil {
		return err
	}
	return nil
}
