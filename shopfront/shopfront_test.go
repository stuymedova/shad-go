package shopfront_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/shopfront"
)

func TestShopfront(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: StartRedis(t),
	})

	ctx := context.Background()

	c := shopfront.New(rdb)

	items, err := c.GetItems(ctx, []shopfront.ItemID{1, 2, 3, 4}, 42)
	require.NoError(t, err)
	require.Equal(t, items, []shopfront.Item{
		{},
		{},
		{},
		{},
	})

	require.NoError(t, c.RecordView(ctx, 3, 42))
	require.NoError(t, c.RecordView(ctx, 2, 42))

	require.NoError(t, c.RecordView(ctx, 2, 4242))

	items, err = c.GetItems(ctx, []shopfront.ItemID{1, 2, 3, 4}, 42)
	require.NoError(t, err)
	require.Equal(t, items, []shopfront.Item{
		{},
		{ViewCount: 2, Viewed: true},
		{ViewCount: 1, Viewed: true},
		{},
	})
}

func BenchmarkShopfront(b *testing.B) {
	const nItems = 1024

	rdb := redis.NewClient(&redis.Options{
		Addr: StartRedis(b),
	})

	ctx := context.Background()
	c := shopfront.New(rdb)

	var ids []shopfront.ItemID
	for i := 0; i < nItems; i++ {
		ids = append(ids, shopfront.ItemID(i))
		require.NoError(b, c.RecordView(ctx, shopfront.ItemID(i), 42))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetItems(ctx, ids, 42)
		require.NoError(b, err)
	}
}