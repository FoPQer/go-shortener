package handlers

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/logger"
	pb "github.com/FoPQer/go-shortener/internal/proto"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/FoPQer/go-shortener/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPCHandler implements the ShortenerServiceServer gRPC interface.
type GRPCHandler struct {
	pb.UnimplementedShortenerServiceServer
	urlService  *service.URLService
	userService *service.UserService
}

// NewGRPCHandler constructs a GRPCHandler with URL and user services.
func NewGRPCHandler(urlService *service.URLService, userService *service.UserService) *GRPCHandler {
	return &GRPCHandler{urlService: urlService, userService: userService}
}

// ShortenURL creates a short URL from the given full URL.
func (h *GRPCHandler) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	userID := getUserIDFromContext(ctx)
	target, err := h.urlService.SetURL(ctx, req.Url, userID)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "URL already exists")
	} else if err != nil {
		logger.GetSugar().Errorf("gRPC ShortenURL error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to shorten URL: %v", err)
	}

	return &pb.URLShortenResponse{Result: target}, nil
}

// ExpandURL resolves a short URL ID to the original URL.
func (h *GRPCHandler) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	original, err := h.urlService.GetURL(ctx, req.Id)
	if errors.Is(err, urls.ErrURLNotFound) {
		return nil, status.Error(codes.NotFound, "URL not found")
	} else if errors.Is(err, urls.ErrURLDeleted) {
		return nil, status.Error(codes.NotFound, "URL has been deleted")
	} else if err != nil {
		logger.GetSugar().Errorf("gRPC ExpandURL error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to expand URL: %v", err)
	}

	return &pb.URLExpandResponse{Result: original}, nil
}

// ListUserURLs returns all non-deleted URLs belonging to the authenticated user.
func (h *GRPCHandler) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user ID")
	}

	userURLs, err := h.urlService.GetUrlsByUserID(ctx, userID)
	if err != nil {
		logger.GetSugar().Errorf("gRPC ListUserURLs error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user URLs: %v", err)
	}

	urlData := make([]*pb.URLData, 0, len(userURLs))
	for _, u := range userURLs {
		short, err := service.MakeShortURL(u.GetShortURL())
		if err != nil {
			logger.GetSugar().Warnf("gRPC ListUserURLs: failed to build short URL for %s: %v", u.GetShortURL(), err)
			continue
		}
		urlData = append(urlData, &pb.URLData{
			ShortUrl:    short,
			OriginalUrl: u.GetOriginal(),
		})
	}

	return &pb.UserURLsResponse{Url: urlData}, nil
}
