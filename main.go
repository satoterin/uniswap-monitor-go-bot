package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/pradeep-selva/uniswap-monitor/UniswapUSDC2Pool"
)

func GetEnvVar(key string) string {
	return os.Getenv(key)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("An error occured", err)
	}

	client, err := ethclient.Dial(GetEnvVar("INFURA_ENDPOINT"))
	checkError(err)

	contractAddress := common.HexToAddress(GetEnvVar("UNISWAP_USDC2_ADDRESS"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	checkError(err)

	log.Println("Listening for Uniswap USDC2 pool's events...")

	uniswapAbi, err := abi.JSON(strings.NewReader(string(UniswapUSDC2Pool.UniswapUSDC2PoolABI)))
	checkError(err)

	for {
		select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs:
				// event := struct {
				// 	Sender [32]byte
				// 	Recipient [32]byte
				// 	Amount0 [32]byte
				// 	Amount1 [32]byte
				// 	SqrtPriceX96 [32]byte
				// 	Liquidity [32]byte
				// 	Tick [32]byte
				// }{}

				log.Println("--------- NEW SWAP ---------")
				swapEvent, err := uniswapAbi.Unpack("Swap", vLog.Data)
				checkError(err)

				log.Println(swapEvent...)
		}
	}
}
