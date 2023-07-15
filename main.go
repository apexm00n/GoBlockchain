package main

import (
	"fmt"
	"goblockchain/wallet"
	"goblockchain/block"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main(){
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// Wallet
	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockChainAddress(), walletB.BlockChainAddress(), 1.0)

	//Blockchain
	blockchain := block.NewBlockhain(walletM.BlockChainAddress())
	isAdded := blockchain.AddTransaction(walletA.BlockChainAddress(), walletB.BlockChainAddress(), 1.0,
	walletA.PublicKey(), t.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("A %.1f \n", blockchain.CalculateTotalAmount(walletA.BlockChainAddress()))
	fmt.Printf("B %.1f \n", blockchain.CalculateTotalAmount(walletB.BlockChainAddress()))
	fmt.Printf("M %.1f \n", blockchain.CalculateTotalAmount(walletM.BlockChainAddress()))
}