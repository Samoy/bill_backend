package main

import (
	"context"
	"fmt"
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	config.Setup("./app.ini")
	dao.Setup(config.DatabaseConf.Type,
		config.DatabaseConf.User,
		config.DatabaseConf.Password,
		config.DatabaseConf.Host,
		config.DatabaseConf.Name,
		config.DatabaseConf.TablePrefix,
	)
}

func main() {
	r := router.InitRouter(config.AppConf.RunMode)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServerConf.HTTPPort),
		Handler:      r,
		ReadTimeout:  config.ServerConf.ReadTimeout,
		WriteTimeout: config.ServerConf.WriteTimeout,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%v\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown:%v\n", err)
	}
	log.Printf("Server exiting")
}
