// Auto-generated by "nex regen" at "2017-11-13 20:58:30.057437 +0800 CST"
// ** DO NOT EDIT **

package player

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/damnever/cc"
	huskarpool "github.com/eleme/huskar-pool"
	resourcepool "github.com/eleme/huskar-pool/pool"
	"github.com/eleme/huskar/config"
	"github.com/eleme/nex"
	"github.com/eleme/nex/circuitbreaker"
	"github.com/eleme/nex/consts"
	"github.com/eleme/nex/consts/ctxkeys"
	"github.com/eleme/nex/endpoint"
	"github.com/eleme/nex/log"
	"github.com/eleme/nex/metric"
	"github.com/eleme/nex/timeout"
	"github.com/eleme/nex/tracking/etrace"
	ttracker "github.com/eleme/nex/tracking/thrift"
	json "github.com/json-iterator/go"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/eleme/purchaseMeiTuan/services/player"
)

// needed to ensure import safety.
var _ = metric.GoUnusedProtection
var _ = player.GoUnusedProtection__

var (
	thriftplayerClient               *ThriftplayerServiceClient
	thriftplayerClientFromHuskar     *ThriftplayerServiceClient
	thriftplayerClientOnce           sync.Once
	thriftplayerClientFromHuskarOnce sync.Once
)

// CliAppErrTypes contains thrift user/system/unknown exception types(reflect).
var CliAppErrTypes = &endpoint.ErrTypes{
	UserErr:  reflect.TypeOf(player.NewplayerUserException()),
	SysErr:   reflect.TypeOf(player.NewplayerSystemException()),
	UnkwnErr: reflect.TypeOf(player.NewplayerUnknownException()),
}

// IsTolerableException checks the error whether is tolerable.
func IsTolerableException(err error) bool {
	if err == nil {
		return true
	}
	switch err.(type) {
	case *player.playerUserException:
		return true
	case *player.playerSystemException:
		return true
	case *player.playerUnknownException:
		return true
	case thrift.TApplicationException:
		return true
	default:
		return false
	}
}

// GetThriftplayerServiceClient returns a pooled client, it use Addr which defines in thriftfs/deps.json.
func GetThriftplayerServiceClient() (*ThriftplayerServiceClient, error) {
	var err error
	thriftplayerClientOnce.Do(func() {
		thriftplayerClient, err = makeThriftplayerServiceClient(false)
	})
	return thriftplayerClient, err
}

// GetThriftplayerServiceClientFromHuskar returns a pooled client, it use the address fetch from Huskar.
func GetThriftplayerServiceClientFromHuskar() (*ThriftplayerServiceClient, error) {
	var err error
	thriftplayerClientFromHuskarOnce.Do(func() {
		thriftplayerClientFromHuskar, err = makeThriftplayerServiceClient(true)
	})
	return thriftplayerClientFromHuskar, err
}

// ThriftplayerServiceClient is wrapper with middlewares, context and pool for thrift's playerServiceClient.
type ThriftplayerServiceClient struct {
	pool                    resourcepool.Pooler
	assigntanksEndpoint     endpoint.Endpoint
	getnewordersEndpoint    endpoint.Endpoint
	lateststateEndpoint     endpoint.Endpoint
	pingEndpoint            endpoint.Endpoint
	uploadmapEndpoint       endpoint.Endpoint
	uploadparamtersEndpoint endpoint.Endpoint
}

// ThriftplayerServiceClientOptions defines optional arguments for NewThriftplayerServiceClient,
// the Addr and AppName is required, IdleTimeout, MaxActive and MaxCap is used for connection pool,
// if Logger is nil, logging is disabled, if Trace is nil, etrace is disabled, if CircuitBreaker
// is nil, circuitbreaker is disabled, if HuskarConfiger is nil, timeout is disabled.
type ThriftplayerServiceClientOptions struct {
	Addr           string
	AppName        string
	ConnectTimeout time.Duration
	RWTimeout      time.Duration
	IdleTimeout    time.Duration
	MaxActive      int
	MaxCap         int
	Logger         log.RPCContextLogger
	Trace          *etrace.Trace
	CircuitBreaker *circuitbreaker.CircuitBreaker
	HuskarConfiger config.Configer
	NexConfig      cc.Configer
}

