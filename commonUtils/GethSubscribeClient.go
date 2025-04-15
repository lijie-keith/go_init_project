package commonUtils

func init() {

	//订阅合约事件
	//go func() {
	//	client, err := ethclient.Dial(config.WSS_CLIENT_URL)
	//	if err != nil {
	//		config.SystemLogger.Error(err.Error())
	//	}
	//	contractAddress := common.HexToAddress("0x2b57d7fFd5708A91d4F8a2Dd1242783744851796")
	//	query := ethereum.FilterQuery{
	//		Addresses: []common.Address{contractAddress},
	//	}
	//	logs := make(chan types.Log)
	//
	//	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	for {
	//		select {
	//		case err := <-sub.Err():
	//			config.SystemLogger.Error(err.Error())
	//		case vLog := <-logs:
	//			config.SystemLogger.Info(string(vLog.Data))
	//		}
	//	}
	//}()
}
