package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Aoluis1005/go-farm-bot/models"
)

// registerUserAPI 注册多用户系统API路由
func registerUserAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/users/login", handleUserLogin)
	mux.HandleFunc("/api/users/register", handleUserRegister)
	mux.HandleFunc("/api/users/claim-card", handleUserClaimCard)
	mux.HandleFunc("/api/users/renew", handleUserRenew)
	mux.HandleFunc("/api/users/me", handleUserMe)
	mux.HandleFunc("/api/users/delete", handleUserDelete)
	mux.HandleFunc("/api/cards", handleCards)
}

// ---- 用户管理 ----

// handleUsers: 用户列表管理（GET: 列出用户, POST: 创建用户）
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := models.GetAllUsers()
		writeJSON(w, map[string]interface{}{"ok": true, "data": users})
	case http.MethodPost:
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
			CardCode string `json:"cardCode,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "参数错误")
			return
		}
		if body.Username == "" || body.Password == "" {
			writeError(w, http.StatusBadRequest, "用户名和密码不能为空")
			return
		}
		result := models.RegisterUser(body.Username, body.Password, body.CardCode)
		if !result.Ok {
			writeError(w, http.StatusBadRequest, result.Error)
			return
		}
		writeJSON(w, map[string]interface{}{"ok": true, "user": result.User})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleUserLogin: 用户登录，返回token
func handleUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "参数错误")
		return
	}
	if body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	result, err := models.ValidateUser(body.Username, body.Password, clientIP)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 生成登录 token
	token, err := models.NewUserToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "生成token失败")
		return
	}

	writeJSON(w, map[string]interface{}{
		"ok": true,
		"token": token,
		"user": map[string]interface{}{
			"username":     result.Username,
			"role":         result.Role,
			"accountLimit": result.AccountLimit,
		},
	})
}

// handleUserRegister: 用户注册接口
func handleUserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		CardCode string `json:"cardCode,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "参数错误")
		return
	}
	if body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}

	result := models.RegisterUser(body.Username, body.Password, body.CardCode)
	if !result.Ok {
		writeError(w, http.StatusBadRequest, result.Error)
		return
	}

	writeJSON(w, map[string]interface{}{
		"ok":  true,
		"user": result.User,
	})
}

// handleUserMe: 获取当前用户信息
func handleUserMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		writeError(w, http.StatusBadRequest, "缺少用户名")
		return
	}

	users := models.GetAllUsers()
	for _, u := range users {
		if u.Username == username {
			writeJSON(w, map[string]interface{}{
				"ok":  true,
				"user": u,
			})
			return
		}
	}
	writeError(w, http.StatusNotFound, "用户不存在")
}

// handleUserRenew: 用户续期（使用卡密）
func handleUserRenew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Username string `json:"username"`
		CardCode string `json:"cardCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "参数错误")
		return
	}

	success, msg, err := models.RenewUser(body.Username, body.CardCode)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, map[string]interface{}{
		"ok":    success,
		"message": msg,
	})
}

// handleUserDelete: 删除用户
func handleUserDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "参数错误")
		return
	}

	err := models.DeleteUser(body.Username)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// ---- 卡密管理 ----

// handleCards: 卡密列表
func handleCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var cards []models.Card
	users := models.GetAllUsers()
	for _, u := range users {
		if u.Card != nil {
			cards = append(cards, *u.Card)
		}
	}
	writeJSON(w, map[string]interface{}{"ok": true, "data": cards})
}

// handleUserClaimCard: 用户领取卡密（防刷机制）
func handleUserClaimCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		Username string `json:"username"`
		UserAgent string `json:"userAgent,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "参数错误")
		return
	}

	ua := body.UserAgent
	if ua == "" {
		ua = "default"
	}

	result := models.ClaimCardByUA(ua, body.Username)
	writeJSON(w, map[string]interface{}{
		"ok":       result.Ok,
		"cardCode": result.CardCode,
		"days":     result.Days,
		"durationMs": result.DurationMs,
		"error":  result.Error,
	})
}