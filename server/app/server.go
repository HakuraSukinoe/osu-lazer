package app

import (
	"github.com/deissh/osu-lazer/server/mlog"
	"github.com/deissh/osu-lazer/server/store"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

type Server struct {
	Store store.Store

	Server     *http.Server
	ListenAddr *net.TCPAddr

	didFinishListen chan struct{}

	goroutineCount      int32
	goroutineExitSignal chan struct{}

	clusterLeaderListeners sync.Map

	licenseValue       atomic.Value
	clientLicenseValue atomic.Value
	licenseListeners   map[string]func()

	newStore func() store.Store

	clientConfig        map[string]string
	clientConfigHash    string
	limitedClientConfig map[string]string

	Log              *mlog.Logger
	NotificationsLog *mlog.Logger

	joinCluster  bool
	startMetrics bool
}
