package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"goblockchain/block"
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

			publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
			privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
			value, err := strconv.ParseFloat(*t.Value, 32)
			if err != nil {
				log.Println("ERROR: parse error")
			}
			value32 := float32(value)		

			w.Header().Add("Content-type", "application/json")
			transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
			signature := transaction.GenerateSignature()
			signatureStr := signature.String()

			bt := &block.TransactionRequest{
				t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.SenderPublicKey,
				&value32, &signatureStr,
			}
			
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			resp, err := http.Post(ws.Gateway() + "/transactions", "application/json", buf)
			if err != nil {
				fmt.Println("ERROR:" + err.Error())
			}
			if resp.StatusCode == 201 {
				io.WriteString(w, string(utils.JsonStatus("success")))
				return
			}
			io.WriteString(w, string(utils.JsonStatus("fail")))

		default:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount (w http.ResponseWriter, req *http.Request) {
	switch req.Method{
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())
		client := &http.Client{}
		bcsReq, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Printf("HTTP ERROR: " + err.Error())
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar block.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}
				m, _ := json.Marshal(struct{
					Message string `json:"message"`
					Amount  float32 `json:"amount"`
				}{
					Message: "success",
					Amount: bar.Amount,
				})
				io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}