package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/domain/accrual"
	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/middleware/logger"
	"github.com/spitfy/gofermart/internal/repository"
	serviceAccrual "github.com/spitfy/gofermart/internal/service/external/accrual"
)

type ctx struct {
	ctx    context.Context
	Cancel func()
}

type App struct {
	cfg *config.Config
	S   handler.Service
	db  *repository.DBStore
	Ctx ctx
	Wg  *sync.WaitGroup
}

func NewApp(cfg *config.Config, db *repository.DBStore) (*App, error) {
	var c ctx
	c.ctx, c.Cancel = context.WithCancel(context.Background())
	a := App{
		cfg: cfg,
		db:  db,
		Ctx: c,
		Wg:  &sync.WaitGroup{},
	}
	return a.Init()
}

func (a *App) Init() (*App, error) {

	accrualStore := accrual.NewStore(a.db)
	as := accrual.NewService(a.cfg, accrualStore)
	extAs := serviceAccrual.NewService(a.Ctx.ctx, a.cfg, as, a.Wg)

	orderStore := order.NewStore(a.db)
	userStore := user.NewStore(a.db)
	withdrawStore := withdraw.NewStore(a.db)

	l, err := logger.Initialize(a.cfg.Logger.LogLevel)
	if err != nil {
		return a, err
	}

	a.S = handler.Service{
		Auth:            auth.New(a.cfg.Auth.SecretKey),
		UserService:     user.NewService(a.cfg, userStore),
		OrderService:    order.NewService(a.cfg, orderStore, extAs),
		WithdrawService: withdraw.NewService(a.cfg, withdrawStore),
		Logger:          l,
	}

	return a, nil
}

func NewServer(cfg *config.Config, h http.Handler) *http.Server {
	return &http.Server{
		Addr:    cfg.Handler.RunAddress,
		Handler: h,
	}
}
