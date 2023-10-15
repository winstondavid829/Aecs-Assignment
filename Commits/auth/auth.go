package auth

import (
	"AECS_Assignment/data"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)

type UserToken struct {
	Accessuuid string  `json:"access_uuid"`
	Aud        string  `json:"aud"`
	Exp        float64 `json:"exp"`
	Iss        string  `json:"iss"`
	Userid     string  `json:"userid"`
	UserToken  string  `json:"usertoken"`
}

func init() {
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
}

// Initialize Redis Start
var redisclient *redis.Client

func CreateToken(userId string) (*data.TokenDetails, error) {
	td := &data.TokenDetails{}
	//td.AtExpires = time.Now().Add(time.Hour * 24).Unix() // commented on 13-01-2021
	td.AtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // New code - 13-01-2021 - Token will expire in 7 days
	//td.AtExpires = time.Now().Add(time.Minute * 5).Unix() // New code - for testing to check access token expire in 5 minutes
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

func CreateTokenforResetPassword(userId string, mysecreat string, useremail string, orgID string) (*data.TokenDetails, error) {
	td := &data.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 10).Unix()
	td.AccessUuid = uuid.NewV4().String()

	//td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // I had commented on 16-12-2020
	td.RtExpires = time.Now().Add(time.Minute * 5).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	//os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires
	atClaims["email"] = useremail
	atClaims["sessionid"] = orgID // we just cheet viewer
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(mysecreat))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