// NewThriftplayerServiceClient creates a new ThriftplayerServiceClient.
func NewThriftplayerServiceClient(options ThriftplayerServiceClientOptions) *ThriftplayerServiceClient {
	clientFactory := func() (resourcepool.Resource, error) {
		return newRawThriftClientWrapper(options.AppName, options.Addr, options.ConnectTimeout, options.RWTimeout)
	}
	pool := resourcepool.NewResourcePool(clientFactory, options.MaxActive, options.MaxCap, options.IdleTimeout)
	return newThriftplayerServiceClient(pool, options)
}

// NewThriftplayerServiceClientWithHuskar creates a new ThriftplayerServiceClient with huskar pool support.
func NewThriftplayerServiceClientWithHuskar(huskarPool *huskarpool.Huskar, options ThriftplayerServiceClientOptions) (*ThriftplayerServiceClient, error) {
	clientFactory := func(_ huskarpool.Meta, conn net.Conn) (resourcepool.Resource, error) {
		socket := thrift.NewTSocketFromConnTimeout(conn, options.ConnectTimeout)
		return newRawThriftClientWrapperByConn(socket, options.AppName, options.RWTimeout)
	}
	poolOptions := huskarpool.PoolOption{
		MaxCap:      options.MaxCap,
		Capacity:    options.MaxActive,
		IdleTimeout: options.IdleTimeout,
		DialTimeout: options.ConnectTimeout,
	}
	cluster := nex.BuildSOAClusterOrIntent("player", options.HuskarConfiger, options.NexConfig)
	pool, err := huskarPool.NewResourcePool("purchaseMeiTuan", cluster, clientFactory, poolOptions)
	if err != nil {
		return nil, err
	}
	return newThriftplayerServiceClient(pool, options), nil
}

