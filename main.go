package main

import (
	"bytes"
	_ "corsMiddlleware/docs" // docs is generated by Swag CLI, you have to import it.
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

// MainHandler godoc
// @title CORS Midlleware
// @summary CORS Midlleware
// @description Сервис-посредник, для работы с rest API, с защитой КОРС
// @accept multipart/form-data

// @title CORS Midlleware
// @version 1.0
// @description Сервис-посредник, для работы с rest API, с защитой КОРС
// @termsOfService http://swagger.io/terms/

// @contact.name Constantine
// @contact.url https://github.com/Constantine-Ka/corsmiddleware

// @license.name Apache 2.0
// @Param _method		formData string false 	"http-метод, с которым нужно отправить запрос"
// @Param _url			formData string true 	"URL по которому нужно выполнить запрос"
// @Param _headers		formData string false 	"Заголовки, которые нужно передать"
// @Param _json			formData string false 	"Тело запроса, которые нужно передать в raw"
// @Param прочее		formData string false 	"Любые поля и значения в формате formdata"
// @Success 200 {object} interface{} "Response from the target URL"
// @Failure 400 {object} map[string]string "Bad Request.  Check the error message for details."
// @Router / [post]
// @host localhost:8008
// @BasePath /
func main() {
	port := flag.Int64("port", 8008, "Порт сервера")
	prometheusPort := flag.Int64("prometheus", 0, "Порт который будет слушать промитэус")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("myapp")) // adds middleware to gather metrics
	if *prometheusPort > 0 {
		go func() {
			metrics := echo.New()                                // this Echo will run on separate port 8081
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if err := metrics.Start(fmt.Sprintf(":%d", *prometheusPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}()
	}
	e.GET("/*", echoSwagger.WrapHandler)
	e.POST("/", MainHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}

func MainHandler(c echo.Context) error {
	//var payload io.Reader
	var req *http.Request
	method := strings.ToUpper(c.FormValue("_method"))
	headersStr := c.FormValue("_headers")
	jsonStr := c.FormValue("_json")
	urlAddr := c.FormValue("_url")
	formValues, err := c.FormParams()
	if err != nil {
		return err
	}

	if urlAddr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"result": "URL required",
		})
	}
	if method == "" {
		method = "GET"
	}

	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"result": "invalid method",
		})
	}
	if method == "GET" {
		var quers []string
		for key, values := range formValues {
			if len(values) == 1 {
				quers = append(quers, fmt.Sprintf("%s=%s", key, values[0]))
			} else {
				for _, value := range values {
					quers = append(quers, fmt.Sprintf("%s[]=%s", key, value))
				}
			}
		}
		urlAddr = fmt.Sprintf("%s?%s", urlAddr, strings.Join(quers, "&"))
		req, err = http.NewRequestWithContext(c.Request().Context(), method, urlAddr, nil)
	} else if jsonStr != "" {
		payload := strings.NewReader(jsonStr)
		req, err = http.NewRequestWithContext(c.Request().Context(), method, urlAddr, payload)
	} else {
		p := &bytes.Buffer{}
		writer := multipart.NewWriter(p)
		for key, values := range formValues {
			if len(values) == 1 {
				_ = writer.WriteField(key, values[0])
			} else {
				for _, value := range values {
					writer.WriteField(fmt.Sprintf("%s[]=", key), value)

				}
			}
		}
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{})
		}
		req, err = http.NewRequestWithContext(c.Request().Context(), method, urlAddr, p)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"result": err.Error(),
				"event":  "NewRequestError",
			})
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

	}

	if headersStr != "" {
		var headerMap map[string]string
		err = json.Unmarshal([]byte(headersStr), &headerMap)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"result": "Header.Invalid json format",
				"event":  "json.Unmarshal([]byte(headersStr), &headerMap)",
			})
		}
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}
	cookies := c.Request().Cookies()
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": err.Error(),
			"event":  "client.Do(req)",
		})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"result": err.Error(),
			"event":  "ReadAllError",
		})
	}
	defer resp.Body.Close()
	return c.JSONBlob(resp.StatusCode, body)

}
