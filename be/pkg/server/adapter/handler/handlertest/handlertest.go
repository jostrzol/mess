package handlertest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/logger"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"github.com/stretchr/testify/suite"
)

func SetupRouter() *gin.Engine {
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

type HandlerSuite struct {
	suite.Suite
	g   *gin.Engine
	jar *cookiejar.Jar
}

func (s *HandlerSuite) SetupTest() {
	var err error
	s.g = SetupRouter()
	s.jar, err = cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
}

func (s *HandlerSuite) ServeHTTPOkAs(method string, url string, body JSON, result interface{}) {
	res := s.ServeHTTPOk(method, url, body)
	bytes, err := io.ReadAll(res.Result().Body)
	s.NoError(err)

	err = json.Unmarshal(bytes, &result)
	s.NoError(err)
}

func (s *HandlerSuite) ServeHTTPOk(method string, url string, body JSON) *httptest.ResponseRecorder {
	res := s.ServeHTTP(method, url, body)
	s.Equal(res.Result().StatusCode, http.StatusOK)
	return res
}

func (s *HandlerSuite) ServeHTTP(method string, url string, body JSON) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()

	bodyBytes, err := json.Marshal(body)
	s.NoError(err)
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	for _, cookie := range s.jar.Cookies(&root) {
		req.AddCookie(cookie)
	}
	s.NoError(err)

	s.g.ServeHTTP(res, req)
	s.jar.SetCookies(&root, res.Result().Cookies())
	return res
}

type JSON = map[string]any

var root = url.URL{Path: "/"}
