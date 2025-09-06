package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/domain/accrual"
	"log"
	"net/http"
)

type Storer interface {
	Add(a accrual.Accrual) error
}

func NewService(cfg *config.Config, s *accrual.Service) *Service {
	return &Service{
		cfg: cfg,
		s:   s,
	}
}

type Service struct {
	cfg    *config.Config
	s      *accrual.Service
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
		a, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(a)
	case http.StatusNoContent:
		//todo
		a, err := s.prepare(userID, resp.Body())
		if err != nil {
			log.Println(err)
		}
		s.save(a)
		return
	case http.StatusTooManyRequests:
		//todo
	case http.StatusInternalServerError:
		log.Println("error response accrual")
		return
	}
}

func (s *Service) prepare(userID int, resp []byte) (accrual.Accrual, error) {
	//todo
	resp = []byte("{\n      \"order\": \"123\",\n      \"status\": \"PROCESSED\",\n      \"accrual\": 500\n  }")
	var a Response
	dec := json.NewDecoder(bytes.NewReader(resp))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return accrual.Accrual{}, err
	}
	return accrual.Accrual{
		UserID: userID,
		Amount: a.Accrual,
		Number: a.Order,
	}, nil
}

func (s *Service) save(a accrual.Accrual) {
	if err := s.s.Add(context.Background(), a); err != nil {
		log.Println(err)
	}
}
