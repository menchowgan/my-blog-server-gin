package weather

import (
	"gmc-blog-server/response"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetWeather(c *gin.Context) error {
	city := c.Query("city")

	if city == "" {
		response.Fail(http.StatusInternalServerError, nil, "城市不能为空", c)
		return nil
	}

	// 调用天气API获取天气信息
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://restapi.amap.com/v3/weather/weatherInfo?key=43a2708e1798446ab929772b924f7ca6&city="+city+"", nil)

	if err != nil {
		return err
	}

	res, _ := client.Do(req)

	defer func() {
		res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	response.Success(string(body), "天气查询成功", c)
	return nil
}
