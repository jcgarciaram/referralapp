package referralapp_api

import (
    
    "github.com/jcgarciaram/general-api/apiutils"
    
    // "github.com/aws/aws-sdk-go/aws/session"
    "github.com/tmaiaroto/aegis/lambda"
    // "github.com/aws/aws-sdk-go/aws"
    "github.com/Sirupsen/logrus"
    // "github.com/guregu/dynamo"
    "gopkg.in/mgo.v2/bson"
    
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "net/url"
    // "time"
    "log"
    "fmt"
)

// PostStore
func PostStore(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

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
    
 
    // Struct to unmarshal body of request into
    var s Store
    

    // Unmarshal body into storeConfig struct defined above
    if err := json.Unmarshal(bodyByte, &s); err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error marshaling JSON to storeConfig struct")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusUnprocessableEntity)
        res.Body = "Error marshaling JSON to storeConfig struct"
        return
    }
    
    s.StoreId = bson.NewObjectId().Hex()
    

    
    for {
        
            
        // Query to insert cake
        query := 
            "INSERT INTO `referralapp_master`.`store` " + 
            "(`store_id`," +
            "`store_name`," +
            "`store_type`," +
            "`store_email`," +
            "`contact_user_id`," +
            "`created_by`," +
            "`last_updated_by`) " +
            "VALUES (?,?,?,?,?,?,?)"
        
        parameters := []interface{}{
            s.StoreId,
            s.StoreName,
            s.StoreType,
            s.StoreEmail,
            s.ContactUserId,
            UserEmail,
            UserEmail,
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
                s.StoreId = bson.NewObjectId().Hex()
                continue
            }
            
            res.Headers["Content-Type"] = "charset=UTF-8"
            res.StatusCode = strconv.Itoa(httpResponse)
            res.Body = errStr
            return
        }

        
        break
    }
    
    
    // Get all queries to create schema in MySQL database
    createSchemaQueries := getCreateSchemaQueries(s.StoreId)
    
    log.Println("About to run all create schema queries")
    
    // Run create schema / table queries
    getLastInsertId := false
    _, _, errStr, httpResponse := apiutils.RunUpsertQueries(createSchemaQueries, getLastInsertId)
    if httpResponse != 0 {
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = fmt.Sprintf(errStr)
        return
        
    }
    
    retJson, _ := json.Marshal(s)
	res.Body = string(retJson)

    res.Headers["Content-Type"] = "application/json"
    
    
}


// GetStore
func GetStore(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    storeId := params.Get("store")


    // Query to get user
    query := "SELECT * from `store` where `store_id` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_master"
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, []interface{}{storeId}, getTotalCount)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }

    // If no users are found, return StatusNoContent
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