package service

import (
	"errors"
	"net/url"
	"regexp"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/xlog"
	"go.uber.org/zap"
)

var (
	influxdbClient             influxdb.Client
	defaultInfluxdbBatchPoints influxdb.BatchPoints
	defaultInfluxdbConfig      *InfluxdbConfig
	defaultWritePointMutex     = new(sync.Mutex)
)

const (
	maxBatchSize = 20
)

type (
	// InfluxdbConfig influxdb's config
	InfluxdbConfig struct {
		Addr      string
		Username  string
		Password  string
		Db        string
		Precision string
	}
)

func init() {
	uri := config.GetString("influxdb")
	if uri != "" {
		conf, err := getInfluxdbConfig(uri)
		if err != nil {
			panic(err)
		}
		defaultInfluxdbConfig = conf
		c, err := newInfluxdbClient(conf)
		if err != nil {
			panic(err)
		}

		influxdbClient = c
		logger := xlog.Logger()
		mask := regexp.MustCompile(`http://(\S+):(\S+)\@`)
		str := mask.ReplaceAllString(uri, "http://:***@")
		_, _, err = influxdbClient.Ping(5 * time.Second)
		if err != nil {
			logger.Error("influxdb ping fail",
				zap.String("uri", str),
				zap.Error(err),
			)
		} else {
			logger.Info("influxdb ping success",
				zap.String("uri", str),
			)
		}

	}
}

func getInfluxdbConfig(uri string) (c *InfluxdbConfig, err error) {
	info, err := url.Parse(uri)
	if err != nil {
		return
	}
	query := info.Query()
	db := query.Get("db")
	if db == "" {
		err = errors.New("db can not be nil")
		return
	}
	precision := query.Get("precision")
	// 设置默认的precision为ns
	if precision == "" {
		precision = "ns"
	}
	c = &InfluxdbConfig{
		Addr:      info.Scheme + "://" + info.Host,
		Db:        db,
		Precision: precision,
	}
	userName := info.User.Username()
	pwd, _ := info.User.Password()
	if userName != "" && pwd != "" {
		c.Username = userName
		c.Password = pwd
	}
	return
}

func newInfluxdbClient(conf *InfluxdbConfig) (client influxdb.Client, err error) {
	httpConfig := influxdb.HTTPConfig{
		Addr: conf.Addr,
	}
	if conf.Username != "" {
		httpConfig.Username = conf.Username
		httpConfig.Password = conf.Password
	}
	client, err = influxdb.NewHTTPClient(httpConfig)
	return
}

// GetInfluxdbClient get influxdb's client
func GetInfluxdbClient() influxdb.Client {
	return influxdbClient
}

// WriteInfluxPoint write point to influxdb
func WriteInfluxPoint(name string, tags map[string]string, fields map[string]interface{}) (err error) {
	if influxdbClient == nil {
		return
	}
	defaultWritePointMutex.Lock()
	defer defaultWritePointMutex.Unlock()
	pt, err := influxdb.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return
	}
	if defaultInfluxdbBatchPoints == nil {
		defaultInfluxdbBatchPoints, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  defaultInfluxdbConfig.Db,
			Precision: defaultInfluxdbConfig.Precision,
		})
		if err != nil {
			return
		}
	}
	defaultInfluxdbBatchPoints.AddPoint(pt)

	if len(defaultInfluxdbBatchPoints.Points()) >= maxBatchSize {
		bp := defaultInfluxdbBatchPoints
		go func() {
			err := influxdbClient.Write(bp)
			if err != nil {
				xlog.Logger().Error("influxdb batch write fail",
					zap.Error(err),
				)
			}
		}()
		defaultInfluxdbBatchPoints = nil
	}
	return
}
