package referralapp_api

import (
    "github.com/jcgarciaram/general-api/apiutils"
    "github.com/tmaiaroto/aegis/lambda"
    // "github.com/Sirupsen/logrus"
    // "gopkg.in/mgo.v2/bson"
    
    "encoding/json"
    "net/http"
    // "strings"
    "strconv"
    "net/url"
    // "time"
    "fmt"
)




// GetOriginatorCode retrieves a particular originator code
func GetOriginatorCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get parameters from URL request
    storeId := params.Get("store")
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
    
    // Query to get originator_code
    query := "SELECT * from `originator_code` where `originator_code_id` = ?"
    
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

    // If no codes are found, return StatusNoContent
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


// GetOriginatorCodes retrieves all codes for a single store
func GetOriginatorCodes(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get parameters from URL request
    storeId := params.Get("store")
    
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

    // Query to get originator_code
    query := "SELECT * from `originator_code`"
    
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

    // If no originator_code are found, return StatusNoContent
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


// UseOriginatorCode verifies whether code is still valid based on expiration date and keeps a count of the times it has been used
func UseOriginatorCode(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get parameters from URL request
    storeId := params.Get("store")
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
    query := fmt.Sprintf("UPDATE `%s`.`originator_code` SET `used_count` = `used_count` + 1 WHERE `originator_code_id` = ? AND CURRENT_TIMESTAMP <= `expiration_date`", ("referralapp_" + storeId)) 
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
    
    // If no affected rows, it means the code did not exist
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
    
    ret := struct {
        Message string `json:"message"`
    }{"Your code was used successfully!"}
    
    // Marshal response and return
    retJson, _ := json.Marshal(ret)
	res.Body = string(retJson)
    
	res.Headers["Content-Type"] = "application/json"

}