package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce 		    int
	previousHash    [32]byte
	timestamp 	    int64
	transactions    []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print(){
	fmt.Printf("timestamp 		%d\n", b.timestamp)
	fmt.Printf("nonce 			%d\n", b.nonce)
	fmt.Printf("previousHash 	%x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Timestamp 	 int64	  		`json:"timestamp"`
		Nonce 		 int	  		`json:"nonce"`
		PreviuosHash [32]byte 		`json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nonce: b.nonce,
		PreviuosHash: b.previousHash,
		Transactions: b.transactions,
	})
}


type BlockChain struct {
	transactionPool []*Transaction
	chain 			[]*Block
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block{
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func NewBlockhain() *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *BlockChain) Print(){
	for i, block := range bc.chain{
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain) - 1]
}

func (bc *BlockChain) AddTransaction (sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}


type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value 					   float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print(){
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address    %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient_blockchain_address %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value                        %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main(){
	blockchain := NewBlockhain()
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 1.0)

	previousHash := blockchain.LastBlock().Hash()
	blockchain.CreateBlock(5, previousHash)
	blockchain.Print()
	previousHash = blockchain.LastBlock().Hash()

	blockchain.AddTransaction("C", "D", 3.0)
	blockchain.AddTransaction("X", "Y", 5.0)
	blockchain.CreateBlock(2, previousHash)
	blockchain.Print()
}