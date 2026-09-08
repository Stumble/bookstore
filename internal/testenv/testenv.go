// Package testenv owns the databases used by the bookstore test suite.
package testenv

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/stumble/wpgx"
)

type Environment struct {
	Postgres     *wpgx.Config
	RedisAddress string
}

func Start(t testing.TB) *Environment {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pgAddress := startContainer(t, ctx, "postgres:14.5", "5432", "POSTGRES_PASSWORD=bookstore-local-test", "PGTZ=America/Los_Angeles")
	redisAddress := startContainer(t, ctx, "redis:7", "6379")
	_, portText, err := net.SplitHostPort(pgAddress)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	pgConfig := &wpgx.Config{
		Username: "postgres", Password: "bookstore-local-test", Host: "127.0.0.1", Port: port,
		DBName: "testdb", AppName: "bookstore", MaxConns: 100,
		ReadReplicas: []wpgx.ReadReplicaConfig{{Name: "R1", Username: "postgres", Password: "bookstore-local-test", Host: "127.0.0.1", Port: port, DBName: "testdb", MaxConns: 100}},
	}
	uri := "postgres://postgres:bookstore-local-test@" + pgAddress + "/postgres?sslmode=disable"
	waitReady(t, ctx, func(ctx context.Context) error {
		conn, err := pgx.Connect(ctx, uri)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close(ctx) }()
		return conn.Ping(ctx)
	})
	rdb := redis.NewClient(&redis.Options{Addr: redisAddress})
	defer func() { _ = rdb.Close() }()
	waitReady(t, ctx, func(ctx context.Context) error { return rdb.Ping(ctx).Err() })
	return &Environment{Postgres: pgConfig, RedisAddress: redisAddress}
}

func startContainer(t testing.TB, ctx context.Context, image, port string, env ...string) string {
	t.Helper()
	args := []string{"run", "--detach", "--publish", "127.0.0.1::" + port, "--label", "bookstore.test=true"}
	for _, v := range env {
		args = append(args, "--env", v)
	}
	args = append(args, image)
	out, err := exec.CommandContext(ctx, "docker", args...).Output()
	if err != nil {
		t.Fatalf("start %s: %v", image, err)
	}
	id := strings.TrimSpace(string(out))
	if !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(id) {
		t.Fatalf("unexpected container ID %q", id)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if out, err := exec.CommandContext(ctx, "docker", "rm", "--force", id).CombinedOutput(); err != nil {
			t.Errorf("cleanup %s: %v: %s", id, err, out)
		}
	})
	out, err = exec.CommandContext(ctx, "docker", "port", id, port+"/tcp").Output()
	if err != nil {
		t.Fatal(err)
	}
	address := strings.TrimSpace(string(out))
	host, _, err := net.SplitHostPort(address)
	if err != nil || host != "127.0.0.1" {
		t.Fatalf("unexpected container address %q", address)
	}
	return address
}

func waitReady(t testing.TB, ctx context.Context, check func(context.Context) error) {
	t.Helper()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		attempt, cancel := context.WithTimeout(ctx, time.Second)
		err := check(attempt)
		cancel()
		if err == nil {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatal(fmt.Errorf("database readiness: %w (last: %v)", ctx.Err(), err))
			return
		case <-ticker.C:
		}
	}
}
