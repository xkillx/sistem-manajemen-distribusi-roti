package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/smdr/backend/internal/config"
	sqlcdb "github.com/smdr/backend/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type SeedService struct {
	queries *sqlcdb.Queries
	cfg     *config.Config
}

func NewSeedService(queries *sqlcdb.Queries, cfg *config.Config) *SeedService {
	return &SeedService{queries: queries, cfg: cfg}
}

func (s *SeedService) SeedOwner(ctx context.Context) error {
	_, err := s.queries.FindUserByUsername(ctx, s.cfg.OwnerUsername)
	if err == nil {
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(s.cfg.OwnerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		Username:           s.cfg.OwnerUsername,
		HashedPassword:     string(hashed),
		Name:               "Owner",
		Phone:              "",
		Role:               "owner",
		MustChangePassword: false,
	})
	return err
}
