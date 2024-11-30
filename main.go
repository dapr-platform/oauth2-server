package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"oauth2-server/config"
	"oauth2-server/dapr"
	_ "oauth2-server/docs"
	"oauth2-server/model"
	_ "oauth2-server/prom"
	"os"
	"strconv"
	"time"

	"github.com/dapr-platform/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/dchest/captcha"
	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "go.uber.org/automaxprocs"
)

var dumpvar bool

var LISTEN_PORT = 80

var oauthServer *server.Server
var clientStore *dapr.ClientStore
var captchaHandler http.Handler
var VERIFY_CAPTCHA = false
var REGISTER_SMS_CODE = false
var BUILD_TIME = ""

// 创建一个配置结构体
type Config struct {
	ListenPort      int
	VerifyCaptcha   bool
	RegisterSmsCode bool
	TokenExpiry     time.Duration
	RefreshExpiry   time.Duration
	JWTSecret       []byte
}

// 从环境变量加载配置
func loadConfig() *Config {
	cfg := &Config{
		ListenPort:    80, // default
		VerifyCaptcha: false,
		TokenExpiry:   time.Duration(common.USER_EXPIRED_SECONDS) * time.Second,
		RefreshExpiry: time.Duration(common.USER_EXPIRED_SECONDS) * 3 * time.Second,
		JWTSecret:     []byte("00000000"),
	}

	if port := os.Getenv("LISTEN_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.ListenPort = p
		}
	}

	if verify := os.Getenv("VERIFY_CAPTCHA"); verify == "true" {
		cfg.VerifyCaptcha = true
	}
	if registerSmsCode := os.Getenv("REGISTER_SMS_CODE"); registerSmsCode == "true" {
		cfg.RegisterSmsCode = true
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWTSecret = []byte(secret)
	}

	return cfg
}

func setupOAuthServer(cfg *Config) (*server.Server, error) {
	manager := manage.NewDefaultManager()

	// 统一的token配置
	tokenConfig := &manage.Config{
		AccessTokenExp:    cfg.TokenExpiry,
		RefreshTokenExp:   cfg.RefreshExpiry,
		IsGenerateRefresh: true,
	}

	manager.SetAuthorizeCodeTokenCfg(tokenConfig)
	manager.SetPasswordTokenCfg(tokenConfig)

	manager.MustTokenStorage(dapr.NewDaprTokenStore())

	// JWT配置
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", cfg.JWTSecret, jwt.SigningMethodHS512))

	// Client存储
	var err error
	clientStore, err = dapr.NewClientStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create client store: %w", err)
	}
	manager.MapClientStorage(clientStore)

	// 创建服务器
	srv := server.NewServer(server.NewConfig(), manager)

	// 配置handlers
	setupOAuthHandlers(srv)

	return srv, nil
}

func setupOAuthHandlers(srv *server.Server) {
	srv.SetPasswordAuthorizationHandler(passwordAuthHandler)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		common.Logger.Error("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		common.Logger.Error("Response Error:", re)
	})
}

func setupRoutes(mux *chi.Mux) {
	if mux == nil {
		panic("mux cannot be nil")
	}
	captchaHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
	// OAuth routes
	mux.HandleFunc("/oauth/authorize", authorizeHandler)
	mux.HandleFunc("/oauth/token", tokenHandler)
	mux.HandleFunc("/oauth/token-by-field", tokenByFieldHandler)

	// Auth routes
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/auth", authHandler)
	mux.HandleFunc("/logout", logoutHandler)

	// Captcha routes
	mux.HandleFunc("/captcha-gen", captchaGen)
	mux.HandleFunc("/captcha/{file}", captchaHandler.ServeHTTP)

	// Test routes
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/test_use_post", testUsePostHandler)

	// Admin routes
	mux.HandleFunc("/initTestClientStore", initTestClientStoreHandler)
	mux.HandleFunc("/refreshClientStoreFromDb", refreshClientFromDbHandler)

	// User routes
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/sms-code/send", smsCodeSendHandler)

	// Swagger
	mux.Handle("/swagger*", httpSwagger.WrapHandler)
}

// @title oauth2-server RESTful API
// @version 1.0
// @description oauth2-server  RESTful API 文档.
// @BasePath /swagger/oauth2-server
func main() {
	// 加载配置
	cfg := loadConfig()

	// 初始化OAuth服务器
	var err error
	oauthServer, err = setupOAuthServer(cfg)
	if err != nil {
		log.Fatalf("Failed to setup OAuth server: %v", err)
	}

	// 设置路由
	mux := chi.NewRouter()
	setupRoutes(mux)

	// 启动服务器
	s := daprd.NewServiceWithMux(":"+strconv.Itoa(cfg.ListenPort), mux)
	if s == nil {
		log.Fatal("Failed to create server")
	}

	common.Logger.Info("server starting on port", cfg.ListenPort)

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func refreshClientInfoFromDb(clientStore *dapr.ClientStore) error {
	if clientStore == nil {
		return fmt.Errorf("clientStore is nil")
	}

	data, err := common.GetDaprClient().InvokeMethod(context.Background(), "db-service", "/"+config.DBNAME+"/public/"+config.CLIENT_INFO_TABLE_NAME, "get")
	if err != nil {
		common.Logger.Error("refreshClientInfoFromDb error:", err.Error())
		return fmt.Errorf("failed to invoke db-service: %w", err)
	}

	var infolist []model.Client_info
	if err = json.Unmarshal(data, &infolist); err != nil {
		common.Logger.Error("refreshClientInfoFromDb Unmarshal data error:", err.Error())
		return fmt.Errorf("failed to unmarshal client info: %w", err)
	}

	for _, info := range infolist {
		if err = clientStore.Set(info.ID, &models.Client{
			ID:     info.ID,
			Secret: info.Password,
		}); err != nil {
			common.Logger.Error("clientStore.Set error:", err.Error())
			return fmt.Errorf("failed to set client info: %w", err)
		}
	}
	return nil
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		if err := dumpRequest(os.Stdout, "userAuthorizeHandler", r); err != nil {
			common.Logger.Error("Failed to dump request:", err)
		}
	}

	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return "", fmt.Errorf("failed to start session: %w", err)
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				return "", fmt.Errorf("failed to parse form: %w", err)
			}
		}

		store.Set("ReturnUri", r.Form)
		if err := store.Save(); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	if err := store.Save(); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}
	return userID, nil
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	if writer == nil || r == nil {
		return fmt.Errorf("invalid parameters")
	}

	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return fmt.Errorf("failed to dump request: %w", err)
	}

	if _, err := writer.Write([]byte("\n" + header + ": \n")); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	return nil
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	if w == nil || req == nil || filename == "" {
		http.Error(w, "Invalid parameters", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
