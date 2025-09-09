package app

import (
	"net/http"

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

type App struct {
	cfg *config.Config
	S   handler.Service
	db  *repository.DBStore
}

func NewApp(cfg *config.Config, db *repository.DBStore) (*App, error) {
	a := App{
		cfg: cfg,
		db:  db,
	}
	return a.Init()
}

func (a *App) Init() (*App, error) {
	userStore := user.NewStore(a.db)
	us := user.NewService(a.cfg, userStore)

	accrualStore := accrual.NewStore(a.db)
	as := accrual.NewService(a.cfg, accrualStore)

	orderStore := order.NewStore(a.db)
	extAs := serviceAccrual.NewService(a.cfg, as)
	os := order.NewService(a.cfg, orderStore, extAs)

	withdrawStore := withdraw.NewStore(a.db)
	ws := withdraw.NewService(a.cfg, withdrawStore)

	l, err := logger.Initialize(a.cfg.Logger.LogLevel)
	if err != nil {
		return a, err
	}

	a.S = handler.Service{
		Auth:            auth.New(a.cfg.Auth.SecretKey),
		UserService:     us,
		OrderService:    os,
		WithdrawService: ws,
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
