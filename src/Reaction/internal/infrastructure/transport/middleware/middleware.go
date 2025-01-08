package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"tinder-reaction/internal/infrastructure/transport/model"
	trace_utils "tinder-reaction/internal/trace"

	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/otel/attribute"
)

type serviceClaims struct {
	Service string `json:"Service"`
	jwt.StandardClaims
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		slog.Info(fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, time.Since(start)))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// todo: print traceId/requestId
				writeError(w, http.StatusInternalServerError, "Internal error")
				slog.Error(
					"request error",
					slog.Any("error", err),
					slog.Any("stack", string(debug.Stack())),
				)
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, "Bearer ")
		if len(splitAuthHeader) != 2 {
			writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		token := splitAuthHeader[1]

		parsedToken, _ := jwt.ParseWithClaims(token, &serviceClaims{}, func(token *jwt.Token) (interface{}, error) {
			// todo: move to secret file
			return []byte("fjg847sdjvnjxcFHdsag38d_d8sj3aqQwfdsph3456v0bjz45ty54gpo3vhjs7234f09Odp"), nil
		})

		claims := parsedToken.Claims.(*serviceClaims)

		// todo: проверить запись спанов
		ctx := req.Context()

		if !parsedToken.Valid {
			trace_utils.AddAttributesToCurrentSpan(
				ctx,
				attribute.String("tndr.warning", "invalid jwt"),
				attribute.String("tndr.jwt", parsedToken.Raw),
			)
			writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// todo: проверять на соответствие разрешенным методам
		// todo: Admin
		if claims.Service != "GeoService" {
			trace_utils.AddAttributesToCurrentSpan(
				ctx,
				attribute.String("tndr.warning", "invalid service"),
				attribute.String("tndr.jwt", parsedToken.Raw),
				attribute.String("tndr.jwt.service", claims.Service),
			)
			writeError(w, http.StatusUnauthorized, "Now allowed")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func writeError(w http.ResponseWriter, statusCode int, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := model.ErrorResponse{Error: reason}
	json.NewEncoder(w).Encode(resp)
}
