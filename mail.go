package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailjet/mailjet-apiv3-go/v3"
)

type MailRequestBody struct {
	Name    string
	Email   string
	Message string
}

type CaptchaResponse struct {
	Success bool `json:"success"`
}

func SendMailHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method == "POST" {
		// Check captcha
		reCaptchaRequestBody := url.Values{}
		reCaptchaRequestBody.Add("secret", "6Lc-chYUAAAAAIrehupr7j3iS8JDFZ3vHfFOSnvt")
		reCaptchaRequestBody.Add("response", r.PostFormValue("g-recaptcha-response"))
		reCaptchaRequestBody.Add("remoteip", r.RemoteAddr)

		var captchaClient = &http.Client{Timeout: 10 * time.Second}
		var captchaResponse = CaptchaResponse{}

		response, _ := captchaClient.PostForm("https://www.google.com/recaptcha/api/siteverify", reCaptchaRequestBody)
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &captchaResponse)

		if captchaResponse.Success != true {
			w.WriteHeader(429)
			fmt.Fprintf(w, "Message sent successful!")
		} else {
			mailRequestBody := MailRequestBody{
				Name:    r.PostFormValue("name"),
				Email:   r.PostFormValue("email"),
				Message: r.PostFormValue("comments"),
			}

			publicKey := os.Getenv("MJ_APIKEY_PUBLIC")
			secretKey := os.Getenv("MJ_APIKEY_PRIVATE")

			mj := mailjet.NewMailjetClient(publicKey, secretKey)

			param := &mailjet.InfoSendMail{
				FromEmail: "w",
				FromName:  "Mailer",
				To:        "jefdauda@gmail.com",
				Subject:   "Resume contact form",
				Headers: map[string]string{
					"Reply-To": mailRequestBody.Name + "<" + mailRequestBody.Email + ">",
				},
				TextPart: mailRequestBody.Message,
			}

			_, err = mj.SendMail(param)

			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "ERROR")
				log.Println(err)
			} else {
				w.WriteHeader(200)
				fmt.Fprintf(w, "SUCCESS")
			}
		}

	} else {
		w.WriteHeader(404)
		fmt.Fprintf(w, "NOT FOUND")
	}
}

func main() {
	home := http.FileServer(http.Dir("./static"))
	http.Handle("/", home)
	http.HandleFunc("/sendmail", SendMailHandler)

	fmt.Printf("Server started at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
	env := godotenv.Load(".env")
	if env != nil {
		log.Fatalf("Error loading .env file")
	}

}