func newThriftplayerServiceClient(pool resourcepool.Pooler, options ThriftplayerServiceClientOptions) *ThriftplayerServiceClient {
	soaArgs := &endpoint.SOAMiddlewareArgs{
		AppID:             options.AppName,
		ThriftServiceName: "playerService",
		RemoteAddr:        options.Addr,
		ErrTypes:          CliAppErrTypes,
	}

	var assigntanksEndpoint endpoint.Endpoint
	{
		assigntanksEndpoint = makeThriftAssignTanksEndpoint()
		// NOTE: Add middlewares here.
		if options.HuskarConfiger != nil {
			assigntanksEndpoint = timeout.EndpointTimeoutSOAClientMiddleware(options.HuskarConfiger)(assigntanksEndpoint)
		}
		if options.CircuitBreaker != nil {
			assigntanksEndpoint = circuitbreaker.EndpointCircuitBreakerSOAClientMiddleware(options.CircuitBreaker)(assigntanksEndpoint)
		}
		if options.Trace != nil {
			assigntanksEndpoint = etrace.EndpointEtraceSOAClientMiddleware(options.Trace, soaArgs)(assigntanksEndpoint)
		}
		assigntanksEndpoint = metric.EndpointStatsdSOAClientMiddleware(soaArgs)(assigntanksEndpoint)
		if options.Logger != nil {
			assigntanksEndpoint = log.EndpointLoggingSOAClientMiddleware(options.Logger, soaArgs)(assigntanksEndpoint)
		}
	}

	var getnewordersEndpoint endpoint.Endpoint
	{
		getnewordersEndpoint = makeThriftGetNewOrdersEndpoint()
		// NOTE: Add middlewares here.
		if options.HuskarConfiger != nil {
			getnewordersEndpoint = timeout.EndpointTimeoutSOAClientMiddleware(options.HuskarConfiger)(getnewordersEndpoint)
		}
		if options.CircuitBreaker != nil {
			getnewordersEndpoint = circuitbreaker.EndpointCircuitBreakerSOAClientMiddleware(options.CircuitBreaker)(getnewordersEndpoint)
		}
		if options.Trace != nil {
			getnewordersEndpoint = etrace.EndpointEtraceSOAClientMiddleware(options.Trace, soaArgs)(getnewordersEndpoint)
		}
		getnewordersEndpoint = metric.EndpointStatsdSOAClientMiddleware(soaArgs)(getnewordersEndpoint)
		if options.Logger != nil {
			getnewordersEndpoint = log.EndpointLoggingSOAClientMiddleware(options.Logger, soaArgs)(getnewordersEndpoint)
		}
	}

	var lateststateEndpoint endpoint.Endpoint
	{
		lateststateEndpoint = makeThriftLatestStateEndpoint()
		// NOTE: Add middlewares here.
		if options.HuskarConfiger != nil {
			lateststateEndpoint = timeout.EndpointTimeoutSOAClientMiddleware(options.HuskarConfiger)(lateststateEndpoint)
		}
		if options.CircuitBreaker != nil {
			lateststateEndpoint = circuitbreaker.EndpointCircuitBreakerSOAClientMiddleware(options.CircuitBreaker)(lateststateEndpoint)
		}
		if options.Trace != nil {
			lateststateEndpoint = etrace.EndpointEtraceSOAClientMiddleware(options.Trace, soaArgs)(lateststateEndpoint)
		}
		lateststateEndpoint = metric.EndpointStatsdSOAClientMiddleware(soaArgs)(lateststateEndpoint)
		if options.Logger != nil {
			lateststateEndpoint = log.EndpointLoggingSOAClientMiddleware(options.Logger, soaArgs)(lateststateEndpoint)
		}
	}

	var pingEndpoint endpoint.Endpoint
	{
		pingEndpoint = makeThriftPingEndpoint()
	}

	var uploadmapEndpoint endpoint.Endpoint
	{
		uploadmapEndpoint = makeThriftUploadMapEndpoint()
		// NOTE: Add middlewares here.
		if options.HuskarConfiger != nil {
			uploadmapEndpoint = timeout.EndpointTimeoutSOAClientMiddleware(options.HuskarConfiger)(uploadmapEndpoint)
		}
		if options.CircuitBreaker != nil {
			uploadmapEndpoint = circuitbreaker.EndpointCircuitBreakerSOAClientMiddleware(options.CircuitBreaker)(uploadmapEndpoint)
		}
		if options.Trace != nil {
			uploadmapEndpoint = etrace.EndpointEtraceSOAClientMiddleware(options.Trace, soaArgs)(uploadmapEndpoint)
		}
		uploadmapEndpoint = metric.EndpointStatsdSOAClientMiddleware(soaArgs)(uploadmapEndpoint)
		if options.Logger != nil {
			uploadmapEndpoint = log.EndpointLoggingSOAClientMiddleware(options.Logger, soaArgs)(uploadmapEndpoint)
		}
	}

	var uploadparamtersEndpoint endpoint.Endpoint
	{
		uploadparamtersEndpoint = makeThriftUploadParamtersEndpoint()
		// NOTE: Add middlewares here.
		if options.HuskarConfiger != nil {
			uploadparamtersEndpoint = timeout.EndpointTimeoutSOAClientMiddleware(options.HuskarConfiger)(uploadparamtersEndpoint)
		}
		if options.CircuitBreaker != nil {
			uploadparamtersEndpoint = circuitbreaker.EndpointCircuitBreakerSOAClientMiddleware(options.CircuitBreaker)(uploadparamtersEndpoint)
		}
		if options.Trace != nil {
			uploadparamtersEndpoint = etrace.EndpointEtraceSOAClientMiddleware(options.Trace, soaArgs)(uploadparamtersEndpoint)
		}
		uploadparamtersEndpoint = metric.EndpointStatsdSOAClientMiddleware(soaArgs)(uploadparamtersEndpoint)
		if options.Logger != nil {
			uploadparamtersEndpoint = log.EndpointLoggingSOAClientMiddleware(options.Logger, soaArgs)(uploadparamtersEndpoint)
		}
	}

	return &ThriftplayerServiceClient{
		pool:                    pool,
		assigntanksEndpoint:     assigntanksEndpoint,
		getnewordersEndpoint:    getnewordersEndpoint,
		lateststateEndpoint:     lateststateEndpoint,
		pingEndpoint:            pingEndpoint,
		uploadmapEndpoint:       uploadmapEndpoint,
		uploadparamtersEndpoint: uploadparamtersEndpoint,
	}
}

