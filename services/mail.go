package services

import (
    "log"
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
	"mime/multipart"

	"github.com/ninjadotorg/handshake-dispatcher/utils"
)

type MailService struct{}

func (s MailService) Send(from string, to string, subject string, content string) (success bool, err error) {
	endpoint, _ := utils.GetServicesEndpoint("mail")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormField("from")
	part.Write([]byte(from))
	part, _ = writer.CreateFormField("to[]")
	part.Write([]byte(to))
	part, _ = writer.CreateFormField("subject")
	part.Write([]byte(subject))
	part, _ = writer.CreateFormField("body")
	part.Write([]byte(content))
	writer.Close()

	request, _ := http.NewRequest("POST", endpoint, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(b, &data)

	fmt.Println(data)
	status, _ := data["status"].(float64)
	success = status >= 1

	return
}

func (s MailService) SendCompleteProfile(email string, username string, hash string) {
    endpoint := utils.GetServerEndpoint()
    refLink := fmt.Sprintf("%s?ref=%s", endpoint, username)

    subject := `You've got 80 Shurikens (SHURI) from Ninja AIRDROP`
    body := fmt.Sprintf(`Greetings,<br/><br/> 
        %s<br/><br/>
        Welcome to the dojo. 80 shurikens (SHURI) have been added to <a href="%s">your wallet</a>.<br/><br/>  
        Your TxHash: <a href="https://etherscan.io/tx/%s">%s</a><br/><br/> 
        Use these tokens to slash fees, learn secrets, unlock special treatment.<br/><br/> 
        This is your unique referral link: <a href="%s">%s</a><br/><br/> 
        Bring your most trusted friends to the dojo and receive 20 extra shurikens (SHURI) for each new recruit.`, username, endpoint, hash, hash, refLink, refLink)
    
    status, err := s.Send("dojo@ninja.org", email, subject, body)
    log.Println("Send mail CompleteProfile status", status, err)
}

func (s MailService) SendCompleteReferrer(email string, username string, hash string) {
    endpoint := utils.GetServerEndpoint()
    refLink := fmt.Sprintf("%s?ref=%s", endpoint, username)
    
    subject := `You got 20 Shurikens (SHURI) for your new ninja recruit.`
    body := fmt.Sprintf(`Hi again,<br/><br/>
        %s<br/><br/>
        Thanks for bringing a new ninja to the dojo. 20 shurikens (SHURI) have been added to <a href="%s">your wallet</a>.<br/><br/>
        Your TxHash: <a href="https://etherscan.io/tx/%s">%s</a><br/><br/> 
        Keep fighting the good fight. Spread your unique referral link: <a href="%s">%s</a><br/><br/> 
        This does not expire. Keep getting extra 20 shurikens (SHURI) for each new recruit.`, username, endpoint, hash, hash, refLink, refLink)
    
    status, err := s.Send("dojo@ninja.org", email, subject, body)
    log.Println("Send mail CompleteReferrer status", status, err)
}

func (s MailService) SendFirstBet(email string, username string, hash string) {
    endpoint := utils.GetServerEndpoint()
    refLink := fmt.Sprintf("%s?ref=%s", endpoint, username)
    
    subject := `Your ninja recruit placed a prediction. That's another 20 Shurikens (SHURI) for you.`
    body := fmt.Sprintf(`Hi again,<br/><br/>
        %s<br/><br/>
        Your new recruit placed a prediction using your referral link. 20 shurikens (SHURI) have been added to <a href="%s">your wallet</a>. Boo ya!<br/><br/>
        Your TxHash: <a href="https://etherscan.io/tx/%s">%s</a><br/><br/> 
        Keep spreading the love: <a href="%s">%s</a><br/><br/> 
        20 shurikens (SHURI) for each new recruit.`, username, endpoint, hash, hash, refLink, refLink)

    status, err := s.Send("dojo@ninja.org", email, subject, body)
    log.Println("Send mail FirstBet status", status, err)
}

func (s MailService) SendFirstBetReferrer(email string, username string, hash string) {
    endpoint := utils.GetServerEndpoint()
    refLink := fmt.Sprintf("%s?ref=%s", endpoint, username)
    
    subject := `Your ninja recruit placed a prediction. That's another 20 Shurikens (SHURI) for you.`
    body := fmt.Sprintf(`Hi again,<br/><br/>
        %s<br/><br/>
        Your new recruit placed a prediction using your referral link. 20 shurikens (SHURI) have been added to <a href="%s">your wallet</a>. Boo ya!<br/><br/>
        Your TxHash: <a href="https://etherscan.io/tx/%s">%s</a><br/><br/> 
        Keep spreading the love: <a href="%s">%s</a><br/><br/> 
        20 shurikens (SHURI) for each new recruit.`, username, endpoint, hash, hash, refLink, refLink)

    status, err := s.Send("dojo@ninja.org", email, subject, body)
    log.Println("Send mail FirstBet status", status, err)
}

func (s MailService) SendFreeBet(email string, username string, hash string) {
    endpoint := utils.GetServerEndpoint()
    refLink := fmt.Sprintf("%s?ref=%s", endpoint, username)
    
    subject := `You've got 20 Shurikens (SHURI) from Ninja for placing your free bet`
    body := fmt.Sprintf(`Greatings,<br/><br/>
        %s<br/><br/>
        Welcome to the dojo. 20 shurikens (SHURI) have been added to <a href="%s">your wallet</a>.<br/><br/>
        Your TxHash: <a href="https://etherscan.io/tx/%s">%s</a><br/><br/>
        Use these tokens to slash fees, learn secrets, unlock special treatment.<br/><br/>
        This is your unique referral link: <a href="%s">%s</a><br/><br/> 
        Bring your most trusted friends to the dojo and receive 20 extra shurikens (SHURI) for each new recruit.`, username, endpoint, hash, hash, refLink, refLink)

    status, err := s.Send("dojo@ninja.org", email, subject, body)
    log.Println("Send mail FirstBet status", status, err)
}
