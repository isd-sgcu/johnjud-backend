package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/internal/cache"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CacheRepositoryTest struct {
	suite.Suite
	cacheDb   *redis.Client
	cacheRepo cache.Repository
	key       string
	value     *dto.AccessTokenCache
}

func TestCacheRepository(t *testing.T) {
	suite.Run(t, new(CacheRepositoryTest))
}

func (t *CacheRepositoryTest) SetupTest() {
	addr := fmt.Sprintf("%s:%d", "localhost", 6379)
	cacheDb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	cacheRepo := cache.NewRepository(cacheDb)
	key := faker.UUIDDigit()
	value := &dto.AccessTokenCache{
		Token:        faker.Word(),
		RefreshToken: faker.UUIDDigit(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := cacheDb.FlushDB(ctx).Err()
	assert.Nil(t.T(), err)

	err = cacheRepo.SetValue(key, value, 60)
	assert.Nil(t.T(), err)

	t.cacheDb = cacheDb
	t.cacheRepo = cacheRepo
	t.key = key
	t.value = value
}

func (t *CacheRepositoryTest) TestSetValueSuccess() {
	key := faker.UUIDDigit()
	value := &dto.AccessTokenCache{
		Token:        faker.Word(),
		RefreshToken: faker.UUIDDigit(),
	}
	err := t.cacheRepo.SetValue(key, value, 60)
	assert.Nil(t.T(), err)
}

func (t *CacheRepositoryTest) TestGetValueSuccess() {
	value := &dto.AccessTokenCache{}
	err := t.cacheRepo.GetValue(t.key, value)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), t.value, value)
}

func (t *CacheRepositoryTest) TestDeleteValueSuccess() {
	err := t.cacheRepo.DeleteValue(t.key)
	assert.Nil(t.T(), err)
}
