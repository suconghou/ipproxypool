package proxy

import (
	"bytes"
	"fmt"
	"io"
	"ipproxypool/util"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/proxy"
)

func init() {
	var address = os.Getenv("PROXY_LISTEN")
	if address != "" {
		go func() {
			if err := serve(address); err != nil {
				util.Log.Print(err)
			}
		}()
	}
}

func serve(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		client, err := l.Accept()
		if err != nil {
			util.Log.Print(err)
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					util.Log.Print(err)
				}
			}()
			err := proxyRequest(client, getdialer())
			if err != nil && err != io.EOF {
				util.Log.Print(err)
			}
		}()
	}
}

func getdialer() proxy.Dialer {
	return proxy.FromEnvironment()
}

func proxyRequest(client net.Conn, dialer proxy.Dialer) error {
	var b [32]byte
	n, err := client.Read(b[:])
	if err != nil {
		return err
	}
	if b[0] == 5 { // socks5
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		if err != nil {
			return err
		}
		var host, port string
		switch b[3] {
		case 0x01: //IPV4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IPV6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		conn, err := dialer.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return err
		}
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		return util.IoCopy(conn, client)
	}
	// try https_proxy and http_proxy
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		return err
	}
	if method == "CONNECT" { // https_proxy
		_, err = client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			return err
		}
		address = fmt.Sprintf("%s", hostPortURL)
	} else { // at last parse as http_proxy
		address = hostPortURL.Host
		if strings.Index(address, ":") == -1 {
			address = address + ":80"
		}
	}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return err
	}
	if method != "CONNECT" {
		_, err = conn.Write(b[:n])
		if err != nil {
			return err
		}
	}
	return util.IoCopy(conn, client)
}
