package irc

import (
	"crypto/tls"
	"fmt"
	"net"
)

// RunServer starts the GGircd IRC server. This method will not return.
func RunServer(cfg Config) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logf(fatal, "Could not create server: %v", err)
	}

	if cfg.Prometheus.Port != 0 {
		go runPrometheus(cfg)
	}

	var lnSSL net.Listener
	certFile := cfg.SSLCertificate.CertFile
	keyFile := cfg.SSLCertificate.KeyFile
	if keyFile != "" && certFile != "" {
		getCertificate := func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return &cert, nil
			} else {
				return nil, err
			}
		}

		sslCfg := &tls.Config{GetCertificate: getCertificate}

		lnSSL, err = tls.Listen("tcp", fmt.Sprintf(":%d", cfg.SSLPort), sslCfg)
		if err != nil {
			logf(fatal, "Could not create TLS server: %v", err)
		}
	}

	state := make(chan state, 1)
	state <- newState(cfg)

	if lnSSL != nil {
		go acceptLoop(cfg, lnSSL, state)
	}
	acceptLoop(cfg, ln, state)
}

func acceptLoop(cfg Config, listener net.Listener, state chan state) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logf(warn, "Could not accept new connection: ", err)
			continue
		}

		c := newConnection(cfg, conn, newFreshHandler(state))
		go c.loop()
	}
}
