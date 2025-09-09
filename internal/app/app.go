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
	Cfg *config.Config
	S   handler.Service
}

func NewApp() (*App, error) {
	a := App{}
	return a.Init()
}

func (a *App) Init() (*App, error) {
	a.Cfg = config.GetConfig()

	store, err := repository.NewDBStore(a.Cfg)
	if err != nil {
		return a, err
	}
	defer store.Close()

	userStore := user.NewStore(store)
	us := user.NewService(a.Cfg, userStore)

	accrualStore := accrual.NewStore(store)
	as := accrual.NewService(a.Cfg, accrualStore)

	orderStore := order.NewStore(store)
	extAs := serviceAccrual.NewService(a.Cfg, as)
	os := order.NewService(a.Cfg, orderStore, extAs)

	withdrawStore := withdraw.NewStore(store)
	ws := withdraw.NewService(a.Cfg, withdrawStore)

	l, err := logger.Initialize(a.Cfg.Logger.LogLevel)
	if err != nil {
		return a, err
	}

	a.S = handler.Service{
		Auth:            auth.New(a.Cfg.Auth.SecretKey),
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
