package service

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// Payload 网关事件通用包装
type Payload struct {
	ID        string          `json:"id"`
	Op        int             `json:"op"`
	Data      json.RawMessage `json:"d"`
	Sequence  int64           `json:"s"`
	EventName string          `json:"t"`
}

// ValidationRequest 验签请求体
type ValidationRequest struct {
	PlainToken string `json:"plain_token"`
	EventTs    string `json:"event_ts"`
}

// ValidationResponse 验签响应体
type ValidationResponse struct {
	PlainToken string `json:"plain_token"`
	Signature  string `json:"signature"`
}

var botSecret = "ayMk8XwLk9YxNnDd3TtKlCd4VwNpHjBd"

func CallbackMSg(rw http.ResponseWriter, r *http.Request) {
	HandleValidation(rw, r, botSecret)
}

func HandleValidation(rw http.ResponseWriter, r *http.Request, botSecret string) {
	httpBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read http body err", err)
		return
	}
	payload := &Payload{}
	if err = json.Unmarshal(httpBody, payload); err != nil {
		log.Println("parse http payload err", err)
		return
	}
	validationPayload := &ValidationRequest{}
	if err = json.Unmarshal(payload.Data, validationPayload); err != nil {
		log.Println("parse http payload failed:", err)
		return
	}
	seed := botSecret
	for len(seed) < ed25519.SeedSize {
		seed = strings.Repeat(seed, 2)
	}
	seed = seed[:ed25519.SeedSize]
	reader := strings.NewReader(seed)
	// GenerateKey 方法会返回公钥、私钥，这里只需要私钥进行签名生成不需要返回公钥
	_, privateKey, err := ed25519.GenerateKey(reader)
	if err != nil {
		log.Println("ed25519 generate key failed:", err)
		return
	}
	var msg bytes.Buffer
	msg.WriteString(validationPayload.EventTs)
	msg.WriteString(validationPayload.PlainToken)
	signature := hex.EncodeToString(ed25519.Sign(privateKey, msg.Bytes()))
	if err != nil {
		log.Println("generate signature failed:", err)
		return
	}
	rspBytes, err := json.Marshal(
		&ValidationResponse{
			PlainToken: validationPayload.PlainToken,
			Signature:  signature,
		})
	if err != nil {
		log.Println("handle validation failed:", err)
		return
	}
	rw.Write(rspBytes)
}
