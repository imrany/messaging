package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/imrany/smart_spore_hub/server/database/crypto"
)

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(rw, r)

		// Log request details
		duration := time.Since(start)
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", duration,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AuthMiddleware validates JWT tokens from Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Check if it starts with "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			respondError(w, http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
			return
		}

		if parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Invalid authorization type. Expected: Bearer")
			return
		}

		// Extract token
		tokenString := parts[1]
		if tokenString == "" {
			respondError(w, http.StatusUnauthorized, "Token is required")
			return
		}

		// Validate token
		_, err := crypto.ValidateToken(tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Call next handler with updated context
		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware handles CORS (Cross-Origin Resource Sharing)
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserEmailFromContext extracts user email from request context
func GetUserEmailFromContext(r *http.Request) (string, error) {
	email, ok := r.Context().Value("email").(string)
	if !ok || email == "" {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}

// GetUserRoleFromContext extracts user role from request context
func GetUserRoleFromContext(r *http.Request) (string, error) {
	role, ok := r.Context().Value("role").(string)
	if !ok || role == "" {
		return "", fmt.Errorf("role not found in context")
	}
	return role, nil
}

// RoleMiddleware checks if user has required role(s)
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, err := GetUserRoleFromContext(r)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Role not found in token")
				return
			}

			// Check if user's role is in allowed roles
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				respondError(w, http.StatusForbidden, "Insufficient permissions for this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware recovers from panics and returns 500
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"method", r.Method,
					"path", r.URL.Path,
				)
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware implements simple rate limiting
// Note: This is a basic implementation. For production, use a proper rate limiter
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	type client struct {
		requests  int
		lastReset time.Time
	}

	clients := make(map[string]*client)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			// Get or create client
			c, exists := clients[ip]
			if !exists {
				c = &client{
					requests:  0,
					lastReset: time.Now(),
				}
				clients[ip] = c
			}

			// Reset counter if minute has passed
			if time.Since(c.lastReset) > time.Minute {
				c.requests = 0
				c.lastReset = time.Now()
			}

			// Check rate limit
			if c.requests >= requestsPerMinute {
				respondError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
				return
			}

			c.requests++
			next.ServeHTTP(w, r)
		})
	}
}

// ContentTypeMiddleware ensures Content-Type is application/json for POST/PUT/PATCH
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				respondError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
