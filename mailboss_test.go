package mailboss

import (
	"bytes"
	"io"
	"net"
	"regexp"
	"testing"
	"time"
)

const (
	TEST_ADDRESS = ":55555"
	SEND         = iota
	RECV
)

type msg struct {
	kind    int
	content string
}

func convo(t *testing.T, expected []msg) {
	s := New()
	go func() {
		err := s.Listen(TEST_ADDRESS)
		if err != nil {
			t.Fatal(err)
		}
	}()
	defer s.Close()
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", TEST_ADDRESS)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range expected {
		conn.SetDeadline(time.Now().Add(time.Second))
		switch m.kind {
		case SEND:
			_, err = io.WriteString(conn, m.content+LINE_END)
			if err != nil {
				t.Fatal(err)
			}
		case RECV:
			recv, err := ReadLine(conn)
			if err != nil {
				t.Fatalf("Error while expecting:\n%s\n%s", m.content, err)
			}
			if !regexp.MustCompile(m.content).Match(recv) {
				t.Fatalf("Expected:\n%s\nGot:\n%s", m.content, string(recv))
			}
		}
	}
}

func TestListen(t *testing.T) {
	convo(t, []msg{
		msg{RECV, "^220"},
		msg{SEND, "HELO localhost"},
		msg{RECV, "^250 Hello localhost"},
	})
}

func TestReadLine(t *testing.T) {
	str := "bkag blag \r bla \n mooo \n \n\n\r\r\r blah"
	line, err := ReadLine(bytes.NewBufferString(str + "\r\nmee"))
	if err != nil {
		t.Fatal(err)
	}
	if str != string(line) {
		t.Fatalf("Expected:\n%s\nGot:\n%s", str, string(line))
	}
}
