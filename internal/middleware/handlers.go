package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go_project/db"
	"go_project/docs"
	_ "go_project/docs"
	"go_project/internal"
	"go_project/internal/entity"
	"go_project/internal/persistance"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type App struct {
	Router           *mux.Router
	DB               *db.MongoConnection
	Logg             log.Logger
	PersonRepository persistance.PersonRepository
}

func (a *App) Run(addr string) error {
	err := http.ListenAndServe(addr, a.Router)
	return err
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv(internal.ORIGIN_ALLOWED))
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func (a *App) Initialize(_user, _password string) (err error) {
	fmt.Println(internal.MsgResponseStartApplication)
	host := os.Getenv(internal.MONGO_HOST) + ":" + os.Getenv(internal.MONGO_PORT)
	dbs := os.Getenv(internal.MONGO_DATABASE)
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Hour,
		Database: dbs,
		Username: _user,
		Password: _password,
	}
	a.DB, _ = db.NewConnection(info)
	fmt.Println(internal.MsgResponseConnectedMongoDB)
	muxObj := mux.NewRouter()
	muxObj.Use(CORS)
	a.Router = muxObj
	values := []interface{}{internal.KeyType, internal.SUCCESS, internal.KeyURL, internal.URLStartingNow, internal.KeyMessage, internal.MsgResponseStartingNow}
	a.LoggingOperation(values...)
	a.initializeRoutes()
	a.initializeSwagger()
	a.initializeRepository()
	return err
}

// routing
func (a *App) initializeRoutes() {
	a.Router.PathPrefix(internal.URLApi).Handler(httpSwagger.WrapHandler)
	a.Router.HandleFunc(internal.URLIndex, a.getIndex).Methods(internal.HTTP_GET)
	a.Router.HandleFunc(internal.URLCreatingOne, a.CreatePerson).Methods(internal.HTTP_POST)
	a.Router.HandleFunc(internal.URLUpdatingOne, a.UpdatePerson).Methods(internal.HTTP_POST)
	a.Router.HandleFunc(internal.URLGettingOne, a.GetPerson).Methods(internal.HTTP_POST)
	a.Router.HandleFunc(internal.URLDeletingOne, a.DeletePerson).Methods(internal.HTTP_POST)
	a.Router.HandleFunc(internal.URLListingAll, a.ListPersons).Methods(internal.HTTP_POST)
}

// swagger
func (a *App) initializeSwagger() {
	docs.SwaggerInfo.Title = internal.MsgApiRestTitle
	docs.SwaggerInfo.Description = internal.MsgApiRestDescription
	docs.SwaggerInfo.Version = internal.MsgApiRestVersion1
	docs.SwaggerInfo.Host = internal.URLLocalhost + ":" + os.Getenv(internal.APP_PORT)
	docs.SwaggerInfo.BasePath = internal.URLStartingNow
	docs.SwaggerInfo.Schemes = []string{internal.SchemaHttp}
}

// Logger
func (a *App) initializeLogger() (f *os.File) {
	f, _ = os.OpenFile(a.RootDir()+"/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	wrt := io.MultiWriter(os.Stdout, f)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(wrt)
		logger = log.With(logger, internal.KeyTs, log.DefaultTimestampUTC)
		logger = log.With(logger, internal.KeyCaller, log.DefaultCaller)
	}
	a.Logg = logger
	return f
}

func (a *App) LoggingOperation(values ...interface{}) {
	f := *a.initializeLogger()
	_ = a.Logg.Log(values...)
	defer f.Close()
}

// Repository
func (a *App) initializeRepository() {
	a.PersonRepository = persistance.NewPersonRepository(a.DB)
}

// Init
func (a *App) getIndex(w http.ResponseWriter, r *http.Request) {

	item := &entity.JsonResponse{
		Message:    internal.MsgApiRestTitle,
		StatusCode: http.StatusOK,
	}
	respondWithJSON(w, http.StatusOK, item)
}

//ROOT DIR
func (a *App) RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
	//d := path.Join(path.Dir(b))
	//return filepath.Dir(d)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
