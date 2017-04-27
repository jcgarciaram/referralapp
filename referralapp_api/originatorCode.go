package referralapp_api

import (
    "fmt"
    "time"
    "errors"
    "strconv"
    "strings"
    "math/rand"
    "gopkg.in/mgo.v2/bson"
    
    qrcode "github.com/skip2/go-qrcode"
    "github.com/jcgarciaram/general-api/apiutils"
)

type OriginatorCode struct {
    OriginatorCodePk            int                 `json:"originator_code_pk"`
    OriginatorCodeId            string              `json:"originator_code_id"`
    StoreId                     string              `json:"store_id"`
    Phone                       string              `json:"phone"`
    ExpirationDate              apiutils.CustomTime `json:"expiration_date"`
    UsedCount                   int                 `json:"used_count"`
    VerificationCode            int                 `json:"verification_code"`
    VerificationCodeDisplay     string              `json:"verification_code_display"`
}

func (oc *OriginatorCode) generateVerificationCode() {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	
	n1 := r1.Intn(100)
	n2 := r1.Intn(100)
	n3 := r1.Intn(100)

	display := fmt.Sprintf("%02d %02d %02d", n1, n2, n3)
	actualInt := fmt.Sprintf("%02d%02d%02d", n1, n2, n3)
    
    oc.VerificationCode, _ = strconv.Atoi(actualInt)
    oc.VerificationCodeDisplay = display
}

func (oc *OriginatorCode) generateQRCode() ([]byte, error) {
    
    url := "2 " + oc.StoreId + " " + oc.OriginatorCodeId

    return qrcode.Encode(url, qrcode.Medium, 256)
}



// buildOriginatorCode stores originator code in MySQL
func (oc *OriginatorCode) buildOriginatorCode() error {

    oc.OriginatorCodeId = bson.NewObjectId().Hex()
    oc.generateVerificationCode()
    expirationDate := time.Now().Add(time.Duration((7*24))*time.Hour)

    for {
        
            
        // Query to insert originator_code
        query := fmt.Sprintf("INSERT INTO `%s`.`originator_code` (`originator_code_id`,`store_id`,`phone`,`expiration_date`,`used_count`,`verification_code`,`verification_code_display`) VALUES (?,?,?,?,?,?,?)",("referralapp_" + oc.StoreId))
        
        parameters := []interface{}{
            oc.OriginatorCodeId,
            oc.StoreId,
            oc.Phone,
            expirationDate,
            0,
            oc.VerificationCode,
            oc.VerificationCodeDisplay,
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
                oc.OriginatorCodeId = bson.NewObjectId().Hex()
                continue
            }
            
            return errors.New(errStr)
        }

        
        break
    }
    
    return nil
}