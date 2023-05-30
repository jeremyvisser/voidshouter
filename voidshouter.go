package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"unicode"
)

var (
	iface = flag.String("interface", "", "the network interface you want to shout from")
	addr  = flag.String("address", "[ff02::401d]:16413", "the multicast address:port you want to shout to")
	nick  = flag.String("nick", "", "nickname you will use whilst shouting")
)

const mtu = 400
const mru = 3 * mtu

type voidShouter struct {
	conn *net.UDPConn
	addr *net.UDPAddr

	Nick string
}

// Receive returns three channels:
//
//	recv (messages received),
//	done (close when done),
//	errc (errors encountered)
//
// recv and errc must be read from, but not closed.
// done must be closed when finished.
func (v *voidShouter) Receive() (recv chan string, done chan struct{}, errc chan error) {
	recv = make(chan string)
	done = make(chan struct{})
	errc = make(chan error)
	go func() {
		defer close(recv)
		defer close(errc)
		var buf = make([]byte, mtu)
		for {
			_, addr, err := v.conn.ReadFromUDP(buf)
			if err != nil {
				errc <- err
				return
			}
			dec := gob.NewDecoder(bytes.NewReader(buf))
			var msg voidMsg
			for dec.Decode(&msg) == nil {
				recv <- fmt.Sprintf("%s <%s> %s", addr.IP, msg.Nick, msg.Message)
				select {
				case <-done:
					return
				default:
				}
			}
		}
	}()
	return recv, done, errc
}

func (v *voidShouter) Send(s string) error {
	if len(s) > mtu {
		return fmt.Errorf("message too long (%d, must be <%d)", len(s), mtu)
	}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(voidMsg{
		Nick:    v.Nick,
		Message: s,
	})
	if err != nil {
		return err
	}
	_, err = v.conn.WriteToUDP(buf.Bytes(), v.addr)
	return err
}

type voidMsg struct {
	Nick    string
	Message string
}

func NewVoidShouter(addr string, iface string, nick string) (v *voidShouter, err error) {
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	var ifi *net.Interface
	if iface != "" {
		ifi, err = net.InterfaceByName(iface)
		if err != nil {
			return nil, err
		}
	} else {
		ifi = nil
	}

	conn, err := net.ListenMulticastUDP("udp", ifi, uaddr)
	if err != nil {
		return nil, err
	}

	v = &voidShouter{
		conn,
		uaddr,
		nick,
	}
	return v, nil
}

// InsertWriter assumes a VT220-compatible terminal, and prepends each write
// before the current line in the terminal, saving/restoring the cursor.
//
// It only works properly when the current line hasn't wrapped. To handle line
// wrapping, you'll need more comprehensive Readline-like handling.
type InsertWriter struct {
	io.Writer
}

const (
	SaveCur    = "\x1b7"
	RestoreCur = "\x1b8"

	ModeInsert  = "\x1b[4h"
	ModeReplace = "\x1b[4l"

	CursorUp   = "\x1b[A"
	CursorDown = "\x1b[B"

	IndexUp   = "\x1bM"
	IndexDown = "\x1bD"

	InsertLine = "\x1b[L"
	DeleteLine = "\x1b[M"
)

func (w *InsertWriter) Write(p []byte) (n int, err error) {
	nl := ""
	if len(p) > 0 && p[len(p)-1] != '\n' {
		nl = "\n"
	}
	// Why do I IndexDown+IndexUp instead of a scroll down (\x1b[S)?
	// Because scroll down messes up the cursor position.
	//
	// That said, I'm confident the below isn't optimal:
	s := ModeInsert + SaveCur + IndexDown + IndexUp + InsertLine +
		string(p) + nl + RestoreCur + CursorDown + ModeReplace
	return w.Writer.Write([]byte(s))
}

func Sanitize(s string) (safe string) {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}

func main() {
	flag.Parse()

	if *nick == "" {
		*nick = os.Getenv("USER")
	}
	if *nick == "" {
		*nick = "nobody"
	}
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	*nick = *nick + "@" + host

	v, err := NewVoidShouter(*addr, *iface, *nick)
	if err != nil {
		log.Fatal("NewVoidShouter: ", err)
	}

	log.Printf("listening on %s %s", *addr, *iface)

	recv, rdone, rerr := v.Receive()
	defer close(rdone)

	iw := InsertWriter{Writer: os.Stdout}
	log.SetOutput(&iw)

	go func() {
		for {
			select {
			case s := <-recv:
				log.Print(Sanitize(s))
			case err := <-rerr:
				log.Print("recv err: ", err)
			}
		}
	}()

	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("<%s> ", v.Nick)
		if !sc.Scan() {
			break
		}

		if err := v.Send(sc.Text()); err != nil {
			log.Print("send err: ", err)
		} else {
			fmt.Print(IndexUp, DeleteLine)           // clear user prompt
			log.Printf("<%s> %s", v.Nick, sc.Text()) // echo what we heard
		}
	}

	fmt.Println() // ensure shell starts on a newline
}
