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
    
    // Get parameters from URL request
    storeId := params.Get("store")
    codeId := params.Get("code")

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

    // If no cakes are found, return StatusNoContent
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
    
    // Get parameters from URL request
    storeId := params.Get("store")

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


// PostReferralCode generated a unique referral code
func PostReferralCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

    // Get parameters from URL request
    storeId := params.Get("store")
    
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
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    // Store QR to S3
    bucket := "referralapp-qrcodes"
    key := storeId + "/" + rc.ReferralCodeId + ".png"
    if err := apiutils.SaveToS3(bucket, key, qrCode); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    url, err := apiutils.GetDownloadURL(bucket, key)
    if err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    
    /*
    messageBody := "Forward this code to your friends and get rewarded when it's used!"
    if err := apiutils.SendMMSMessage(rc.GeneratedPhone, messageBody, url); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    */
    
    var messageBody string
    
    if shortUrl, err := apiutils.ShortenURL(url); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    } else {
    
        messageBody = "Forward this code to your friends and get rewarded when it's used: " + shortUrl
    
    }
    
    
    if err = apiutils.SendSMSMessage(rc.GeneratedPhone, messageBody); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = err.Error()
        return
    }
    

    res.StatusCode = strconv.Itoa(http.StatusOK)
}



// UseReferralCode verifies whether code is still valid based on expiration date and keeps a count of the times it has been used
func UseReferralCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get parameters from URL request
    storeId := params.Get("store")
    codeId := params.Get("code")
    
    
    
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
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }
    
    if affectedRows == 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusBadRequest)
        res.Body = "Invalid or expired code"
        return
    }
    
    res.Headers["Content-Type"] = "charset=UTF-8"
    res.StatusCode = strconv.Itoa(http.StatusOK)
    res.Body = "Your code was used successfully!"

}