func makeThriftplayerServiceClient(fromHuskar bool) (*ThriftplayerServiceClient, error) {
	logger, err := log.GetContextLogger("soa_client.player")
	if err != nil {
		return nil, err
	}

	svcCfg, err := findConfigFromRPCDepsFile("player")
	if err != nil {
		return nil, err
	}
	nexConfig := nex.GetNexConfig()
	poolOpts := svcCfg.Config("PoolOptions")
	maxActive := poolOpts.IntOr("MaxActive", 60)
	options := ThriftplayerServiceClientOptions{
		AppName:        nexConfig.String("app_name"),
		ConnectTimeout: poolOpts.DurationOr("ConnectTimeout", 3000) * time.Millisecond,
		RWTimeout:      poolOpts.DurationOr("RWTimeout", 3000) * time.Millisecond,
		IdleTimeout:    poolOpts.DurationOr("IdleTimeout", 1200) * time.Second,
		MaxActive:      maxActive,
		MaxCap:         poolOpts.IntAndOr("MaxCap", fmt.Sprintf("N>%d", maxActive), maxActive),
		Logger:         logger,
		Trace:          nil,
		CircuitBreaker: circuitbreaker.New("playerService", CliAppErrTypes),
		HuskarConfiger: nex.GetHuskarConfiger(),
		NexConfig:      nexConfig,
	}
	if nexConfig.Config("plugins").Bool("etrace") {
		options.Trace = nex.GetETrace()
	}
	if !fromHuskar {
		options.Addr = nex.GetNexConfig().String("addr")
		return NewThriftplayerServiceClient(options), nil
	}
	return NewThriftplayerServiceClientWithHuskar(nex.GetHuskarPool(), options)
}

func findConfigFromRPCDepsFile(name string) (cc.Configer, error) {
	content, err := ioutil.ReadFile("thriftfs/deps.json")
	if err != nil {
		return nil, err
	}
	var deps []map[string]interface{}
	if err := json.Unmarshal(content, &deps); err != nil {
		return nil, err
	}
	for _, cfg := range deps {
		if cfg["Name"] == name {
			return cc.NewConfigFrom(cfg), nil
		}
	}
	return nil, fmt.Errorf("No config section found for '%v'", name)
}

type rawThriftClientWrapper struct {
	addr      string
	transport thrift.TTransport
	raw       *player.playerServiceClient
}

func newRawThriftClientWrapper(appName, addr string, connectTimeout, rwTimeout time.Duration) (*rawThriftClientWrapper, error) {
	socket, err := thrift.NewTSocket(addr)
	if err != nil {
		return nil, err
	}
	// set connect timeout
	socket.SetTimeout(connectTimeout)
	if err := socket.Open(); err != nil {
		return nil, err
	}
	return newRawThriftClientWrapperByConn(socket, appName, rwTimeout)
}

func newRawThriftClientWrapperByConn(socket *thrift.TSocket, appName string, rwTimeout time.Duration) (*rawThriftClientWrapper, error) {
	// set read/write timeout
	socket.SetTimeout(rwTimeout)
	transportFactory := thrift.NewTBufferedTransportFactory(consts.BufferSize)
	transport := transportFactory.GetTransport(socket)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tracker := ttracker.NewTracker(appName)
	client, err := player.NewplayerServiceClientFactory(tracker, transport, protocolFactory)
	if err != nil {
		return nil, err
	}
	return &rawThriftClientWrapper{
		addr:      socket.Conn().RemoteAddr().String(),
		raw:       client,
		transport: transport,
	}, nil
}

func (rtc *rawThriftClientWrapper) Ping() error {
	_, err := rtc.raw.Ping(context.TODO())
	return err
}

func (rtc *rawThriftClientWrapper) Close() error {
	return rtc.transport.Close()
}

// AssignTanks is used for client.
func (c ThriftplayerServiceClient) AssignTanks(ctx context.Context, tanks []int32) error {
	select {
	case <-ctx.Done():
		return timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		return err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "assignTanks")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := assigntanksRequest{client: client, Tanks: tanks}
	_, err = c.assigntanksEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		return err
	}
	return nil
}

// GetNewOrders is used for client.
func (c ThriftplayerServiceClient) GetNewOrders(ctx context.Context) ([]*player.Order, error) {
	select {
	case <-ctx.Done():
		var v []*player.Order
		return v, timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		var v []*player.Order
		return v, err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "getNewOrders")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := getnewordersRequest{client: client}
	response, err := c.getnewordersEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		var v []*player.Order
		return v, err
	}
	resp := response.(getnewordersResponse)
	return resp.V, err
}

