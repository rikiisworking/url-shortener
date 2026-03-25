package service

import (
	"errors"
	"time"

	"context"

	"gitgub.com/rikiisworking/url-shortener/internal/model"
	"gitgub.com/rikiisworking/url-shortener/internal/storage"
	"gitgub.com/rikiisworking/url-shortener/internal/util"
	"github.com/jackc/pgx/v5"
)

type URLService struct {
	postgresRepo *storage.PostgresRepo
	redisRepo    *storage.RedisRepo
	codeLength   int
}

func NewURLService(p *storage.PostgresRepo, r *storage.RedisRepo, codeLen int) *URLService {
	return &URLService{
		postgresRepo: p,
		redisRepo:    r,
		codeLength:   codeLen,
	}
}

func (s *URLService) Shorten(ctx context.Context, originalURL string) (string, error) {
	for i := 0; i < 5; i++ {
		shortCode, err := util.GenerateShortCode(s.codeLength)
		if err != nil {
			return "", err
		}

		_, err = s.postgresRepo.GetByShortCode(ctx, shortCode)
		if errors.Is(err, pgx.ErrNoRows) {
			url := &model.URL{
				OriginalURL: originalURL,
				ShortCode:   shortCode,
			}

			if err := s.postgresRepo.Create(ctx, url); err != nil {
				return "", err
			}

			s.redisRepo.Set(shortCode, originalURL, 24*time.Hour)
			return shortCode, nil
		}
	}
	return "", errors.New("failed to generate unique short code after retries")
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	if url, err := s.redisRepo.Get(shortCode); err == nil {
		go func() {
			_ = s.redisRepo.IncrementClick(shortCode)
			_ = s.postgresRepo.IncrementClick(context.Background(), shortCode)
		}()
		return url, nil
	}

	u, err := s.postgresRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	s.redisRepo.Set(shortCode, u.OriginalURL, 24*time.Hour)

	go func() {
		_ = s.postgresRepo.IncrementClick(context.Background(), shortCode)
	}()

	return u.OriginalURL, nil
}
