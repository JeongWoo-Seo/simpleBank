package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/JeongWoo-Seo/simpleBank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (s *Server) authorizationUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	value := md.Get(authorizationHeader)
	if len(value) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := value[0]
	field := strings.Fields(authHeader)
	if len(field) < 2 {
		return nil, fmt.Errorf("invalid auth header format")
	}

	authType := strings.ToLower(field[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := field[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
