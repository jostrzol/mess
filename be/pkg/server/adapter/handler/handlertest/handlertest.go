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
		IsProduction:  false,
		SessionSecret: "secret",
		Port:          54321,
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

func CloneWithEmptyJar[T Client](c T) T {
	result := c
	newClient := result.client().cloneWithEmptyJar()
	result.initClient(newClient)
	return result
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

func (c *httpClient) ServeHTTPOkAs(method string, url string, body JSON, result interface{}) {
	res := c.ServeHTTPOk(method, url, body)
	bytes, err := io.ReadAll(res.Result().Body)
	c.NoError(err)

	err = json.Unmarshal(bytes, &result)
	c.NoError(err)
}

func (c *httpClient) ServeHTTPOk(method string, url string, body JSON) *httptest.ResponseRecorder {
	res := c.ServeHTTP(method, url, body)
	c.Equal(res.Result().StatusCode, http.StatusOK)
	return res
}

func (c *httpClient) ServeHTTP(method string, url string, body JSON) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()

	bodyBytes, err := json.Marshal(body)
	c.NoError(err)
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	for _, cookie := range c.jar.Cookies(&root) {
		req.AddCookie(cookie)
	}
	c.NoError(err)

	c.g.ServeHTTP(res, req)
	c.jar.SetCookies(&root, res.Result().Cookies())
	return res
}

type JSON = map[string]any

var root = url.URL{Scheme: "http", Path: "/"}
