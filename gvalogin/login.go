package gvalogin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	proto "github.com/tscuite/tscuite-operator/operator-proto"
	"google.golang.org/grpc"
)

type Login struct {
}

var LoginApp = new(Login)

type Captcha struct {
	Code int `json:"code"`
	Data struct {
		CaptchaID     string `json:"captchaId"`
		PicPath       string `json:"picPath"`
		CaptchaLength int    `json:"captchaLength"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Loginres struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			ID          int       `json:"ID"`
			CreatedAt   time.Time `json:"CreatedAt"`
			UpdatedAt   time.Time `json:"UpdatedAt"`
			UUID        string    `json:"uuid"`
			UserName    string    `json:"userName"`
			NickName    string    `json:"nickName"`
			SideMode    string    `json:"sideMode"`
			HeaderImg   string    `json:"headerImg"`
			BaseColor   string    `json:"baseColor"`
			ActiveColor string    `json:"activeColor"`
			AuthorityID string    `json:"authorityId"`
			Authority   struct {
				CreatedAt       time.Time   `json:"CreatedAt"`
				UpdatedAt       time.Time   `json:"UpdatedAt"`
				DeletedAt       interface{} `json:"DeletedAt"`
				AuthorityID     string      `json:"authorityId"`
				AuthorityName   string      `json:"authorityName"`
				ParentID        string      `json:"parentId"`
				DataAuthorityID interface{} `json:"dataAuthorityId"`
				Children        interface{} `json:"children"`
				Menus           interface{} `json:"menus"`
				DefaultRouter   string      `json:"defaultRouter"`
			} `json:"authority"`
			Authorities []struct {
				CreatedAt       time.Time   `json:"CreatedAt"`
				UpdatedAt       time.Time   `json:"UpdatedAt"`
				DeletedAt       interface{} `json:"DeletedAt"`
				AuthorityID     string      `json:"authorityId"`
				AuthorityName   string      `json:"authorityName"`
				ParentID        string      `json:"parentId"`
				DataAuthorityID interface{} `json:"dataAuthorityId"`
				Children        interface{} `json:"children"`
				Menus           interface{} `json:"menus"`
				DefaultRouter   string      `json:"defaultRouter"`
			} `json:"authorities"`
		} `json:"user"`
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"data"`
	Msg string `json:"msg"`
}

//grpc url??????
const PORT string = "127.0.0.1:7001"


//??????token
func (login *Login) Xtoken(url string, user string, passwd string) (string, string) {
	captcha, captchaId := login.Ddddorc(url)
	fmt.Println(captcha)
	content := `{"username": "`+ user +`", "Password":"`+ passwd +`", "captcha": "` + captcha + `", "captchaId": "` + captchaId + `"}`
	body, err := login.Client(url+"/api/base/login", content)
	if err != nil {
		fmt.Printf("login: %v\n", err)
	}
	return login.JsonLogin(body).Data.Token, login.JsonLogin(body).Msg
}

//???????????????
func (login *Login) Ddddorc(url string) (string, string) {
	content := `{}`
	body, err := login.Client(url+"/api/base/captcha", content)
	if err != nil {
		fmt.Printf("Ddddorc: %v\n", err)
	}
	return login.Grpc(login.JsonCaptcha(body).Data.PicPath), login.JsonCaptcha(body).Data.CaptchaID
}

//?????????json
func (login *Login) JsonCaptcha(body []byte) *Captcha {
	var a Captcha
	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Printf("JsonCaptcha: %v\n", err)
	}
	return &a
}

//?????????json
func (login *Login) JsonLogin(body []byte) *Loginres {
	var a Loginres
	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Printf("JsonLogin: %v\n", err)
	}
	return &a
}

//grpc?????????????????????
func (login *Login) Grpc(request string) string {
	conn, err := grpc.Dial(PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc????????????!", err.Error())
	}
	defer conn.Close()
	client := proto.NewSearchServiceClient(conn)
	req, err := client.Search(context.Background(), &proto.SearchRequest{
		Request: request,
	})
	if err != nil {
		log.Fatal("grpc??????????????????!", err.Error())
	}
	return req.GetResponse()
}

//post??????????????????
func (login *Login) Client(url, content string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		fmt.Printf("NewRequest: %v\n", err)
	}
	req.Header.Set("x-user-id", "1")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: %v\n", err)
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

//post??????????????????
func (login *Login) AllClient(xtoken, url, content string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		fmt.Printf("NewRequest: %v\n", err)
	}
	req.Header.Set("x-user-id", "1")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-token", xtoken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: %v\n", err)
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func Gvalogin(url string,user string, passwd string ,Map map[string]string ) {
	for {
		if xtoken, err := LoginApp.Xtoken(url,user,passwd); err != "????????????" {
			fmt.Printf("login: %v\n", err)
		} else {
			var Json map[string]interface{}
			fmt.Println(xtoken)
			for url, key := range Map {
				body, err := LoginApp.AllClient(xtoken, url, key)
				if err != nil {
					fmt.Printf("login: %v\n", err)
				} else {
					json.Unmarshal(body, &Json)
					fmt.Println(Json)
				}
			}
			break
		}
	}
}