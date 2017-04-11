package referralapp_api

import (
    // "time"
)

type User struct {
    UserEmail           string         `json:"user_email" dynamo:"user_email"`
    Password            string         `json:"password" dynamo:"password"`
    Store               string         `json:"store_id" dynamo:"store_id"`
    Role                string         `json:"role" dynamo:"role"`
    FirstName           string         `json:"first_name" dynamo:"first_name"`
    MiddleName          string         `json:"middle_name" dynamo:"middle_name"`
    LastName            string         `json:"last_name" dynamo:"last_name"`
    SecondLastName      string         `json:"second_last_name" dynamo:"second_last_name"`
}
