package referralapp_api

import (
    "github.com/jcgarciaram/general-api/apiutils"
    "github.com/tmaiaroto/aegis/lambda"
    "github.com/Sirupsen/logrus"
    "gopkg.in/mgo.v2/bson"
    
    "encoding/json"
    "net/http"
    "strings"
    "strconv"
    "net/url"
    "time"
    "fmt"
)




// GetReferralCode retrieves a particular referral code
func GetReferralCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get storeId from JWT
    var storeId string
    if s, ok := getStoreId(evt); !ok {
        res.StatusCode = strconv.Itoa(http.StatusForbidden)
        ret := struct {
            Message string `json:"message"`
        }{"Something is wrong man"}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    } else {
        storeId = s
    }

    codeId := params.Get("code")
    
    // Verify Store Exists
    if errStr, httpResponse := verifyStoreExists(storeId); httpResponse != 0 {
        res.StatusCode = strconv.Itoa(httpResponse)

        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }

    // Query to get cake
    query := "SELECT * from `referral_code` where `referral_code_id` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_" + storeId
    parameters :=  []interface{}{codeId}
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, parameters, getTotalCount)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }

    // If no code is found, return StatusNoContent
    if len(rowMapSlice) == 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusNoContent)
        res.Body = ""
        return
    }

    // Marshal response and return
    retJson, _ := json.Marshal(rowMapSlice[0])
	res.Body = string(retJson)
    
	res.Headers["Content-Type"] = "application/json"
    
}


// GetReferralCodes retrieves all codes for a single store
func GetReferralCodes(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get storeId from JWT
    var storeId string
    if s, ok := getStoreId(evt); !ok {
        res.StatusCode = strconv.Itoa(http.StatusForbidden)
        ret := struct {
            Message string `json:"message"`
        }{"Something is wrong man"}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    } else {
        storeId = s
    }
    
    // Verify Store Exists
    if errStr, httpResponse := verifyStoreExists(storeId); httpResponse != 0 {
        res.StatusCode = strconv.Itoa(httpResponse)

        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }

    // Query to get cake
    query := "SELECT * from `referral_code`"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_" + storeId
    parameters :=  []interface{}{}
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, parameters, getTotalCount)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }

    // If no cakes are found, return StatusNoContent
    if len(rowMapSlice) == 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusNoContent)
        res.Body = ""
        return
    }

    // Marshal response and return
    retJson, _ := json.Marshal(rowMapSlice)
	res.Body = string(retJson)
    
	res.Headers["Content-Type"] = "application/json"
    
}


// GenerateReferralCode generated a unique referral code
func GenerateReferralCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

    // Get storeId from JWT
    var storeId string
    if s, ok := getStoreId(evt); !ok {
        res.StatusCode = strconv.Itoa(http.StatusForbidden)
        ret := struct {
            Message string `json:"message"`
        }{"Something is wrong man"}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    } else {
        storeId = s
    }
    
    // Verify Store Exists
    if errStr, httpResponse := verifyStoreExists(storeId); httpResponse != 0 {
        res.StatusCode = strconv.Itoa(httpResponse)

        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }
    
    // Read body from request
    var bodyByte []byte
    if tBody, err := apiutils.GetBodyFromEvent(evt); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    } else {
        bodyByte = tBody
    }

    // Struct to unmarshal body of request
    var rc ReferralCode

    
    // Unmarshal body into ReferralCode struct defined above
    if err := json.Unmarshal(bodyByte, &rc); err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error marshaling JSON to ReferralCode struct")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusUnprocessableEntity)
        res.Body = "Error marshaling JSON to ReferralCode struct"
        return
    }
    
    rc.ReferralCodeId = bson.NewObjectId().Hex()
    rc.StoreId = storeId
    rc.generateVerificationCode()
    expirationDate := time.Now().Add(time.Duration((7*24))*time.Hour)
    
    
    for {
        
            
        // Query to insert referral_code
        query := fmt.Sprintf("INSERT INTO `%s`.`referral_code` (`referral_code_id`,`store_id`,`generated_phone`,`expiration_date`,`used_count`,`verification_code`,`verification_code_display`) VALUES (?,?,?,?,?,?,?)",("referralapp_" + storeId))
        
        parameters := []interface{}{
            rc.ReferralCodeId,
            storeId,
            rc.GeneratedPhone,
            expirationDate,
            0,
            rc.VerificationCode,
            rc.VerificationCodeDisplay,
        }
     
        
        // Build query to run in MySQL
        upsertQueries := []apiutils.UpsertQuery{
            {
                Query: query,
                Parameters: parameters,
            },

        }
            
        // Run queries
        getLastInsertId := false
        _, _, errStr, httpResponse := apiutils.RunUpsertQueries(upsertQueries, getLastInsertId)
        if httpResponse != 0 {
            
            if strings.Contains(errStr, "1062") {
                rc.ReferralCodeId = bson.NewObjectId().Hex()
                continue
            }
            
            res.Headers["Content-Type"] = "charset=UTF-8"
            res.StatusCode = strconv.Itoa(httpResponse)
            res.Body = errStr
            return
        }

        
        break
    }
    
    // Generate QR png image
    qrCode, err := rc.generateQRCode()
    if err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error generating QR code")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Store QR to S3
    bucket := "referralapp-qrcodes"
    key := storeId + "/referral/" + rc.ReferralCodeId + ".png"
    if err := apiutils.SaveToS3(bucket, key, qrCode); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error storing QR png to S3")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Make it publicly accessible
    if err := apiutils.GiveS3ObjectPublicRead(bucket, key); err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error storing QR png to S3")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Public URL for object
    url := "https://s3.amazonaws.com/referralapp-qrcodes/" + storeId + "/referral/" + rc.ReferralCodeId + ".png"
    
    var messageBody string
    
    // Shorten URL using Google URL Shortener API
    if shortUrl, err := apiutils.ShortenURL(url); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error shortening URL")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    } else {
    
        messageBody = "Forward this code to your friends and get rewarded when it's used: " + shortUrl
    
    }
    
    // Send SMS Message using AWS SNS
    if err = apiutils.SendAwsSMSMessage(rc.GeneratedPhone, messageBody); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error sending SMS")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Write OK
    res.StatusCode = strconv.Itoa(http.StatusOK)
}