func CreateAuth(userid string, td *data.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := redisclient.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := redisclient.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// Extract Token from Header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	//println("Received Primary Token", strArr[1])
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func DecodeBase64String(tokenString string) (string, bool) {
	//tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjMwZTUwZWRkLTFhOWItNDMxYy1iZWM1LTNjYzdmNTc3YjYwMCIsImF1ZCI6Im9yZ3JlZ2lzdHJ5LWRldi5lbmRwb2ludHMubmV1cmFsd2F5cy5jbG91ZC5nb29nIiwiZXhwIjoxNTkwOTE4MTg5LCJpc3MiOiJ3aXRtZWctZGV2QG5ldXJhbHdheXMuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJ1c2VyaWQiOiI1ZWFkNDNlMzBjNjBiOGRlYzc4MjA5OGQiLCJ1c2VydG9rZW4iOiJleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05mZFhWcFpDSTZJakJtWWpNMlpUVXlMVFV4TTJZdE5HWmhZeTA1TVRBekxUVmpaRE0wWXpObFpEazBOeUlzSW1GMWRHaHZjbWw2WldRaU9uUnlkV1VzSW1WNGNDSTZNVFU1TURreE5EWXhOQ3dpZFhObGNtbGtJam9pTldWaFpEUXpaVE13WXpZd1lqaGtaV00zT0RJd09UaGtJbjAuWFp3MzI2OWFhRTJlMW9Pdi1KSVNMZm8tU2ZDZ2swTlVRdzZDandjcnZYTSJ9.ewIH04WabPA-XXLQoYzBcYPHKyQqSvzFAAkPZw0l3q6VaCXlz6o4JuqKNGGvcwIOKSeo8pLcSaUQTRgjUUFGoBxA_Wb5htIG_y4nksgb8K4BcChdQt4UDtLUhofF6lwhZ5YdQMqNrQL5Z-0YpiHAQUGciKgL5wX-Jw5L3DpxkVLcoPIaRRXZzUy4KdktMJ4AmsrXXRSo0WDBrfVdp-f05LFtoKTosqXN4w4SNdJKBmTaCYyS30qFHkCLjdL0HV0a5xpMhimtacAGiuV5Phkbc6R1218SZEfTqnosyjebedd20KNEIPBh11fUxKK99NItOKhflgDV_LcmQzY-geT7iA"

	strArr := strings.Split(tokenString, ".")

	//fmt.Println("source B64:", strArr)
	//byteString, err := b64.URLEncoding.DecodeString(strArr[1])
	byteString, err := b64.StdEncoding.WithPadding(b64.NoPadding).DecodeString(strArr[1])

	if err != nil {
		fmt.Println("Error in decoing 64BasedUrl:", err)
		return "", false
	}
	//println("Payload is : ", fmt.Sprintf("%s", byteString))

	userToken := UserToken{}
	if err := json.Unmarshal(byteString, &userToken); err != nil {
		//panic(err)
		println("Error in Unmarshal userToken")
		return "", false
	}

	//fmt.Println(userToken)
	return userToken.UserToken, true
}

func DecodeBase64UserTokenString(tokenString string) (data.CustomerToken, bool) {
	//tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjMwZTUwZWRkLTFhOWItNDMxYy1iZWM1LTNjYzdmNTc3YjYwMCIsImF1ZCI6Im9yZ3JlZ2lzdHJ5LWRldi5lbmRwb2ludHMubmV1cmFsd2F5cy5jbG91ZC5nb29nIiwiZXhwIjoxNTkwOTE4MTg5LCJpc3MiOiJ3aXRtZWctZGV2QG5ldXJhbHdheXMuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJ1c2VyaWQiOiI1ZWFkNDNlMzBjNjBiOGRlYzc4MjA5OGQiLCJ1c2VydG9rZW4iOiJleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05mZFhWcFpDSTZJakJtWWpNMlpUVXlMVFV4TTJZdE5HWmhZeTA1TVRBekxUVmpaRE0wWXpObFpEazBOeUlzSW1GMWRHaHZjbWw2WldRaU9uUnlkV1VzSW1WNGNDSTZNVFU1TURreE5EWXhOQ3dpZFhObGNtbGtJam9pTldWaFpEUXpaVE13WXpZd1lqaGtaV00zT0RJd09UaGtJbjAuWFp3MzI2OWFhRTJlMW9Pdi1KSVNMZm8tU2ZDZ2swTlVRdzZDandjcnZYTSJ9.ewIH04WabPA-XXLQoYzBcYPHKyQqSvzFAAkPZw0l3q6VaCXlz6o4JuqKNGGvcwIOKSeo8pLcSaUQTRgjUUFGoBxA_Wb5htIG_y4nksgb8K4BcChdQt4UDtLUhofF6lwhZ5YdQMqNrQL5Z-0YpiHAQUGciKgL5wX-Jw5L3DpxkVLcoPIaRRXZzUy4KdktMJ4AmsrXXRSo0WDBrfVdp-f05LFtoKTosqXN4w4SNdJKBmTaCYyS30qFHkCLjdL0HV0a5xpMhimtacAGiuV5Phkbc6R1218SZEfTqnosyjebedd20KNEIPBh11fUxKK99NItOKhflgDV_LcmQzY-geT7iA"

	strArr := strings.Split(tokenString, ".")

	//fmt.Println("source B64:", strArr)
	//byteString, err := b64.URLEncoding.DecodeString(strArr[1])
	byteString, err := b64.StdEncoding.WithPadding(b64.NoPadding).DecodeString(strArr[1])

	if err != nil {
		fmt.Println("Error in decoing 64BasedUrl:", err)
		return data.CustomerToken{}, false
	}
	//println("Payload is : ", fmt.Sprintf("%s", byteString))

	customerToken := data.CustomerToken{}
	if err := json.Unmarshal(byteString, &customerToken); err != nil {
		//panic(err)
		println("Error in Unmarshal userToken")
		return data.CustomerToken{}, false
	}

	//fmt.Println(userToken)
	return customerToken, true
}

// Verify Token - signing method HS256
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r) // Result is checking in Redis
	//println("UserToken is : ", tokenString)
	//Here we decode Cloud Token and get our user token from payload Start
	tokenString, status := DecodeBase64String(tokenString) //You can comment for local testing. But Uncomment if use in Cloud
	if status == false {
		println("Error in Decode UserToken")
		return nil, fmt.Errorf("Unexpected Decode Error in received token")
	}
	//println("Token after decode from Primary Token : ", tokenString)
	//Here we decode Cloud Token and get our user token from payload End

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//println("Current ACCESS_SECREAT : ", os.Getenv("ACCESS_SECRET"))

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		println("Token Verified failed : ", fmt.Sprintf("%s", err))
		return nil, err
	}
	println("Token Verified")
	return token, nil
}

