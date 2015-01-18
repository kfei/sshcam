package sshd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

func generatePrivateKey(keyFile string) []byte {
	log.Println("Generating new key pair...")
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	err = ioutil.WriteFile(keyFile, pemBytes, 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Saved to ", keyFile)

	return pemBytes
}

func readLocalPrivateKey() ssh.Signer {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".sshcam"), os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	keyFile := filepath.Join(usr.HomeDir, ".sshcam", "id_rsa")

	privateBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("Private key for sshcam server does not exists")
		privateBytes = generatePrivateKey(keyFile)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	return private
}

func Run(user, pass, host, port string, sshcamArgs []string) {
	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PasswordCallback: func(c ssh.ConnMetadata,
			password []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(password) == pass {
				return nil, nil
			}
			return nil, fmt.Errorf("Password rejected for %q", c.User())
		},
	}

	config.AddHostKey(readLocalPrivateKey())

	bindAddress := strings.Join([]string{host, port}, ":")
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatal("Failed to listen on " + bindAddress)
	}

	// Accepting connections
	log.Print("Listening on " + bindAddress + "...")
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connections (%s)", err)
			continue
		}
		// Handshaking
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}

		log.Printf("New ssh connection established from %s (%s)",
			sshConn.RemoteAddr(),
			sshConn.ClientVersion())

		// Print incoming out-of-band Requests
		go handleRequests(reqs)
		// Accept all channels
		go handleChannels(chans, sshcamArgs)
	}
}

func handleRequests(reqs <-chan *ssh.Request) {
	for req := range reqs {
		log.Printf("Recieved out-of-band request: %+v", req)
	}
}

func handleChannels(chans <-chan ssh.NewChannel, sshcamArgs []string) {
	// Service the incoming Channel channel.
	for newChannel := range chans {
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType,
				fmt.Sprintf("Unknown channel type: %s", t))
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel (%s)", err)
			continue
		}

		// Allocate a terminal for this channel
		log.Print("Creating pty...")

		// Can this always work without PATH specified?
		c := exec.Command("sshcam", sshcamArgs...)

		f, err := pty.Start(c)
		if err != nil {
			log.Printf("Could not start pty (%s)", err)
			continue
		}

		// Teardown session
		var once sync.Once
		close := func() {
			channel.Close()
			_, err := c.Process.Wait()
			if err != nil {
				log.Printf("Failed to exit session (%s)", err)
			}
			log.Printf("Session closed")
		}

		// Pipe session to sshcam and visa versa
		go func() {
			io.Copy(channel, f)
			once.Do(close)
		}()
		go func() {
			io.Copy(f, channel)
			once.Do(close)
		}()

		// Deal with session requests
		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch req.Type {
				case "shell":
					// Don't accept any commands (payload)
					if len(req.Payload) == 0 {
						ok = true
					}
				case "pty-req":
					// Responding 'ok' here will let the client
					// know we have a pty ready for input
					ok = true
					// Parse body...
					termLen := req.Payload[3]
					termEnv := string(req.Payload[4 : termLen+4])
					w, h := parseDims(req.Payload[termLen+4:])
					SetWinsize(f.Fd(), w, h)
					log.Printf("pty-req '%s'", termEnv)
				case "window-change":
					w, h := parseDims(req.Payload)
					SetWinsize(f.Fd(), w, h)
					continue
				}
				if !ok {
					log.Printf("Declining %s request...", req.Type)
				}
				req.Reply(ok, nil)
			}
		}(requests)
	}
}

// Extracts two uint32s from the provided buffer
func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

// Winsize stores the Height and Width of a terminal
type Winsize struct {
	Height uint16
	Width  uint16
}

// SetWinsize sets the size of the given pty
func SetWinsize(fd uintptr, w, h uint32) {
	log.Printf("Window resize %dx%d", w, h)
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(ws)))
}
