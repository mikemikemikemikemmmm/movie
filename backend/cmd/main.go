package main

import (
	"backend/internal"
	"backend/internal/kafka"
	"backend/internal/otel"
	"backend/internal/promethus"
	"backend/internal/redis"
	"backend/internal/sql"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func listenSign(rootCancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		rootCancel()
	}()
}
func closeAllGoroutine(rootWg *sync.WaitGroup) {
	allGoroutineDoneCh := make(chan struct{})
	go func() {
		rootWg.Wait()
		close(allGoroutineDoneCh)
	}()
	waitAllGoroutineDoneCtx, waitAllGoroutineDoneCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer waitAllGoroutineDoneCancel()
	select {
	case <-allGoroutineDoneCh:
		log.Println("所有背景工作完成")
	case <-waitAllGoroutineDoneCtx.Done():
		log.Println("等待背景工作超時，強制退出")
	}
}
func main() {
	var rootWg sync.WaitGroup
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()
	listenSign(rootCancel)

	tracerProvider, err := otel.InitTracer(rootCtx)
	if err != nil {
		log.Fatalf("tracer初始化失敗 : %v", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(rootCtx); err != nil {
			log.Fatal(err)
		}
	}()
	if err := kafka.InitKafkaWriter(); err != nil {
		log.Fatalf("kafka初始化失敗 : %v", err)
	}
	kafka.InitKafkaReader(rootCtx, &rootWg)
	if err := redis.InitRedis(); err != nil {
		log.Fatalf("redis初始化失敗 : %v", err)
	}
	if err := sql.InitSQL(); err != nil {
		log.Fatalf("Database初始化失敗 : %v", err)
	}
	promethus.InitPromethus()
	apiServer := internal.InitRouter()
	apiServerErrCh := make(chan error, 1)

	rootWg.Add(1)
	go func() {
		defer rootWg.Done()
		log.Println("HTTP server starting...")
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			apiServerErrCh <- err
		}
	}()

	select {
	case <-rootCtx.Done():
		log.Println("收到關機信號，開始優雅關機...")
	case err := <-apiServerErrCh:
		if err != nil {
			log.Printf("HTTP server error: %v", err)
			rootCancel()
		}
	}
	shutdownApiServerCtx, shutdownApiServerCtxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownApiServerCtxCancel()

	// 6. 優雅關閉 HTTP server
	if err := apiServer.Shutdown(shutdownApiServerCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shutdown complete")
	}

	// 7. 關閉其他資源
	if err := redis.RedisClient.Close(); err != nil {
		log.Printf("Redis close error: %v", err)
	}

	sqlDB, err := sql.SqlDB.DB()
	if err != nil {
		log.Printf("get sqldb error: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("SQLDB close error: %v", err)
	}

	closeAllGoroutine(&rootWg)

	log.Println("應用程式已優雅關機")
}
