package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Caso decida colocar isso em um container, recomendo que remova o loadDataFromDotEnv e suas dependências.
	// Não se esqueça de incluir as variáveis de ambiente no comando docker run
	loadDataFromDotEnv()

	http.HandleFunc("/generate-token", generateTokenHandler)
	http.ListenAndServe(":8080", nil)
}

func loadDataFromDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if initialToken != correctToken {
		http.Error(w, "Authentication Failed", http.StatusUnauthorized)
		return
	}
	
	// Você pode pegar a credencial da AWS Secret Manager
	pemSecret, err := getAWSSecret()
	
	// Ou de sua variavel de ambiente
	pemSecret := os.Getenv("PEM_SECRET")
	if err != nil {
		http.Error(w, "AWS Process Failed", http.StatusInternalServerError)
	}

	appID := os.Getenv("APP_ID")
	appInstallationID := os.Getenv("APP_INSTALLATION_ID")

	jwtToken, err := generateJWTToken(pemSecret, appID)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	appAccessToken, err := generateAppAccessToken(jwtToken, appInstallationID)
	if err != nil {
		http.Error(w, "Failed to generate app access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": appAccessToken})
}

func getAWSSecret() (string, error) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	secretName := os.Getenv("SECRET_NAME")

	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), 
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	secretString := *result.SecretString
	return secretString, nil
}

func generateJWTToken(pemSecret string, appID string) (string, error) {
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pemSecret))
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 10).Unix(),
		"iss": appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateAppAccessToken(jwtToken string, appInstallationID string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", appInstallationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	token, ok := data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}
