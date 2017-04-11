package referralapp_api

import (
    "github.com/jcgarciaram/general-api/apiutils"
    
    // "time"
    "fmt"
)

type Store struct {
    StoreId         string              `json:"store_id" dynamo:"store_id"`
    StoreName       string              `json:"store_name" dynamo:"store_name"`
    StoreType       string              `json:"store_type" dynamo:"store_type"`
    StoreEmail      string              `json:"store_email" dynamo:"store_email"`
    ContactUserId   string              `json:"contact_user_id" dynamo:"contact_user_id"`
}



// getCreateSchemaQueries returns all the queries for a 
func getCreateSchemaQueries(storeId string) []apiutils.UpsertQuery {

    // Build query to run in MySQL
    createQueries := []apiutils.UpsertQuery{
        {
            Query: fmt.Sprintf("DROP SCHEMA IF EXISTS `%s`", "referralapp_" + storeId),
            Parameters: []interface{}{},
        },
        {
            Query: fmt.Sprintf("CREATE SCHEMA `%s`", "referralapp_" + storeId),
            Parameters: []interface{}{},
        },
        {
            Query: fmt.Sprintf("GRANT SELECT ON `%s`.* TO 'referralappread'", "referralapp_" + storeId),
            Parameters: []interface{}{},
        },
        {
            // referral_code
            Query: fmt.Sprintf(
            "CREATE TABLE `%s`.`referral_code` " +
            "(`referral_code_pk` INT NOT NULL AUTO_INCREMENT COMMENT ''," +
            "`referral_code_id` VARCHAR(45) DEFAULT NULL COMMENT ''," +
            "`store_id` VARCHAR(45) DEFAULT NULL COMMENT ''," +
            "`generated_phone` VARCHAR(45) DEFAULT NULL COMMENT ''," +
            "`expiration_date` TIMESTAMP DEFAULT NULL COMMENT ''," +
            "`used_count` INT NOT NULL COMMENT ''," +
            "`verification_code` INT DEFAULT NULL COMMENT ''," +
            "`verification_code_display` VARCHAR(8) DEFAULT NULL COMMENT ''," +
            "`last_updated_timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''," +
            "`created_timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT ''," +
            "PRIMARY KEY (`referral_code_pk`)  COMMENT '')",
            
            "referralapp_" + storeId),
            
            Parameters: []interface{}{},
        },
    }
    
    return createQueries
}





