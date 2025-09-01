package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model/accrual"
	bt "github.com/spitfy/gofermart/internal/model/balance_transaction"
	"log"
	"net/http"
)

type Storer interface {
	Add(trx bt.BalanceTransaction) error
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

type Service struct {
	cfg    *config.Config
	s      Storer
	userID int
}

func (s *Service) Call(userID int, orderNumber string) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Length", "0").
		Get(fmt.Sprintf("%s/%s", s.cfg.Accrual.SystemAddress, orderNumber))
	if err != nil {
		log.Println(err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		b, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(b)
	case http.StatusNoContent:
		return
	case http.StatusTooManyRequests:
		//todo
	case http.StatusInternalServerError:
		log.Println("error response accrual")
		return
	}
}

func (s *Service) prepare(userID int, resp []byte) (bt.BalanceTransaction, error) {
	var a accrual.Response
	dec := json.NewDecoder(bytes.NewReader(resp))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return bt.BalanceTransaction{}, err
	}
	return bt.BalanceTransaction{
		OrderNum: a.Order,
		Type:     bt.TypeAccrual,
		Amount:   a.Accrual,
		UserID:   userID,
	}, nil
}

func (s *Service) save(trx bt.BalanceTransaction) {
	if err := s.s.Add(trx); err != nil {
		log.Println(err)
	}
}
