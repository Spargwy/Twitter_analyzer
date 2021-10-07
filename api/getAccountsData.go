package api

import (
	"dev-team/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RateLimitError struct{}

func (m *RateLimitError) Error() string { return "Rate limit exceeded" }

func RateLimit() error {
	return &RateLimitError{}
}

func GetAllAccountsData(usernames []string) (accountsData []storage.Account, err error) {
	if err != nil {
		return accountsData, err
	}
	for i := range usernames {
		accountData, err := GetMainAccountData(usernames[i])
		if err != nil {
			return accountsData, err
		}
		accountsData = append(accountsData, accountData)
	}
	for i := range accountsData {
		GetAdditionAccountData(&accountsData[i], "followers", "")
		GetAdditionAccountData(&accountsData[i], "following", "")
		GetAdditionAccountData(&accountsData[i], "tweets", "")

	}

	return accountsData, nil
}

func GetMainAccountData(username string) (accountData storage.Account, err error) {
	url := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}
	bearerToken := fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN"))
	req.Header.Add("Authorization", bearerToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &accountData)

	if err != nil {
		log.Print("ERROR IN UNMARSHALL: ", err)
		return accountData, err
	}
	return accountData, nil
}

func GetAdditionAccountData(user *storage.Account, endpoint string, params string) error {

	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/%s?max_results=100%s", user.Data.ID, endpoint, params)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return err
	}
	bearerToken := fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN"))
	req.Header.Add("Authorization", bearerToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if string(body) == "Rate limit exceeded" {
		err := RateLimit()
		if err != nil {
			return err
		}
	}
	switch endpoint {
	case "followers":
		var followers storage.Followers
		err = json.Unmarshal(body, &followers)
		if err != nil {
			log.Print("ERROR IN UNMARSHALL FOLLOWERS: ", err)
			return err
		}
		user.Followers += followers.Meta.ResultCount
		if followers.Meta.NextToken != "" {
			params = fmt.Sprintf("&pagination_token=%s", followers.Meta.NextToken)
			err = GetAdditionAccountData(user, endpoint, params)
			if err != nil {
				return err
			}
		}
	case "following":
		var following storage.Following
		err = json.Unmarshal(body, &following)
		if err != nil {
			log.Print("ERROR IN UNMARSHALL FOLLOWING: ", err)
			return err
		}
		user.Following += following.Meta.ResultCount
		if following.Meta.NextToken != "" {
			params = fmt.Sprintf("&pagination_token=%s", following.Meta.NextToken)
			err = GetAdditionAccountData(user, endpoint, params)
			if err != nil {
				return err
			}
		}
	case "tweets":
		var tweets storage.Tweets
		err = json.Unmarshal(body, &tweets)
		if err != nil {
			log.Print("ERROR IN UNMARSHALL TWEETS: ", err)
			return err
		}
		user.Tweets += tweets.Meta.ResultCount
		if tweets.Meta.NextToken != "" {
			params = fmt.Sprintf("&pagination_token=%s", tweets.Meta.NextToken)
			err = GetAdditionAccountData(user, endpoint, params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
