package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotwire/pkg/log"
	"hotwire/pkg/vars"

	"github.com/digital-dream-labs/api/go/tokenpb"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type TokenServer struct {
	tokenpb.UnimplementedTokenServer
}

var (
	TimeFormat     = time.RFC3339Nano
	ExpirationTime = time.Hour * 24
)

// returns session token, session cert, robot name ("Vector-####"), then thing ("vic:esn")
func getBotDetailsFromTokReq(ctx context.Context, req *tokenpb.AssociatePrimaryUserRequest) (token string, cert []byte, name string, esn string, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", nil, "", "", errors.New("no peer info found in context")
	}
	if p.AuthInfo != nil {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.PeerCertificates) == 0 {
				return "", nil, "", "", errors.New("no peer certificates found")
			}
			clientCert := tlsInfo.State.PeerCertificates[0]
			esn = clientCert.Subject.CommonName
		}
	}
	cert = req.SessionCertificate
	block, _ := pem.Decode(cert)
	certParsed, err := x509.ParseCertificate(block.Bytes)
	if err == nil {
		name = certParsed.Issuer.CommonName
	}

	// get metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil, "", "", errors.New("no metadata found in context")
	}
	token = md["anki-user-session"][0]
	return token, cert, name, esn, nil
}

func getIPfromReq(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", errors.New("unable to get peer details from ctx (getipfromreq)")
	}
	return strings.Split(p.Addr.String(), ":")[0], nil
}

func GenJWT(genSTS bool, genJWT bool, userID, esnThing string) *tokenpb.TokenBundle {
	bundle := &tokenpb.TokenBundle{}

	currentTime := time.Now().Format(TimeFormat)
	expiresAt := time.Now().AddDate(0, 6, 0).Format(TimeFormat)

	if genJWT {
		var tokenJson ClientTokenManager
		guid, tokenHash, _ := CreateTokenAndHashedToken()
		ajdoc, err := vars.ReadJdoc(vars.Thingifier(esnThing), "vic.AppTokens")
		if err != nil {
			log.Debug("new vic.AppTokens jdoc:", err)
			ajdoc.DocVersion = 1
			ajdoc.FmtVersion = 1
			ajdoc.ClientMetadata = "wirepod-new-token"
		}
		json.Unmarshal([]byte(ajdoc.JsonDoc), &tokenJson)
		var clientToken ClientToken
		clientToken.IssuedAt = time.Now().Format(TimeFormat)
		clientToken.ClientName = "hotwire"
		clientToken.Hash = tokenHash
		clientToken.AppId = "SDK"
		tokenJson.ClientTokens = append(tokenJson.ClientTokens, clientToken)
		var finalTokens []ClientToken
		// limit tokens to 6, don't fill the db
		if len(tokenJson.ClientTokens) == 6 {
			log.Debug("shaving a token off the top (", esnThing, ")")
			for i, tok := range tokenJson.ClientTokens {
				if i != 0 {
					finalTokens = append(finalTokens, tok)
				}
			}
			tokenJson.ClientTokens = finalTokens
		}
		jdocJsoc, err := json.Marshal(tokenJson)
		ajdoc.JsonDoc = string(jdocJsoc)
		ajdoc.DocVersion++
		vars.WriteJdoc(vars.Thingifier(esnThing), "vic.AppTokens", ajdoc)

		bundle.ClientToken = guid

		requestUUID := uuid.New().String()
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
			"expires":      expiresAt,
			"iat":          currentTime,
			"permissions":  nil,
			"requestor_id": esnThing,
			"token_id":     requestUUID,
			"token_type":   "user+robot",
			"user_id":      userID,
		})
		rsaKey, _ := rsa.GenerateKey(rand.Reader, 1024)
		tokenString, _ := jwtToken.SignedString(rsaKey)
		bundle.Token = tokenString
	}
	if genSTS {
		bundle.StsToken = &tokenpb.StsToken{
			AccessKeyId:     "placeholder",
			SecretAccessKey: "placeholder",
			SessionToken:    "placeholder",
			Expiration:      expiresAt,
		}
	}
	return bundle
}

// this is something no one should ever do
func decodeJWT(tokenString string) (string, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid token structure")
	}

	headerPart := parts[0]
	payloadPart := parts[1]

	headerBytes, err := base64.RawURLEncoding.DecodeString(headerPart)
	if err != nil {
		return "", "", fmt.Errorf("error decoding header: %w", err)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return "", "", fmt.Errorf("error decoding payload: %w", err)
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", "", fmt.Errorf("error unmarshaling header JSON: %w", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", "", fmt.Errorf("error unmarshaling payload JSON: %w", err)
	}

	requestorID, ok := payload["requestor_id"].(string)
	if !ok {
		return "", "", errors.New("payload missing 'requestor_id'")
	}

	userID, ok := payload["user_id"].(string)
	if !ok {
		return "", "", errors.New("payload missing 'user_id'")
	}

	return requestorID, userID, nil
}

