package testEnv

import (
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"github.com/valyala/tcplisten"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Init() {
	log.InitLogger(log.DebugLevel, "")
	config.Set(config.TestMode, true)
	config.Set(config.MongoUrl, "mongodb://localhost:27001")
	config.Set(config.MongoDB, "blip")
	config.Set(config.RedisUrl, "localhost:6379")
	config.Set(config.LogLevel, log.DebugLevel)
	config.Set(config.SongsIndexDir, "./_hdd")

	_ = os.MkdirAll(config.GetString(config.SongsIndexDir), os.ModePerm)

	// Initialize MongoDB
	mongoClient, err := mongo.Connect(
		nil,
		options.Client().ApplyURI(config.GetString(config.MongoUrl)).SetDirect(true),
	)
	if err != nil {
		log.Fatal("Error On MongoConnect", zap.Error(err))
	}
	err = mongoClient.Ping(nil, nil)
	if err != nil {
		log.Fatal("Error On MongoConnect (Ping)", zap.Error(err))
	}
	auth.InitMongo(mongoClient)
	crawler.InitMongo(mongoClient)
	music.InitMongo(mongoClient)
	session.InitMongo(mongoClient)
	token.InitMongo(mongoClient)
	user.InitMongo(mongoClient)

	// Initialize RedisCache
	redisConfig := redis.DefaultConfig
	redisConfig.Host = config.GetString(config.RedisUrl)
	redisConfig.Password = config.GetString(config.RedisPass)
	redisCache := redis.New(redisConfig)
	auth.InitRedisCache(redisCache)
	user.InitRedisCache(redisCache)

	// Initialize Modules
	auth.Init()
	crawler.Init()
	music.Init()
	session.Init()
	token.Init()
	user.Init()

}

type mockCrawler struct {
	MaxDelay time.Duration
	Port     int
}

func (m mockCrawler) ServeHTTP(httpRes http.ResponseWriter, httpReq *http.Request) {
	if m.MaxDelay > 0 {
		time.Sleep(time.Duration(tools.RandomInt(int(m.MaxDelay))))
	}
	reqData, _ := ioutil.ReadAll(httpReq.Body)
	req := crawler.SearchRequest{}
	_ = req.UnmarshalJSON(reqData)
	res := crawler.SearchResponse{
		RequestID: req.RequestID,
		Source:    fmt.Sprintf("Source %d", m.Port%3),
	}
	for i := 0; i < 10; i++ {
		res.Result = append(res.Result, crawler.FoundSong{
			SongUrl:  "http://url.com",
			CoverUrl: "http://cover-url.com",
			Lyrics:   "This is some lyrics text",
			Artists:  fmt.Sprintf("Some Famous Artist (%s)", tools.RandomID(4)),
			Title:    fmt.Sprintf("Title %s", tools.RandomID(4)),
			Genre:    "Rock",
		})
	}

	resData, _ := res.MarshalJSON()
	httpRes.Write(resData)
}

func InitMockCrawler(maxDelay time.Duration, port int) {
	s := httptest.NewUnstartedServer(mockCrawler{
		MaxDelay: maxDelay,
		Port:     port,
	})
	tcpConfig := tcplisten.Config{}
	s.Listener, _ = tcpConfig.NewListener("tcp4", fmt.Sprintf(":%d", port))
	s.Start()

}

func InitMultiCrawlers(n int, maxDelay time.Duration, startPort int) {
	for i := 0; i < n; i++ {
		InitMockCrawler(maxDelay, startPort+i)
		crawlerX := &crawler.Crawler{
			ID:          primitive.NewObjectID(),
			Url:         fmt.Sprintf("http://localhost:%d", startPort+i),
			Name:        fmt.Sprintf("Crawler %d", i),
			Description: "This is a Mock Crawler",
			Source:      fmt.Sprintf("Source %d", i%3),
		}
		_, _ = crawler.Save(crawlerX)
	}
}
