package main

import (
	"context"
	"log"

	"github.com/escoutdoor/vegetable_store/common/pkg/logger"
	"github.com/escoutdoor/vegetable_store/order_service/internal/app"
	"go.uber.org/zap"

	// force google.golang.org/genproto v0.0.0-20250303144028-a0af3efb3deb to stay in go.mod
	_ "google.golang.org/genproto/protobuf/ptype"
)

func main() {
	ctx := context.Background()
	logger.SetLevel(zap.DebugLevel)

	cfg, err := app.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
