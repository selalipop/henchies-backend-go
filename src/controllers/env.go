package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
)

type Controllers struct {
	 PlayerRepository repository.PlayerRepository
	 GameRepository repository.GameRepository
}
