package ssh

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func StartServer(privateKey []byte, authorizedKey []byte) error {
	authorizedKeysMap := map[string]bool{}
	for len(authorizedKey) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKey)
		if err != nil {
			return fmt.Errorf("unable to parse authorized key: %v", err)
		}

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKey = rest
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	private, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	config.AddHostKey(private)

	fmt.Println("Listening on port 2023")
	listener, err := net.Listen("tcp", "0.0.0.0:2023")
	if err != nil {
		return fmt.Errorf("unable to listen: %v", err)
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("unable to accept: %v\n", err)
		}

		conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			fmt.Printf("unable to handshake: %v\n", err)
		}
		if conn != nil && conn.Permissions != nil {
			log.Printf("logged in with key: %s", conn.Permissions.Extensions["pubkey-fp"])
		}

		go ssh.DiscardRequests(reqs)

		go handleConnection(conn, chans)

	}
}

func handleConnection(conn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("could not accept channel: %v", err)
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Printf("request type made by client: %s\n", req.Type)
				switch req.Type {
				case "exec":
					payload := bytes.TrimPrefix(req.Payload, []byte{0, 0, 0, 6})
					channel.Write([]byte(execCommand(conn, payload)))
					channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					req.Reply(true, nil)
					channel.Close()
				case "shell":
					req.Reply(req.Type == "shell", nil)
				case "pty-req":
					createTerminal(conn, channel)
					req.Reply(true, nil)
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}

func createTerminal(conn *ssh.ServerConn, channel ssh.Channel) {
	termInstace := term.NewTerminal(channel, "> ")
	go func() {
		defer channel.Close()
		termInstace.Write([]byte("Welcome to the server\n"))
		for {
			line, err := termInstace.ReadLine()
			if err != nil {
				fmt.Printf("unable to read line: %v", err)
				break
			}
			switch line {
			case "whoami":
				termInstace.Write([]byte(execCommand(conn, []byte("whoami"))))
			case "exit":
				termInstace.Write([]byte("Bye\n"))
				return
			case "clear":
				termInstace.Write([]byte("\033[H\033[2J"))
			case "ls":
				termInstace.Write([]byte("No files\n"))
			case "help":
				termInstace.Write([]byte("Available commands:\n"))
				termInstace.Write([]byte("whoami\n"))
				termInstace.Write([]byte("exit\n"))
			default:
				termInstace.Write([]byte(fmt.Sprintf("Unknown command: %s\n", line)))
			}
		}
	}()
}

func execCommand(conn *ssh.ServerConn, payload []byte) string {
	switch string(payload) {
	case "whoami":
		return fmt.Sprintf("You are %s\n", conn.Conn.User())
	default:
		return fmt.Sprintf("Unknown command: %s\n", string(payload))
	}
}
