package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"goblockchain/utils"
	"goblockchain/wallet"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(build.Default.GOPATH + "/src/goblockchain/wallet_server/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case http.MethodPost:
			w.Header().Add("Content-Type", "application/json")
			myWallet := wallet.NewWallet()
			m, _ := myWallet.MarshalJSON()
			io.WriteString(w, string(m[:]))
		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case http.MethodPost:
			decoder := json.NewDecoder(req.Body)
			var t *wallet.TransactionRequest
			err := decoder.Decode(&t)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}
			if !t.Validate() {
				log.Printf("ERROR: missing field(s)")
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			fmt.Println(*t.SenderPrivateKey)
			fmt.Println(*t.SenderPublicKey)
			fmt.Println(*t.SenderBlockchainAddress)
			fmt.Println(*t.RecipientBlockchainAddress)
			fmt.Println(*t.Value)

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}