package referralapp_api

import (
    
    "github.com/jcgarciaram/general-api/apiutils"
    // "github.com/aws/aws-sdk-go/aws/session"
    "github.com/tmaiaroto/aegis/lambda"
    // "github.com/aws/aws-sdk-go/aws"
    "golang.org/x/crypto/bcrypt"
    "github.com/Sirupsen/logrus"
    // "github.com/guregu/dynamo"
    
    "encoding/json"
    "net/url"
    "strconv"
    // "time"
    "fmt"
)



// PostUser
func PostUser(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    
    // Read body from request
    var bodyByte []byte
    if tBody, err := apiutils.GetBodyFromEvent(evt); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusInternalServerError
        res.Body = err.Error()
        return
    } else {
        bodyByte = tBody
    }
    

    
    // Struct to unmarshal body of request into
    var u User
    

    // Unmarshal body into userConfig struct defined above
    if err := json.Unmarshal(bodyByte, &u); err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error marshaling JSON to userConfig struct")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusUnprocessableEntity
        res.Body = "Error marshaling JSON to userConfig struct"
        return
    }

    // Hash password
    hashByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 6)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error hashing password")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusInternalServerError
        res.Body = fmt.Sprintf("Error hashing password")
        return
    }
    
    // We will only store the hashed password
    u.Password = string(hashByte)
    
    
    
    // Query to insert cake
    query := 
        "INSERT INTO `referralapp_master`.`user` " + 
        "(`user_email`," +
        "`password`," +
        "`store_id`," +
        "`role`," +
        "`first_name`," +
        "`middle_name`," +
        "`last_name`," +
        "`second_last_name`) " +
        "VALUES (?,?,?,?,?,?,?,?)"
    
    parameters := []interface{}{
        u.UserEmail,
        u.Password,
        u.Store,
        u.Role,
        u.FirstName,
        u.MiddleName,
        u.LastName,
        u.SecondLastName,
    }
 
    
    // Build query to run in MySQL
    upsertQueries := []apiutils.UpsertQuery{
        {
            Query: query,
            Parameters: parameters,
        },

    }
    
    // Run queries
    getLastInsertId := true
    lastInsertId, _, errStr, httpResponse := apiutils.RunUpsertQueries(upsertQueries, getLastInsertId)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }
    res.Body = strconv.Itoa(lastInsertId)
    res.StatusCode = StatusOK

}


// GetUser
func GetUser(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
    
    // Get parameters from URL request
    userEmail := params.Get("useremail")

    // Query to get user
    query := "SELECT * from `user` where `user_email` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_master"
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, []interface{}{userEmail}, getTotalCount)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }

    // If no users are found, return StatusNoContent
    if len(rowMapSlice) == 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusNoContent
        res.Body = ""
        return
    }

    // Marshal response and return
    retJson, _ := json.Marshal(rowMapSlice)
	res.Body = string(retJson)
    
	res.Headers["Content-Type"] = "application/json"
    
}


// LogInUser
func LogInUser(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

    
    // Read body from request
    var bodyByte []byte
    if tBody, err := apiutils.GetBodyFromEvent(evt); err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusInternalServerError
        res.Body = err.Error()
        return
    } else {
        bodyByte = tBody
    }
    

    
    // Struct to unmarshal body of request into
    var u User
    

    // Unmarshal body into userConfig struct defined above
    if err := json.Unmarshal(bodyByte, &u); err != nil {
    
        logrus.WithFields(logrus.Fields{
            "err": err,
        }).Warn("Error marshaling JSON to userConfig struct")   
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusUnprocessableEntity
        res.Body = "Error marshaling JSON to userConfig struct"
        return
    }

    // Query to get user
    query := "SELECT password from `user` where `user_email` = ?"
    
    // Run query from MySQL
    getTotalCount := false
    schema := "referralapp_master"
    rowMapSlice, _, errStr, httpResponse := apiutils.RunSelectQuery(schema, query, []interface{}{u.UserEmail}, getTotalCount)
    if httpResponse != 0 {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = strconv.Itoa(httpResponse)
        res.Body = errStr
        return
    }

    // If no users are found, return StatusNoContent
    if len(rowMapSlice) == 0 {
        
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusNoContent
        res.Body = "No content"
        return
    }
    
    storedPassword := rowMapSlice[0]["password"].(string)
    
    // Verify password is correct. If not, return error
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password)); err != nil {
    
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusNoContent
        res.Body = "Invalid username/password"
        return
    }
    
    token, err := apiutils.GenerateJWT(u.UserEmail)
    if err != nil {
        res.Headers["Content-Type"] = "charset=UTF-8"
        res.StatusCode = StatusInternalServerError
        res.Body = "Error generating JSON Web Token"
        return
    }
    
    res.Headers["Set-Cookie"] = "access-token=" + token
    res.StatusCode = StatusOK
    
}

