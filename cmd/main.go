package main

import (
	"crypto/tls"
	"hotwire/pkg/api"
	"hotwire/pkg/hotwire"
	"hotwire/pkg/log"
	"hotwire/pkg/servers/jdocs"
	"hotwire/pkg/servers/token"
	"hotwire/pkg/stts/vosk"
	"hotwire/pkg/vars"
	"net"
	"net/http"
	"os"

	chipperserver "hotwire/pkg/servers/chipper"

	chipperpb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/api/go/jdocspb"
	"github.com/digital-dream-labs/api/go/tokenpb"
	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
)

func main() {
	certPub, err := os.ReadFile(vars.CertPath)
	if err != nil {
		panic(err)
	}
	certPriv, err := os.ReadFile(vars.KeyPath)
	if err != nil {
		panic(err)
	}
	cert, err := tls.X509KeyPair(certPub, certPriv)
	if err != nil {
		panic(err)
	}

	srv, err := grpcserver.New(
		grpcserver.WithViper(),
		grpcserver.WithReflectionService(),
		grpcserver.WithCertificate(cert),
		grpcserver.WithClientAuth(tls.RequestClientCert),
	)
	if err != nil {
		panic(err)
	}
	p, err := hotwire.New(vosk.NewVoskSTT())
	if err != nil {
		panic(err)
	}
	s, _ := chipperserver.New(
		chipperserver.WithIntentProcessor(p),
		chipperserver.WithKnowledgeGraphProcessor(p),
		chipperserver.WithIntentGraphProcessor(p),
	)

	tokenServer := token.NewTokenServer()
	jdocsServer := jdocs.NewJdocsServer()
	//jdocsserver.IniToJson()

	chipperpb.RegisterChipperGrpcServer(srv.Transport(), s)
	jdocspb.RegisterJdocsServer(srv.Transport(), jdocsServer)
	tokenpb.RegisterTokenServer(srv.Transport(), tokenServer)

	listenerOne, err := net.Listen("tcp", ":4045")
	if err != nil {
		panic(err)
	}
	log.Normal("started")
	go srv.Transport().Serve(listenerOne)
	api.InitConncheck()
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Error("ListenAndServe error (:80):", err)
		os.Exit(1)
	}
}
