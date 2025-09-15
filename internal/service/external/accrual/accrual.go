package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/domain/accrual"
)

type Storer interface {
	Add(a accrual.Accrual) error
}

func NewService(ctx context.Context, cfg *config.Config, as *accrual.Service, wg *sync.WaitGroup) *Service {
	t := time.Now()
	s := Service{
		cfg: cfg,
		s:   as,
		Ctx: ctx,
		Wg:  wg,
	}
	s.Await.Store(t)
	return &s
}

type Service struct {
	cfg   *config.Config
	s     *accrual.Service
	Ctx   context.Context
	Wg    *sync.WaitGroup
	Await atomic.Time
}

func (s *Service) Call(userID int, orderNumber string) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Length", "0").
		Get(fmt.Sprintf("%s/api/orders/%s", s.cfg.Accrual.SystemAddress, orderNumber))
	if err != nil {
		log.Println(err)
		return
	}
	if resp == nil {
		log.Println("received nil response from accrual service")
		return
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		log.Println("========= accrual StatusOK orderNumber: ", orderNumber)
		a, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(a)
		return
	case http.StatusNoContent:
		log.Println("========= accrual StatusNoContent orderNumber: ", orderNumber)
		return
	case http.StatusTooManyRequests:
		log.Println("========= accrual StatusTooManyRequests orderNumber: ", orderNumber)
		retryAfter := resp.Header().Get("Retry-After")
		delay, err := strconv.Atoi(retryAfter)
		if err != nil || delay <= 0 {
			delay = 1
		}
		t := time.Now().Add(time.Duration(delay) * time.Second)
		s.Await.Store(t)
	case http.StatusInternalServerError:
		log.Println("========= accrual StatusInternalServerError orderNumber: ", orderNumber)
		return
	default:
		return
	}
}

func (s *Service) prepare(userID int, resp []byte) (accrual.Accrual, error) {
	//todo
	//resp = []byte("{\n      \"order\": \"48212146351759\",\n      \"status\": \"PROCESSED\",\n      \"accrual\": 500\n  }")
	var a Response
	dec := json.NewDecoder(bytes.NewReader(resp))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return accrual.Accrual{}, err
	}
	if a.Status == "" {
		a.Status = accrual.StatusNew
	}
	return accrual.Accrual{
		UserID: userID,
		Amount: a.Accrual,
		Number: a.Order,
		Status: a.Status,
	}, nil
}

func (s *Service) save(a accrual.Accrual) {
	if err := s.s.Add(context.Background(), a); err != nil {
		log.Println(err)
	}
}
