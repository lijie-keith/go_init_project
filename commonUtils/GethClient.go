package commonUtils

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lijie-keith/go_init_project/config"
)

var Client *ethclient.Client

func init() {
	client, err := ethclient.Dial(config.CLIENT_URL)

	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	Client = client
}
