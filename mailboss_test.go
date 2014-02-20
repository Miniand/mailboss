package mailboss

import (
	"net/smtp"
	"testing"
	"time"
)

const (
	TEST_ADDRESS = ":55555"
)

func TestListen(t *testing.T) {
	s := New()
	go func() {
		err := s.Listen(TEST_ADDRESS)
		if err != nil {
			t.Fatal(err)
		}
	}()
	defer s.Close()
	time.Sleep(time.Second)
	conn, err := smtp.Dial(TEST_ADDRESS)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Hello("localhost"); err != nil {
		t.Fatal(err)
	}
	if err := conn.Rcpt("blah@localhost"); err != nil {
		t.Fatal(err)
	}
}
