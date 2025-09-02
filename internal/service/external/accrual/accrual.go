package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model/accrual"
	"github.com/spitfy/gofermart/internal/model/balance"
	"log"
	"net/http"
)

type Storer interface {
	Add(bt balance.BalanceTransaction) error
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
		b, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(b)
	case http.StatusNoContent:
		b, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(b)
		return
	case http.StatusTooManyRequests:
		//todo
	case http.StatusInternalServerError:
		log.Println("error response accrual")
		return
	}
}

func (s *Service) prepare(userID int, resp []byte) (balance.BalanceTransaction, error) {
	resp = []byte("{\n      \"order\": \"123\",\n      \"status\": \"PROCESSED\",\n      \"accrual\": 500\n  }")
	var a accrual.Response
	dec := json.NewDecoder(bytes.NewReader(resp))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return balance.BalanceTransaction{}, err
	}
	return balance.BalanceTransaction{
		OrderNum: a.Order,
		Type:     balance.TypeAccrual,
		Amount:   a.Accrual,
		UserID:   userID,
	}, nil
}

func (s *Service) save(bt balance.BalanceTransaction) {
	if err := s.s.Add(bt); err != nil {
		log.Println(err)
	}
}
