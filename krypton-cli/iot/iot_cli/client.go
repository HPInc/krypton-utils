package iot_cli

import (
	"cli/common"
	"cli/logging"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type SchedulerClient struct {
	Server      string
	TlsCertFile string
	ClientId    string
	Token       string
	Topic       string
	CliConfig   *autopaho.ClientConfig
	Timeout     uint // timeout in seconds
	Protocol    string
}

const (
	userNameFormat   = "username?x-amz-customauthorizer-name=KryptonIoTAuthorizer&device_token=%s"
	loginTimeout     = time.Second * 3
	subscribeTimeout = time.Second * 30
	Qos              = 1
	protocolWS       = "ws"
	websocketPort    = 9001
)

var log = logging.GetLogger()

func NewClient(server, certFile, clientId, token string, timeout uint,
	protocol string) *SchedulerClient {
	return &SchedulerClient{
		Server:      server,
		TlsCertFile: certFile,
		ClientId:    clientId,
		Token:       token,
		Timeout:     timeout,
		Protocol:    protocol,
	}
}

// create client config
func (c *SchedulerClient) GetClientConfig() (*autopaho.ClientConfig, error) {
	u, err := c.parseURL()
	if err != nil {
		return nil, err
	}
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{u},
		TlsCfg:            &tls.Config{MinVersion: tls.VersionTLS12},
		KeepAlive:         600, // aws iot keep alive values 30 to 1200
		ConnectRetryDelay: time.Second * 3,
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { log.Debug("mqtt connection up") },
		OnConnectError:    func(err error) { log.Printf("error whilst attempting connection: %s\n", err) },
		Debug:             log,
		ClientConfig: paho.ClientConfig{
			ClientID:      c.ClientId,
			OnClientError: func(err error) { log.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	certs, err := loadTlsCert(c.TlsCertFile)
	if err != nil {
		return nil, err
	}
	cliCfg.TlsCfg = &tls.Config{
		RootCAs:    certs,
		NextProtos: []string{"mqtt"},
		MinVersion: tls.VersionTLS12,
	}

	cliCfg.SetUsernamePassword(fmt.Sprintf(userNameFormat, c.Token), nil)

	return &cliCfg, nil
}

// do login
func (c *SchedulerClient) Login() error {
	cliCfg, err := c.GetClientConfig()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), loginTimeout)
	defer cancel()

	// Connect to the broker - this will return immediately after initiating the connection process
	cm, err := autopaho.NewConnection(ctx, *cliCfg)
	if err != nil {
		return err
	}

	err = cm.AwaitConnection(ctx)
	if err != nil {
		return err
	}
	return nil
}

// load cert from disk file
func loadTlsCert(rootCertPath string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()

	pemData, err := common.ReadFile(rootCertPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !certs.AppendCertsFromPEM(pemData) {
		err := errors.New("failed to append root ca cert")
		log.Println(err)
		return nil, err
	}

	return certs, nil
}

func (c *SchedulerClient) parseURL() (*url.URL, error) {
	u, err := url.ParseRequestURI(c.Server)
	if err != nil {
		return nil, err
	}

	switch c.Protocol {
	case protocolWS:
		u.Scheme = protocolWS
		u.Path = "/mqtt"
		// split host and port
		host, _, _ := net.SplitHostPort(u.Host)
		u.Host = fmt.Sprintf("%s:%d", host, websocketPort)
	}

	return u, nil
}
