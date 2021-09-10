package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/config"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/convertor"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/router"
	"github.com/EpicStep/avito-autumn-2021-intern-task/pkg/database"
	"github.com/EpicStep/avito-autumn-2021-intern-task/pkg/exchangeratesapi"
	"github.com/EpicStep/avito-autumn-2021-intern-task/pkg/server"
	"github.com/go-chi/chi/v5"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cfg, err := config.New()
	if err != nil {
		return err
	}

	addr := ":" + cfg.Port

	db, err := database.New(context.Background(), cfg.PgURL)
	if err != nil {
		return err
	}

	currencyAPI := exchangeratesapi.New(cfg.EAPIToken)

	currency, err := currencyAPI.GetCurrencyList()
	if err != nil {
		return err
	}

	cConvertor := convertor.NewCurrencyConvertor(currency.Rates)

	r := router.New()

	service := balance.New(db, cConvertor)

	r.Route("/api", func(r chi.Router) {
		service.Routes(r)
	})

	srv := server.New(addr, r)

	fmt.Printf("Service has been started on %s\n", addr)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("server shutdown failed")
	}

	db.Close()

	return nil
}
