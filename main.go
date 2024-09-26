package main

import (
	"context"
	"encoding/json"
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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"oauth2-server/config"
	"oauth2-server/dapr"
	_ "oauth2-server/docs"
	"oauth2-server/model"
	"oauth2-server/mycaptcha"
	_ "oauth2-server/prom"
	"os"
	"strconv"
	"strings"
	"time"
)

var dumpvar bool

var LISTEN_PORT = 80

var oauthServer *server.Server
var clientStore *dapr.ClientStore
var captchaHandler http.Handler
var VERIFY_CAPTCHA = false

func init() {
	buf, _ := ioutil.ReadFile("build.time")

	t, _ := strconv.ParseInt(strings.TrimSpace(string(buf))+"000", 10, 64)
	s := time.UnixMilli(t).Format("2006-01-02 15:04:05")
	common.Logger.Info("build.time:", s, string(buf))

	if val := os.Getenv("LISTEN_PORT"); val != "" {
		LISTEN_PORT, _ = strconv.Atoi(val)
	}
	if val := os.Getenv("LISTEN_PORT"); val != "" {
		LISTEN_PORT, _ = strconv.Atoi(val)
	}

	if val := os.Getenv("VERIFY_CAPTCHA"); val != "" {
		VERIFY_CAPTCHA = val == "true"
	}
	captcha.SetCustomStore(&mycaptcha.DaprCaptchaStore{
		Expiration: time.Minute,
	})

	captchaHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
}

// @title oauth2-server RESTful API
// @version 1.0
// @description oauth2-server  RESTful API 文档.
// @BasePath /swagger/oauth2-server
func main() {

	manager := manage.NewDefaultManager()
	config := &manage.Config{AccessTokenExp: time.Second * time.Duration(common.USER_EXPIRED_SECONDS), RefreshTokenExp: time.Second * time.Duration(common.USER_EXPIRED_SECONDS) * 3, IsGenerateRefresh: true}
	manager.SetAuthorizeCodeTokenCfg(config)
	// token store
	ptcfg := &manage.Config{AccessTokenExp: time.Second * time.Duration(common.USER_EXPIRED_SECONDS), RefreshTokenExp: time.Second * time.Duration(common.USER_EXPIRED_SECONDS) * 3, IsGenerateRefresh: true}
	//common.Logger.Info("token expired:time.Second * 30")
	manager.SetPasswordTokenCfg(ptcfg)
	manager.MustTokenStorage(dapr.NewDaprTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	//manager.MapAccessGenerate(generates.NewAccessGenerate())
	clientStore, err := dapr.NewClientStore()
	if err != nil {
		panic(err)
	}

	manager.MapClientStorage(clientStore)
	oauthServer = server.NewServer(server.NewConfig(), manager)
	oauthServer.SetPasswordAuthorizationHandler(passwordAuthHandler)

	oauthServer.SetClientInfoHandler(server.ClientFormHandler)
	oauthServer.SetUserAuthorizationHandler(userAuthorizeHandler)

	oauthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		common.Logger.Error("Internal Error:", err.Error())
		return
	})

	oauthServer.SetResponseErrorHandler(func(re *errors.Response) {

		common.Logger.Error("Response Error:", re)
	})
	mux := chi.NewRouter()
	//api.InitUserRoute(mux)
	//api.InitClient_infoRoute(mux)
	mux.HandleFunc("/captcha-gen", captchaGen)
	mux.HandleFunc("/captcha/{file}", captchaHandler.ServeHTTP)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/auth", authHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.Handle("/swagger*", httpSwagger.WrapHandler)
	mux.HandleFunc("/oauth/authorize", authorizeHandler)

	mux.HandleFunc("/oauth/token", tokenHandler)
	mux.HandleFunc("/oauth/token-by-field", tokenByFieldHandler)
	mux.HandleFunc("/test", testHandler)

	mux.HandleFunc("/test_use_post", testUsePostHandler)

	mux.HandleFunc("/initTestClientStore", initTestClientStoreHandler)
	mux.HandleFunc("/refreshClientStoreFromDb", refreshClientFromDbHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	s := daprd.NewServiceWithMux(":"+strconv.Itoa(LISTEN_PORT), mux)
	common.Logger.Info("server start")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
	}

}

func refreshClientInfoFromDb(clientStore *dapr.ClientStore) error {
	data, err := common.GetDaprClient().InvokeMethod(context.Background(), "db-service", "/"+config.DBNAME+"/public/"+config.CLIENT_INFO_TABLE_NAME, "get")
	if err != nil {
		log.Printf("refreshClientInfoFromDb error.%s", err.Error())
		return err
	}
	var infolist []model.Client_info
	err = json.Unmarshal(data, &infolist)
	if err != nil {
		log.Printf("refreshClientInfoFromDb Unmarshal data error.%s", err.Error())
		return err
	}
	for _, info := range infolist {
		err = clientStore.Set(info.ID, &models.Client{
			ID:     info.ID,
			Secret: info.Password,
		})
		if err != nil {
			log.Printf("clientStore.Set error.%s", err.Error())
			return err
		}
	}
	return nil

}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
