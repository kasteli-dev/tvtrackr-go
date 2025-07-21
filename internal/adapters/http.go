package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasteli-dev/tvtrackr-go/internal/core"
)

type HTTPServer struct {
	userService core.UserService
	jwtSecret   []byte
	tvdbClient  *TheTVDBClient
}

func NewHTTPServer(userService core.UserService, tvdbClient *TheTVDBClient) *HTTPServer {
	// In production, use a secure random secret from config/env
	return &HTTPServer{userService: userService, jwtSecret: []byte("supersecretkey"), tvdbClient: tvdbClient}
}

func (s *HTTPServer) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/v1/register", s.handleRegister)
	r.Post("/api/v1/login", s.handleLogin)
	r.Get("/api/v1/series/search", s.handleSeriesSearch)

	// Protected routes
	r.Group(func(pr chi.Router) {
		pr.Use(s.jwtAuthMiddleware)
		pr.Post("/api/v1/series/follow", s.handleFollowSeries)
		pr.Get("/api/v1/series/followed", s.handleListFollowedSeries)
	})
	return r
}

// Series Search handler (real)
func (s *HTTPServer) handleSeriesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}
	results, err := s.tvdbClient.SearchSeries(query)
	if err != nil {
		http.Error(w, "thetvdb error: "+err.Error(), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(results)
}

// JWT middleware
func (s *HTTPServer) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := authHeader[7:]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return s.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "invalid token subject", http.StatusUnauthorized)
			return
		}
		// Store userID in context
		ctx := r.Context()
		ctx = setUserIDInContext(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string

const userIDKey contextKey = "userID"

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func getUserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	id, ok := v.(string)
	return id, ok
}

// Follow Series handler
type followSeriesRequest struct {
	SeriesID string `json:"seriesId"`
}

const errInvalidRequest = "invalid request"

func (s *HTTPServer) handleFollowSeries(w http.ResponseWriter, r *http.Request) {
	var req followSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	err := s.userService.FollowSeries(r.Context(), userID, req.SeriesID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"series followed"}`))
}

// List Followed Series handler
func (s *HTTPServer) handleListFollowedSeries(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	series, err := s.userService.ListFollowedSeries(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(series)
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func toUserResponse(u *core.User) userResponse {
	return userResponse{
		ID:       u.ID,
		Username: u.Username,
	}
}

func (s *HTTPServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}
	user, err := s.userService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toUserResponse(user))
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}
	user, err := s.userService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{
		User:  toUserResponse(user),
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(resp)
}
