package server

import (
	"encoding/json"
	"net/http"
	"sync"
)

const MigrationsEndpoint = "/migrations"

type versionFunc struct {
	lock          sync.Mutex
	unsafeVersion func() (int64, int64, error)
}

func (v *versionFunc) Set(version func() (int64, int64, error)) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.unsafeVersion = version
}

func (v *versionFunc) Get() func() (int64, int64, error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.unsafeVersion
}

type MigrationsInfoEndpoint struct {
	version versionFunc
}

func NewMigrationsInfoEndpoint() *MigrationsInfoEndpoint {
	return &MigrationsInfoEndpoint{
		version: versionFunc{
			lock: sync.Mutex{},
			unsafeVersion: func() (int64, int64, error) {
				return 0, 0, nil
			},
		},
	}
}

func (e *MigrationsInfoEndpoint) SetVersion(version func() (int64, int64, error)) {
	e.version.Set(version)
}

func (e *MigrationsInfoEndpoint) ApplyHandlers(mux *http.ServeMux) error {
	mux.Handle(MigrationsEndpoint, http.HandlerFunc(e.infoEndpoint))

	return nil
}

func (e *MigrationsInfoEndpoint) infoEndpoint(w http.ResponseWriter, _ *http.Request) {
	current, expected, err := e.version.Get()()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info := MigrationsInfo{
		Version:           current,
		MigrationsApplied: current == expected,
	}
	rsp, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
	w.Write([]byte("\n"))
}

type MigrationsInfo struct {
	Version           int64
	MigrationsApplied bool
}
