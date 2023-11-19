package handlertest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/logger"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite[T Client] struct {
	suite.Suite
	defaultClient *T
	g             *gin.Engine
}

func (s *HandlerSuite[T]) SetupTest() {
	s.g = setupRouter()
	s.defaultClient = s.NewClient()
	setupDir()
}

func setupRouter() *gin.Engine {
	config := &serverconfig.Config{
		IsProduction:   false,
		SessionSecret:  "secret",
		Port:           54321,
		IncomingOrigin: "http://localhost:4000",
	}
	ioc.MustSingleton(config)
	logger, err := logger.New(config.IsProduction)
	if err != nil {
		panic(err)
	}
	ioc.MustSingleton(logger)
	return ioc.MustResolve[*gin.Engine]()
}

// setupDir changes working directory to module's root.
// It is necessary, so that we can load rules correctly.
// TODO: remove when rules are handled dynamically.
func setupDir() {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "../../../../..")
	err := os.Chdir(root)
	if err != nil {
		panic(err)
	}
}

func (s *HandlerSuite[T]) Client() *T {
	return s.defaultClient
}

func (s *HandlerSuite[T]) NewClient() *T {
	var client T
	v := reflect.Indirect(reflect.ValueOf(&client))
	v.FieldByName("BaseClient").Set(reflect.ValueOf(new(BaseClient)))
	client.initClient(newHTTPClient(&s.Suite, s.g))
	return &client
}

type Client interface {
	initClient(*httpClient)
	client() *httpClient
}

type BaseClient struct {
	*httpClient
}

func (c *BaseClient) initClient(httpClient *httpClient) {
	c.httpClient = httpClient
}

func (c *BaseClient) client() *httpClient {
	return c.httpClient
}

func CloneWithEmptyJar[T Client, TP interface {
	Client
	*T
}](c TP) TP {
	var client T
	v := reflect.Indirect(reflect.ValueOf(&client))
	v.FieldByName("BaseClient").Set(reflect.ValueOf(new(BaseClient)))
	client.initClient(newHTTPClient(c.client().Suite, c.client().g))
	return &client
}

type httpClient struct {
	*suite.Suite
	g   *gin.Engine
	jar *cookiejar.Jar
}

func newHTTPClient(suite *suite.Suite, g *gin.Engine) *httpClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	return &httpClient{Suite: suite, g: g, jar: jar}
}

func (c *httpClient) cloneWithEmptyJar() *httpClient {
	return newHTTPClient(c.Suite, c.g)
}

func (c *httpClient) ServeJSONOkAs(method string, url string, body any, result interface{}) {
	c.T().Helper()
	res := c.ServeJSONOk(method, url, body)
	bytes, err := io.ReadAll(res.Result().Body)
	c.NoError(err)

	err = json.Unmarshal(bytes, &result)
	c.NoError(err)
}

func (c *httpClient) ServeJSONOk(method string, url string, body any) *httptest.ResponseRecorder {
	c.T().Helper()
	return c.invokeMarshalled(c.ServeOk, method, url, body)
}

func (c *httpClient) ServeJSON(method string, url string, body any) *httptest.ResponseRecorder {
	c.T().Helper()
	return c.invokeMarshalled(c.Serve, method, url, body)
}

func (c *httpClient) invokeMarshalled(
	action func(string, string, []byte) *httptest.ResponseRecorder,
	method string,
	url string,
	body any,
) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	c.NoError(err)
	return action(method, url, bodyBytes)
}

func (c *httpClient) ServeOk(method string, url string, body []byte) *httptest.ResponseRecorder {
	c.T().Helper()
	res := c.Serve(method, url, body)
	status := res.Result().StatusCode
	c.True(200 <= status && status < 300)

	return res
}

func (c *httpClient) Serve(method string, url string, body []byte) *httptest.ResponseRecorder {
	c.T().Helper()
	res := httptest.NewRecorder()

	req := c.request(method, url, body)
	c.g.ServeHTTP(res, req)
	c.jar.SetCookies(&root, res.Result().Cookies())

	c.logRequest(method, url, body, res)
	return res
}

func (c *httpClient) request(method string, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	c.NoError(err)
	for _, cookie := range c.jar.Cookies(&root) {
		req.AddCookie(cookie)
	}
	return req
}

func (c *httpClient) logRequest(method string, url string, reqBody []byte, res *httptest.ResponseRecorder) {
	req := c.request(method, url, reqBody)
	c.T().Logf("-> request url: %v %v", req.Method, req.URL)
	c.T().Logf("-> request headers: %+v", spew.Sdump(req.Header))
	if reqBody != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			c.T().Logf("-> request body: <can't read>")
		} else {
			c.T().Logf("-> request body: %v", string(body))
		}
	}
	c.T().Logf("--------------------------------------------------------------------------------")
	var newBody bytes.Buffer
	resReader := io.TeeReader(res.Body, &newBody)
	body, err := io.ReadAll(resReader)
	if err != nil {
		c.T().Logf("<- response body: <can't read>")
	} else {
		c.T().Logf("<- response body: %v", string(body))
	}
	res.Body = &newBody
	c.T().Logf("================================================================================")
}

var root = url.URL{Scheme: "http", Path: "/"}
