package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// extract the request headers
// extract jwt claims
// validate jwt

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// extract headers
		token := extractTokenFromHeaders(request.Headers)

		if token == "" {
			return events.APIGatewayProxyResponse{
				Body:       "missing auth token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		// parse the token to extract claims
		claims, err := parseToken(token)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "user unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// DO NOT USE IT IN PRODUCTION -> use .env or AWS Secret
		secret := "secret-test-string"
		return []byte(secret), nil
	})

	if err != nil {
		return jwt.MapClaims{}, fmt.Errorf("unauthorized %w", err)
	}

	if !token.Valid {
		return jwt.MapClaims{}, fmt.Errorf("token not valid - unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("claims unauthorized")
	}

	return claims, nil
}
