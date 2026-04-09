package store

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreInit(t *testing.T) {
	miniRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer miniRedis.Close()

	testStoreService, err := NewStorageService(miniRedis.Addr(), "", 0)
	require.NoError(t, err)
	assert.NotNil(t, testStoreService.redisClient)
}

func TestInsertionAndRetrieval(t *testing.T) {
	miniRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer miniRedis.Close()

	testStoreService, err := NewStorageService(miniRedis.Addr(), "", 0)
	require.NoError(t, err)

	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortURL := "Jsz4k57oAX"

	err = testStoreService.SaveURLMapping(shortURL, initialLink)
	require.NoError(t, err)

	retrievedURL, err := testStoreService.RetrieveInitialURL(shortURL)
	require.NoError(t, err)

	assert.Equal(t, initialLink, retrievedURL)
}
