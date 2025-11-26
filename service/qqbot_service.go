package service

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"github.com/tencent-connect/botgo/token"
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
	log.Println("callback start", "method=", r.Method, "path=", r.URL.Path, "remote=", r.RemoteAddr)
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		log.Println("content-type", ct)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("callback read body err", err)
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}
	var p Payload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Println("callback unmarshal payload err", err)
	} else {
		log.Printf("payload id=%s op=%d t=%s s=%d", p.ID, p.Op, p.EventName, p.Sequence)
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	HandleValidation(rw, r, botSecret)
	log.Println("callback end", "method=", r.Method, "path=", r.URL.Path)
}

// InitQQTokenFromEnv 使用 botgo 的 token_source 初始化并自动刷新
func InitQQTokenFromEnv() error {
	appID := "102816513"
	appSecret := botSecret

	ts := token.NewQQBotTokenSource(&token.QQBotCredentials{AppID: appID, AppSecret: appSecret})
	if err := token.StartRefreshAccessToken(context.Background(), ts); err != nil {
		log.Println("start refresh access token failed:", err)
		return err
	}
	log.Println("qq token refresh started")
	return nil
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