func (s *TokenServer) AssociatePrimaryUser(ctx context.Context, req *tokenpb.AssociatePrimaryUserRequest) (*tokenpb.AssociatePrimaryUserResponse, error) {
	_, cert, name, esn, err := getBotDetailsFromTokReq(ctx, req)
	if err != nil {
		return nil, err
	}
	ip, err := getIPfromReq(ctx)
	if err != nil {
		return nil, err
	}
	log.Important("Robot being authenticated. ESN: "+esn, ", name: "+name)
	log.Debug("incoming primary user")
	log.SuperDebug(cert, name, esn, err)
	thing := esn
	esn = strings.TrimPrefix(esn, "vic:")
	os.Mkdir(vars.SessionCertsStorage, 0777)
	os.WriteFile(filepath.Join(vars.SessionCertsStorage, name+"_"+esn), cert, 0777)
	bundle := GenJWT(req.GenerateStsToken, true, "hotwire", thing)
	log.SuperDebug(bundle)
	vars.SaveRobot(
		vars.Robot{
			Active:      true,
			TSLC:        0,
			ResetTimer:  true,
			IP:          ip,
			ESN:         esn,
			CurrentGUID: bundle.ClientToken,
			Name:        name,
		},
	)
	return &tokenpb.AssociatePrimaryUserResponse{
		Data: bundle,
	}, nil
}

func (s *TokenServer) AssociateSecondaryClient(ctx context.Context, req *tokenpb.AssociateSecondaryClientRequest) (*tokenpb.AssociateSecondaryClientResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no request metadata")
	}
	jwtToken := md["anki-access-token"]
	log.SuperDebug(md)
	thing, userId, err := decodeJWT(jwtToken[0])
	if err != nil {
		return nil, err
	}
	ip, err := getIPfromReq(ctx)
	esn := strings.Split(thing, ":")[1]
	log.Important("Robot being authenticated (secondary). ESN: " + thing)
	log.Debug("Incoming secondary client")
	log.SuperDebug(jwtToken[0])
	bundle := GenJWT(false, true, userId, thing)
	var name string
	rob, err := vars.GetRobot(ip, esn, "")
	if err == nil {
		log.Debug("using "+rob.Name+" as name for secondary client response (ip, esn:", ip, esn+")")
		name = rob.Name
	} else {
		log.Debug("this robot's userdata was not cleared. i still need to figure out a way to find the name")
	}
	log.SuperDebug("this bot:", ip, esn, bundle.ClientToken, name)
	vars.SaveRobot(
		vars.Robot{
			Active:      true,
			TSLC:        0,
			ResetTimer:  true,
			IP:          ip,
			ESN:         esn,
			CurrentGUID: bundle.ClientToken,
			Name:        name,
		},
	)
	return &tokenpb.AssociateSecondaryClientResponse{
		Data: bundle,
	}, nil
}

// INSECURE!
// i don't have a way to verify the incoming JWT, unless i save the generated key from the primary request.. that's an idea
func (s *TokenServer) RefreshToken(ctx context.Context, req *tokenpb.RefreshTokenRequest) (*tokenpb.RefreshTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no request metadata")
	}
	log.Debug("Incoming refresh token")
	log.SuperDebug(md)
	jwtToken := md["anki-access-token"]
	thing, userId, err := decodeJWT(jwtToken[0])
	if err != nil {
		return nil, err
	}
	log.SuperDebug(jwtToken, thing, userId)
	ip, err := getIPfromReq(ctx)
	if err != nil {
		return nil, err
	}
	esn := strings.Split(thing, ":")[1]
	bundle := GenJWT(req.RefreshStsTokens, req.RefreshJwtTokens, userId, thing)
	var name string
	rob, err := vars.GetRobot(ip, esn, "")
	if err == nil {
		log.Debug("using "+rob.Name+" as name for secondary client response (ip, esn:", ip, esn+")")
		name = rob.Name
	}
	log.SuperDebug("this bot:", ip, esn, bundle.ClientToken, name)
	vars.SaveRobot(
		vars.Robot{
			Active:      true,
			TSLC:        0,
			ResetTimer:  true,
			IP:          ip,
			ESN:         esn,
			CurrentGUID: bundle.ClientToken,
			Name:        name,
		},
	)
	return &tokenpb.RefreshTokenResponse{
		Data: bundle,
	}, nil
}

func NewTokenServer() *TokenServer {
	return &TokenServer{}
}
