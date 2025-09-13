package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const userAgent = "Sikx2i/0.1.0 (https://www.example.com/contact; zezima; to check if user exists or not)"

type ipData struct {
	City    string `json:"city"`
	Country string `json:"country"`
	ISP     string `json:"isp"`
}

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	URL   string `json:"url"`
}

type GitLabUser struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	WebURL   string `json:"web_url"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	Number    string `json:"number"`
}

type Database struct {
	Users []User `json:"users"`
}

type RedditResponse struct {
	Available bool `json:"available"`
}

func main() {
	fmt.Println("Welcome to passive")

	fullName := flag.String("fn", "", "Search with full-name")
	ipAddress := flag.String("ip", "", "Search with IP")
	username := flag.String("u", "", "Search with username")
	help := flag.Bool("help", false, "Display help message")

	flag.Parse()

	if *help {
		fmt.Println("OPTIONS")
		flag.PrintDefaults()
		return
	}

	if *fullName != "" {
		fmt.Println("Processing full name search...")
		newName := strings.Title(strings.ToLower(*fullName))
		searchByFullName(newName)
	} else if *ipAddress != "" {
		fmt.Println("Processing IP search...")
		searchByIp(*ipAddress)
	} else if *username != "" {
		fmt.Println("Processing username search...")
		results := searchByUsername(*username)
		createResultTXT(results)
	} else {
		fmt.Println("No valid input detected. Use --help for more info")
	}
}

func searchByUsername(username string) (result string) {
	hasYoutube := youtubeSearch(username)
	hasTikTok := tikTokSearch(username)
	hasReddit := redditSearch(username)
	hasGithub := githubSearch(username)
	hasGitlab := gitLabSearch(username)

	result = fmt.Sprintf(
		"Search results for username '%s':\n- Reddit: %v\n- YouTube: %v\n- TikTok: %v\n- GitLab: %v\n- GitHub: %v\n",
		username, hasReddit, hasYoutube, hasTikTok, hasGitlab, hasGithub,
	)
	return result
}

func searchByIp(ip string) {
	url := "http://ip-api.com/json/" + ip
	response, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var ipInfo ipData
	errs := json.Unmarshal(responseData, &ipInfo)
	if errs != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}
	result := "ISP: " + ipInfo.ISP + "\nCity: " + ipInfo.City + "\nCountry: " + ipInfo.Country
	if ipInfo.ISP == "" {
		log.Println("No info found with that IP")
		return
	}
	createResultTXT(result)
}

func searchByFullName(fullname string) {
	splittedName := strings.Split(fullname, " ")

	if len(splittedName) <= 1 {
		log.Println("Error: Please write both first and last name")
		return
	}
	FirstName := splittedName[0]
	LastName := splittedName[1]

	database, err := readDatabase()
	if err != nil {
		log.Println(err)
		return
	}

	for _, user := range database.Users {
		if user.FirstName == FirstName && user.LastName == LastName {
			result := fmt.Sprintf(
				"First name: %s\nLast name: %s\nAddress: %s\nNumber: %s", user.FirstName, user.LastName, user.Address, user.Number,
			)
			createResultTXT(result)
			return
		}
	}
	fmt.Printf("No data with name %s.", fullname)
}

func createResultTXT(result string) {
	files, _ := os.ReadDir("results")
	numberOfFiles := strconv.Itoa(len(files))

	newFile, err := os.Create("results/result" + numberOfFiles + ".txt")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = newFile.WriteString(result)
	if err != nil {
		log.Println(err)
		return
	}

	err = newFile.Close()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
}

func readDatabase() (Database, error) {
	var db Database
	file, err := os.ReadFile("dummy_database.json")
	if err != nil {
		fmt.Println(err)
		return db, err
	}
	err = json.Unmarshal(file, &db)
	return db, err
}

func githubSearch(rawusername string) (hasGithub bool) {
	username := removeSymbol(rawusername)
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		hasGithub = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get response in github search: %v", resp.Status)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode response in githubSearch: %v", err)
	}

	if user.Login != "" {
		hasGithub = true
	} else {
		hasGithub = false
	}
	return hasGithub
}

func gitLabSearch(rawusername string) (hasGitlab bool) {
	username := removeSymbol(rawusername)
	url := fmt.Sprintf("https://gitlab.com/api/v4/users?username=%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get response in gitlab search: %v", resp.Status)
	}

	var users []GitLabUser
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		log.Printf("Failed to decode response in gitlabSearch: %v", err)
	}

	if len(users) > 0 {
		hasGitlab = true
	} else {
		hasGitlab = false
	}
	return hasGitlab
}

func removeSymbol(rawUsername string) (usernameWithoutSymbol string) {
	if strings.Contains(rawUsername, "@") {
		cutUsername := strings.Split(rawUsername, "@")
		return cutUsername[1]
	}
	return rawUsername
}

func redditSearch(rawusername string) (hasReddit bool) {
	username := removeSymbol(rawusername)
	url := fmt.Sprintf("https://www.reddit.com/api/username_available.json?user=%s", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get response in Redditsearch: %v", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&hasReddit); err != nil {
		hasReddit = true
		log.Print("Error occurred in redditSearch:", err)
		log.Printf("Probably invalid input for reddit user search!")
	}
	return !hasReddit
}

func tikTokSearch(rawusername string) (hasTikTok bool) {
	_, username, _ := strings.Cut(rawusername, "@")
	tiktok := "http://www.tiktok.com/@"
	response, err := http.Get(tiktok + username)
	if err != nil {
		log.Println(err.Error())
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	if (len(string(responseData)) > 186500 && len(string(responseData)) < 187700) || (len(string(responseData)) < 184500) {
		hasTikTok = false
	} else {
		hasTikTok = true
	}
	return hasTikTok
}

func youtubeSearch(rawusername string) (hasYoutube bool) {
	APIkey := "AIzaSyD0kppK-iB0j_1Jgu_n2QWM9J5CRnoyCW4"

	_, username, _ := strings.Cut(rawusername, "@")
	youtube := "https://youtube.googleapis.com/youtube/v3/channels?part=snippet,contentDetails,statistics&forUsername=" + username + "&key=" + APIkey
	response, err := http.Get(youtube)
	if err != nil {
		log.Println(err.Error())
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	if len(string(responseData)) > 350 {
		hasYoutube = true
	} else {
		hasYoutube = false
	}
	return hasYoutube
}
