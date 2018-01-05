package gorilla

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/qiangli/go2/logging"
	"github.com/qiangli/go2/web"
)

type GorillaServer struct {
	web.BasicServer

	Router *mux.Router
}

var log = logging.Logger()

func (r *GorillaServer) Serve() {
	port := r.Port()

	//
	if r.Router == nil {
		r.Router = mux.NewRouter()
		r.Router.HandleFunc("/", r.home)
	}

	log.Infof("Server listening on port: %s", port)

	log.Fatal(http.ListenAndServe(":" + port, r.Router))
}

func (r *GorillaServer) home(res http.ResponseWriter, req *http.Request) {
	type message struct {
		Server    string `json:"server"`
		Name      string `json:"name"`
		Version   string `json:"version"`
		Build     string `json:"build"`
		Timestamp int64 `json:"timestamp"`
	}
	n := r.Ctx.Env.GetStringEnv("VCAP_APPLICATION", "name")
	v := r.Ctx.Env.GetStringEnv("VCAP_APPLICATION", "version")
	b := r.Ctx.Env.GetStringEnv("build")
	t := web.CurrentTimestamp()
	m := &message{Server: "gorilla", Name: n, Version: v, Build: b, Timestamp: t}

	r.HandleJson(m, res, req)
}

func NewGorillaServer(router ...*mux.Router) *GorillaServer {
	ctx := web.CreateAppContext()

	if len(router) == 0 {
		return &GorillaServer{BasicServer: web.BasicServer{Ctx: ctx}, Router: nil}
	} else {
		return &GorillaServer{BasicServer: web.BasicServer{Ctx: ctx}, Router: router[0]}
	}
}