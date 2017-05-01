package referralapp_api

import (
    
    "github.com/jcgarciaram/general-api/apiutils"
    
    // "github.com/aws/aws-sdk-go/aws/session"
    "github.com/tmaiaroto/aegis/lambda"
    // "github.com/aws/aws-sdk-go/aws"
    "golang.org/x/crypto/bcrypt"
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

// CreateStoreAccount
func CreateStoreAccount(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

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
    
    // Hash password
    hashByte, err := bcrypt.GenerateFromPassword([]byte(s.Password), 6)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error hashing password")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = fmt.Sprintf("Error hashing password")
        return
    }
    
    // We will only store the hashed password
    s.Password = string(hashByte)
    

    
    for {
        
            
        // Query to insert cake
        query := 
            "INSERT INTO `referralapp_master`.`store` " + 
            "(`store_id`," +
            "`store_name`," +
            "`store_type`," +
            "`store_email`," +
            "`password`)" +
            "VALUES (?,?,?,?,?)"
        
        parameters := []interface{}{
            s.StoreId,
            s.StoreName,
            s.StoreType,
            s.StoreEmail,
            s.Password,
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



// LogInStore
func LogInStore(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

    
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

    // Query to get store
    query := "SELECT * from `store` where `store_email` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_master"
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, []interface{}{s.StoreEmail}, getTotalCount)
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
        res.Body = "No content"
        return
    }
    
    storedPassword := rowMapSlice[0]["password"].(string)
    storeId := rowMapSlice[0]["store_id"].(string)
    
    // Verify password is correct. If not, return error
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(s.Password)); err != nil {
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusNoContent)
        res.Body = "Invalid username/password"
        return
    }
    
    
    js := JWTStruct{s.StoreEmail, storeId}
    
    token, err := apiutils.GenerateJWT(js)
    if err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(http.StatusInternalServerError)
        res.Body = "Error generating JSON Web Token"
        return
    }
    
    res.Headers["Set-Cookie"] = apiutils.GenerateCookieToken(token)
    res.StatusCode = strconv.Itoa(http.StatusOK)
    
}