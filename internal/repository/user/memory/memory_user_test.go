package memory

import (
	"context"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	repo "github.com/FoPQer/go-shortener/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*MemoryUserRepository)
		user       *model.User
		wantID     string
		wantErr    error
		assertRepo func(*testing.T, *MemoryUserRepository, *model.User)
	}{
		{
			name:    "saves new user",
			setup:   func(*MemoryUserRepository) {},
			user:    model.NewUser("user-1"),
			wantID:  "user-1",
			wantErr: nil,
			assertRepo: func(t *testing.T, repository *MemoryUserRepository, user *model.User) {
				storedUser, err := repository.FindByID(context.Background(), "user-1")
				require.NoError(t, err)
				assert.Same(t, user, storedUser)
			},
		},
		{
			name: "returns error for duplicate id",
			setup: func(repository *MemoryUserRepository) {
				_, err := repository.Save(context.Background(), model.NewUser("user-1"))
				require.NoError(t, err)
			},
			user:    model.NewUser("user-1"),
			wantID:  "",
			wantErr: repo.ErrUserAlreadyExists,
			assertRepo: func(t *testing.T, repository *MemoryUserRepository, _ *model.User) {
				assert.Len(t, repository.users, 1)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := NewRepository()
			testCase.setup(repository)

			id, err := repository.Save(context.Background(), testCase.user)

			require.ErrorIs(t, err, testCase.wantErr)
			assert.Equal(t, testCase.wantID, id)
			testCase.assertRepo(t, repository, testCase.user)
		})
	}
}

func TestFindByID(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*MemoryUserRepository) *model.User
		id         string
		wantErr    error
		assertUser func(*testing.T, *model.User, *model.User)
	}{
		{
			name: "returns user by id",
			setup: func(repository *MemoryUserRepository) *model.User {
				user := model.NewUser("user-1")
				_, err := repository.Save(context.Background(), user)
				require.NoError(t, err)
				return user
			},
			id:      "user-1",
			wantErr: nil,
			assertUser: func(t *testing.T, wantUser *model.User, gotUser *model.User) {
				assert.Same(t, wantUser, gotUser)
			},
		},
		{
			name:    "returns error when user does not exist",
			setup:   func(*MemoryUserRepository) *model.User { return nil },
			id:      "missing-user",
			wantErr: repo.ErrUserNotFound,
			assertUser: func(t *testing.T, _ *model.User, gotUser *model.User) {
				assert.Nil(t, gotUser)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := NewRepository()
			wantUser := testCase.setup(repository)

			result, err := repository.FindByID(context.Background(), testCase.id)

			require.ErrorIs(t, err, testCase.wantErr)
			testCase.assertUser(t, wantUser, result)
		})
	}
}

func TestGetUserURLs(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*MemoryUserRepository) []*model.Urls
		userID     string
		wantErr    error
		assertURLs func(*testing.T, []*model.Urls, []*model.Urls)
	}{
		{
			name: "returns urls for existing user",
			setup: func(repository *MemoryUserRepository) []*model.Urls {
				user := model.NewUser("user-1")
				firstURL := model.NewUrls("https://example.com", "abc123")
				secondURL := model.NewUrls("https://google.com", "def456")
				user.AddURL(firstURL)
				user.AddURL(secondURL)

				_, err := repository.Save(context.Background(), user)
				require.NoError(t, err)
				return []*model.Urls{firstURL, secondURL}
			},
			userID:  "user-1",
			wantErr: nil,
			assertURLs: func(t *testing.T, wantURLs []*model.Urls, gotURLs []*model.Urls) {
				assert.Len(t, gotURLs, 2)
				assert.Same(t, wantURLs[0], gotURLs[0])
				assert.Same(t, wantURLs[1], gotURLs[1])
			},
		},
		{
			name:    "returns error when user does not exist",
			setup:   func(*MemoryUserRepository) []*model.Urls { return nil },
			userID:  "missing-user",
			wantErr: repo.ErrUserNotFound,
			assertURLs: func(t *testing.T, _ []*model.Urls, gotURLs []*model.Urls) {
				assert.Nil(t, gotURLs)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := NewRepository()
			wantURLs := testCase.setup(repository)

			result, err := repository.GetUserURLs(context.Background(), testCase.userID)

			require.ErrorIs(t, err, testCase.wantErr)
			testCase.assertURLs(t, wantURLs, result)
		})
	}
}