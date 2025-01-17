package gapi

import (
	"fmt"
	db "project/simplebank/db/sqlc"
	"project/simplebank/pb"
	"project/simplebank/token"
	"project/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		fmt.Printf("Key length in bytes: %d\n", len([]byte(config.TokenSymmetricKey)))
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