// LatestState is used for client.
func (c ThriftplayerServiceClient) LatestState(ctx context.Context, state *player.GameState) error {
	select {
	case <-ctx.Done():
		return timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		return err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "latestState")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := lateststateRequest{client: client, State: state}
	_, err = c.lateststateEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		return err
	}
	return nil
}

// Ping is used for client.
func (c ThriftplayerServiceClient) Ping(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		var v bool
		return v, timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		var v bool
		return v, err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "ping")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := pingRequest{client: client}
	response, err := c.pingEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		var v bool
		return v, err
	}
	resp := response.(pingResponse)
	return resp.V, err
}

// UploadMap is used for client.
func (c ThriftplayerServiceClient) UploadMap(ctx context.Context, gamemap [][]int32) error {
	select {
	case <-ctx.Done():
		return timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		return err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "uploadMap")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := uploadmapRequest{client: client, Gamemap: gamemap}
	_, err = c.uploadmapEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		return err
	}
	return nil
}

// UploadParamters is used for client.
func (c ThriftplayerServiceClient) UploadParamters(ctx context.Context, arguments *player.Args_) error {
	select {
	case <-ctx.Done():
		return timeout.ErrTimeout
	default:
	}
	rawClient, err := c.pool.Get(ctx)
	if err != nil {
		etrace.LogErrorContext(ctx, "ConnectFailed", fmt.Sprintf("purchaseMeiTuan -- %s", c.pool.StatsReadableJSON()))
		return err
	}
	client := rawClient.(*rawThriftClientWrapper)
	// Passing context values.
	ctx = context.WithValue(ctx, ctxkeys.CliAPIName, "uploadParamters")
	ctx = context.WithValue(ctx, ctxkeys.RemoteAddr, client.addr)

	request := uploadparamtersRequest{client: client, Arguments: arguments}
	_, err = c.uploadparamtersEndpoint(ctx, request)
	if IsTolerableException(err) {
		c.pool.Put(client)
	} else {
		client.Close()
		c.pool.Put(nil)
	}
	if err != nil {
		return err
	}
	return nil
}

// Close closes all connections in pool.
func (c ThriftplayerServiceClient) Close() {
	c.pool.Close()
}

func makeThriftAssignTanksEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(assigntanksRequest)
		err := req.client.raw.AssignTanks(ctx, req.Tanks)
		response := assigntanksResponse{}
		return response, err
	}
}

func makeThriftGetNewOrdersEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getnewordersRequest)
		v, err := req.client.raw.GetNewOrders(ctx)
		response := getnewordersResponse{V: v}
		return response, err
	}
}

func makeThriftLatestStateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(lateststateRequest)
		err := req.client.raw.LatestState(ctx, req.State)
		response := lateststateResponse{}
		return response, err
	}
}

func makeThriftPingEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pingRequest)
		v, err := req.client.raw.Ping(ctx)
		response := pingResponse{V: v}
		return response, err
	}
}

func makeThriftUploadMapEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadmapRequest)
		err := req.client.raw.UploadMap(ctx, req.Gamemap)
		response := uploadmapResponse{}
		return response, err
	}
}

func makeThriftUploadParamtersEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadparamtersRequest)
		err := req.client.raw.UploadParamters(ctx, req.Arguments)
		response := uploadparamtersResponse{}
		return response, err
	}
}

type assigntanksRequest struct {
	Tanks  []int32
	client *rawThriftClientWrapper
}

type assigntanksResponse struct {
}

type getnewordersRequest struct {
	client *rawThriftClientWrapper
}

type getnewordersResponse struct {
	V []*player.Order
}

type lateststateRequest struct {
	State  *player.GameState
	client *rawThriftClientWrapper
}

type lateststateResponse struct {
}

type pingRequest struct {
	client *rawThriftClientWrapper
}

type pingResponse struct {
	V bool
}

type uploadmapRequest struct {
	Gamemap [][]int32
	client  *rawThriftClientWrapper
}

type uploadmapResponse struct {
}

type uploadparamtersRequest struct {
	Arguments *player.Args_
	client    *rawThriftClientWrapper
}

type uploadparamtersResponse struct {
}
