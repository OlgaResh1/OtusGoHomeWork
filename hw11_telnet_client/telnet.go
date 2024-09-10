package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var errConnectionIsNil = errors.New("connection not inited")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tcpClient *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tcpClient.address, tcpClient.timeout)
	if err != nil {
		return err
	}
	tcpClient.connection = conn
	return nil
}

func (tcpClient *telnetClient) Close() error {
	if tcpClient.connection == nil {
		return errConnectionIsNil
	}
	return tcpClient.connection.Close()
}

func (tcpClient *telnetClient) Send() error {
	if tcpClient.connection == nil {
		return errConnectionIsNil
	}
	for {
		buf := make([]byte, 1000)
		readed, err := tcpClient.in.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		for {
			writed, err := tcpClient.connection.Write(buf[:readed])
			if err != nil {
				return err
			}
			if writed == readed {
				break
			}
			buf = buf[writed:]
		}
	}
}

func (tcpClient *telnetClient) Receive() error {
	if tcpClient.connection == nil {
		return errConnectionIsNil
	}
	for {
		buf := make([]byte, 1000)
		readed, err := tcpClient.connection.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		for {
			writed, err := tcpClient.out.Write(buf[:readed])
			if err != nil {
				return err
			}
			if writed == readed {
				break
			}
			buf = buf[writed:]
		}
	}
}
