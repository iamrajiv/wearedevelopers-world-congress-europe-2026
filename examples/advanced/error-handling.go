package main

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User represents a user in the system
type User struct {
	UserID    string
	Name      string
	Email     string
	CreatedAt string
}

// GetUserRequest represents a request to get a user
type GetUserRequest struct {
	UserId string
}

// server implements the user service
type server struct {
	db *sql.DB
}

// GetUser demonstrates proper error handling with gRPC status codes
func (s *server) GetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
	// Simulate database lookup
	user, err := s.db.QueryRowContext(ctx, "SELECT user_id, name, email, created_at FROM users WHERE user_id = ?", req.UserId).Scan()

	// Handle not found error
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.UserId)
	}

	// Handle database errors
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Validate user data
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user data")
	}

	return user, nil
}

// Additional error handling examples

// ValidateEmail demonstrates input validation
func ValidateEmail(email string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	// More validation logic...
	return nil
}

// CheckPermission demonstrates authorization errors
func CheckPermission(ctx context.Context, userID string) error {
	// Check if user has permission
	hasPermission := false // simulate permission check

	if !hasPermission {
		return status.Error(codes.PermissionDenied, "user does not have permission")
	}

	return nil
}

// RateLimitCheck demonstrates rate limiting errors
func RateLimitCheck(ctx context.Context, userID string) error {
	// Check rate limit
	exceeded := true // simulate rate limit check

	if exceeded {
		return status.Error(codes.ResourceExhausted, "rate limit exceeded, please try again later")
	}

	return nil
}