// UseReferralCode verifies whether code is still valid based on expiration date and keeps a count of the times it has been used
func UseReferralCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get storeId from JWT
    var storeId string
    if s, ok := getStoreId(evt); !ok {
        res.StatusCode = strconv.Itoa(http.StatusForbidden)
        ret := struct {
            Message string `json:"message"`
        }{"Something is wrong man"}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    } else {
        storeId = s
    }
    
    codeId := params.Get("code")
    
    // Verify Store Exists
    if errStr, httpResponse := verifyStoreExists(storeId); httpResponse != 0 {
        res.StatusCode = strconv.Itoa(httpResponse)

        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }
    
    
    // Query to run
    query := fmt.Sprintf("UPDATE `%s`.`referral_code` SET `used_count` = `used_count` + 1 WHERE `referral_code_id` = ? AND CURRENT_TIMESTAMP <= `expiration_date`", ("referralapp_" + storeId)) 
    parameters := []interface{}{codeId}
    
    // Build query to run in MySQL
    upsertQueries := []apiutils.UpsertQuery{
        {
            Query: query,
            Parameters: parameters,
        },

    }
    
    // Run queries
    getLastInsertId := false
    _, affectedRows, errStr, httpResponse := apiutils.RunUpsertQueries(upsertQueries, getLastInsertId)
    if httpResponse != 0 {
        
        res.StatusCode = strconv.Itoa(http.StatusBadRequest)
        
        
        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }
    
    if affectedRows == 0 {
        res.StatusCode = strconv.Itoa(http.StatusBadRequest)
        
        
        ret := struct {
            Message string `json:"message"`
        }{"Invalid or expired code"}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }
    
    
    // Query to get referral code
    query = "SELECT * from `referral_code` where `referral_code_id` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_" + storeId
    parameters =  []interface{}{codeId}
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, parameters, getTotalCount)
    if httpResponse != 0 {
        res.StatusCode = strconv.Itoa(http.StatusBadRequest)
        
        
        ret := struct {
            Message string `json:"message"`
        }{errStr}
        
        // Marshal response and return
        retJson, _ := json.Marshal(ret)
        res.Body = string(retJson)
        
        res.Headers["Content-Type"] = "application/json"
        return
    }
    
    // Get phone that generated original code to send originator reward
    originatorPhone := rowMapSlice[0]["generated_phone"].(string)
    
    // Create new variable
    var oc OriginatorCode
    oc.Phone = originatorPhone
    oc.StoreId = storeId
    
    // Build all originator fields and store in DB
    if err := oc.buildOriginatorCode(); err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error generating QR code")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Generate QR png image
    qrCode, err := oc.generateQRCode()
    if err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error generating QR code")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Store QR to S3
    bucket := "referralapp-qrcodes"
    key := storeId + "/originator/" + oc.OriginatorCodeId + ".png"
    if err := apiutils.SaveToS3(bucket, key, qrCode); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error storing QR png to S3")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }

    // Make it publicly accessible
    if err := apiutils.GiveS3ObjectPublicRead(bucket, key); err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error storing QR png to S3")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Public URL for object
    url := "https://s3.amazonaws.com/referralapp-qrcodes/" + storeId + "/originator/" + oc.OriginatorCodeId + ".png"
    
    var messageBody string
    
    // Shorten URL using Google URL Shortener API
    if shortUrl, err := apiutils.ShortenURL(url); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error shortening URL")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    } else {
    
        messageBody = "Someone has used your QR code! Show this and get 10% off your purchase! " + shortUrl
    
    }
    
    // Send SMS Message using AWS SNS
    if err = apiutils.SendAwsSMSMessage(oc.Phone, messageBody); err != nil {
        
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error sending SMS")  
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Write OK
    res.StatusCode = strconv.Itoa(http.StatusOK)

}