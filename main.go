package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		// panic(err.Error())
	}
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func download(index, n int, token string) {
	defer wg.Done()
	up := index + n
	for index < up {
		index++
		indexStr := fmt.Sprintf("%d", index)
		for len(indexStr) != 4 {
			indexStr = "0" + indexStr
		}

		// if ok := Exist("./" + indexStr + ".jpeg"); ok {
		// 	fmt.Println(indexStr + "存在")
		// 	continue
		// }

		fmt.Println(indexStr)

		reqUrl := "https://ia903107.us.archive.org/BookReader/BookReaderImages.php?zip=/4/items/shuzhifenxi0005unse/shuzhifenxi0005unse_jp2.zip&file=shuzhifenxi0005unse_jp2/shuzhifenxi0005unse_" + indexStr + ".jp2&scale=4&rotate=0"
		method := "GET"

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {
					return url.Parse("http://127.0.0.1:1080")
				},
			},
			Timeout: 10 * time.Second,
		}
		req, err := http.NewRequest(method, reqUrl, nil)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Cookie", "test-cookie=1; PHPSESSID=spibqj5l387dmv6agfjt6r5ut3; logged-in-sig=1614138253+1582602253+E%2BkjWrUXPKbhT0FnHAF53q6SPqQbIAvR%2Fm3oC4WV%2BH9bzWmUZcGxWk5GN6eiAjrmwTOBhXWpZB0jqx0cGfCjCdFUoN58MY1BAZ4zpTYMCBDyOc2CL%2BxE8VOFBdQ9uiMWVlTfjcuGaBxtdV4JHeAzz04r9u3oIKXnfSdF8Ly%2FfkU%3D; logged-in-user=905529001%40qq.com; br-loan-shuzhifenxi0005unse=1631169966; ol-auth-url=%2F%2Farchive.org%2Fservices%2Fborrow%2FXXX%3Fmode%3Dauth;")
		req.Header.Add("Referer", "https://archive.org/details/shuzhifenxi0005unse/page/n12/mode/1up")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
		req.Header.Add("Cookie", "loan-shuzhifenxi0005unse="+token)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		handleErr(err)

		_, suffix, _ := image.DecodeConfig(strings.NewReader(string(bytes)))
		if suffix == "" { //判断下载是是否为图片
			fmt.Println(string(bytes))
			fmt.Println(indexStr, "非图片")
		}

		fmt.Println(suffix)
		ioutil.WriteFile(indexStr+"."+suffix, bytes, 0777)
	}
}

var wg sync.WaitGroup

func getNewToken(oldToken string) string {
	urll := "https://archive.org/services/loans/beta/loan/"
	payload := url.Values{
		"action":     {"create_token"},
		"identifier": {"shuzhifenxi0005unse"},
	}.Encode()
	req, err := http.NewRequest("POST", urll, strings.NewReader(payload))
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:1080")
			},
		},
	}
	req.Header = map[string][]string{
		"Cookie":       {"test-cookie=1; PHPSESSID=spibqj5l387dmv6agfjt6r5ut3; logged-in-sig=1614138253+1582602253+E%2BkjWrUXPKbhT0FnHAF53q6SPqQbIAvR%2Fm3oC4WV%2BH9bzWmUZcGxWk5GN6eiAjrmwTOBhXWpZB0jqx0cGfCjCdFUoN58MY1BAZ4zpTYMCBDyOc2CL%2BxE8VOFBdQ9uiMWVlTfjcuGaBxtdV4JHeAzz04r9u3oIKXnfSdF8Ly%2FfkU%3D; logged-in-user=905529001%40qq.com; br-loan-shuzhifenxi0005unse=1631169966; ol-auth-url=%2F%2Farchive.org%2Fservices%2Fborrow%2FXXX%3Fmode%3Dauth;"},
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Referer":      {"https://archive.org/details/shuzhifenxi0005unse/page/n12/mode/1up"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"},
	}
	req.Header.Add("Cookie", "loan-shuzhifenxi0005unse="+oldToken)
	resp, err := client.Do(req)
	handleErr(err)
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	var data map[string]interface{}
	json.Unmarshal(bytes, &data)
	ret := data["token"].(string)
	return ret
}

var token = "1582617257-4574e29ef62fda3433f2c12a1600ed9b"
var n = 10

func main() {
	// index := 300
	// for {
	// 	token = getNewToken(token)
	// 	fmt.Println("newToken", token)
	// 	wg.Add(6)
	// 	go download(index, n, token)
	// 	index += n
	// 	go download(index, n, token)
	// 	index += n
	// 	go download(index, n, token)
	// 	index += n
	// 	go download(index, n, token)
	// 	index += n
	// 	go download(index, n, token)
	// 	index += n
	// 	go download(index, n, token)
	// 	wg.Wait()
	// 	if index > 330 {
	// 		return
	// 	}
	// }

	// token = getNewToken(token)
	// wg.Add(1)
	// download(315, 1, token)
	// findLost(token)

	makePDF()
}

func findLost(token string) {
	for index := 1; index < 342; index++ {
		fmt.Println("check ", index)
		indexStr := fmt.Sprintf("%d", index)
		for len(indexStr) != 4 {
			indexStr = "0" + indexStr
		}
		if ok := Exist(indexStr + ".jpeg"); !ok {
			fmt.Println(indexStr, "不存在 , 尝试下载:")

			// wg.Add(1)
			// download(index, 1, token)
			return
		}
	}
}

func makePDF() {
	var (
		w float64 = float64(820) * 0.23
		h float64 = float64(1075) * 0.23
	)

	pdf := gofpdf.New("P", "mm", "A4", "")
	for index := 1; index < 343; index++ {
		indexStr := fmt.Sprintf("%d", index)
		for len(indexStr) != 4 {
			indexStr = "0" + indexStr
		}
		pdf.AddPage()
		pdf.Image(indexStr+".jpeg", 10, 10, w, h, true, "", 0, "")
	}
	err := pdf.OutputFileAndClose("test.pdf")
	if err != nil {
		fmt.Println(err)
	}
}
