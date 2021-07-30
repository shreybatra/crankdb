package main

import (
	"math/rand"
	"testing"
	"time"

	cql "github.com/ahsanbarkati/crankdb/cql"
	"github.com/ahsanbarkati/crankdb/server"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	c, err := NewCrankConnection("localhost:9876")
	require.NoError(t, err)

	c.Set("foo", "bar")

	resp, err := c.Get("foo")
	require.NoError(t, err)

	require.Equal(t, resp.DataType, cql.DataType_STRING)
	require.Equal(t, resp.StringVal, "bar")
}

func BenchmarkWriteWithClient(b *testing.B) {
	c, err := NewCrankConnection("localhost:9876")
	require.NoError(b, err)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("foo", "bar")
		}
	})
}

func BenchmarkWriteNormal(b *testing.B) {
	db := server.NewDatabase()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			db.Add("foo", "bar", cql.DataType_STRING)
		}
	})
}

func BenchmarkWriteSM(b *testing.B) {
	db := server.NewDatabase()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			db.AddSM("foo", "bar", cql.DataType_STRING)
		}
	})
}

func BenchmarkRead(b *testing.B) {
	c, err := NewCrankConnection("localhost:9876")
	require.NoError(b, err)

	c.Set("foo", "bar")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("foo")
		}
	})
}

func BenchmarkReadNormal(b *testing.B) {
	db := server.NewDatabase()
	db.Add("foo", "bar", cql.DataType_STRING)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			db.Retrieve("foo")
		}
	})
}

func BenchmarkReadSM(b *testing.B) {
	db := server.NewDatabase()
	db.AddSM("foo", "bar", cql.DataType_STRING)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			db.RetrieveSM("foo")
		}
	})
}

func BenchmarkParWriteRead(b *testing.B) {
	c, err := NewCrankConnection("localhost:9876")
	require.NoError(b, err)

	b.RunParallel(func(pb *testing.PB) {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			if rng.Intn(2) == 0 {
				c.Set("foo", "bar")
			} else {
				c.Get("foo")
			}
		}
	})
}