type KeyResetCustomerPassword struct{}

func VerifyTokenES256(r *http.Request) (*jwt.Token, error) {
	//For RS256 alg start
	var privatePem = os.Getenv("ACCESS_PRIVATE")
	block, _ := pem.Decode([]byte(privatePem)) //Here we need to pass Public Key Pem, Not Private Pem
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		//log.Fatal(err)
		println("VerifyTokenES256 : Key Arror", err)
		return nil, err
	}
	//For RS256 alg end

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		println("Error in VerifyTokenES256 : ", fmt.Sprintf("%v", err))
		return nil, err
	}
	return token, nil
}

// TokenValid Checks Token is valid
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// Extract Meta Data from Token for HS256 for Redis Check
func ExtractTokenMetadata(r *http.Request) (*data.AccessDetails, error) {
	token, err := VerifyToken(r)
	//token, err := VerifyTokenES256(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		//userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		userID := claims["userid"].(string)
		if err != nil {
			return nil, err
		}
		return &data.AccessDetails{
			AccessUuid: accessUUID,
			UserId:     userID,
		}, nil
	}
	return nil, err
}

// Extract Meta Data from Token for HS256, to check the exp time
func ExtractTokenMetadataExpTime(r *http.Request) (bool, error) {
	token, err := VerifyToken(r)

	if err != nil {
		println("ExtractTokenMetadataExpTime->VerifyToken :", fmt.Sprintf("%s", err))
		return false, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if getTokenRemainingValidity(claims["exp"]) > 0 {
			return true, nil
		}
		return false, nil

	}
	println("ExtractTokenMetadataExpTime->End :", fmt.Sprintf("%s", err))
	return false, err
}

func getTokenRemainingValidity(timestamp interface{}) int {
	var expireOffset = 3600 * 24
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())
		//println("getTokenRemainingValidity->Remainder :", remainder)
		if remainder > 0 {
			//println("getTokenRemainingValidity->Remainder :", remainder)
			return int(remainder.Seconds() + float64(expireOffset))
		}
	}
	//println("getTokenRemainingValidity->expireOffset :", expireOffset)
	return 0
}
func getTokenRemainingValidityforPasswordReset(timestamp interface{}) int {
	var expireOffset = 60 * 10
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())
		//println("getTokenRemainingValidity->Remainder :", remainder)
		if remainder > 0 {
			//println("getTokenRemainingValidity->Remainder :", remainder)
			return int(remainder.Seconds() + float64(expireOffset))
		}
	}
	//println("getTokenRemainingValidity->expireOffset :", expireOffset)
	return 0
}
func FetchAuth(authD *data.AccessDetails) (string, error) {
	userid, err := redisclient.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	//userID, _ := strconv.ParseUint(userid, 10, 64)
	userID := userid
	return userID, nil
}

// ValidateUserInHeader checks that user request token is genue and which is in redis
func ValidateUserInHeader(r *http.Request) string {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		println("unauthorized user request - Extract Token: %v", err)
		return ""
	}
	userID, err := FetchAuth(tokenAuth)
	if err != nil {
		println("unauthorized user request - FetchAuth: %v", err)
		return ""
	}
	return userID
}

// ValidateUserTokenInHeader is used to validate ourself using exp time in the token, It will not check in redis.
func ValidateUserTokenInHeader(r *http.Request) bool {
	//tokenAuth, err := ExtractTokenMetadata(r)//This is used with redis
	tokenAuth, err := ExtractTokenMetadataExpTime(r) //This is used ourself by checking exp time
	if err != nil {
		println("unauthorized user request - Extract Token: %v", err)
		return false
	}
	if tokenAuth {

		return true
	}
	return false
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := redisclient.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func CreateRSA256Token(audience, issuer, privatekey, userToken string) (string, error) {
	var err error
	//Setting Environmental Data at start up
	//data.IsProductionEnvironment(false)
	var privatePem = privatekey
	block, _ := pem.Decode([]byte(privatePem))
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	atClaims := jwt.MapClaims{}
	atClaims["iss"] = issuer
	atClaims["aud"] = audience
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	atClaims["userToken"] = userToken
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	token, err := at.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}
