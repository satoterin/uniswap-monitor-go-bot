package main

import (
	"context"
	"log"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
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
	
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case event := <-logs:
			log.Println(event)
		}
	}
}
