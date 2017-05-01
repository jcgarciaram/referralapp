package referralapp_api

import (
    r "github.com/jcgarciaram/general-api/routes"
)

const (
        UserEmail               = "jcgarciaram@gmail.com"
)

var routes = r.Routes{
    
    // ReferralCode
    r.Route{
        "GetReferralCode",
        "GET",
        "/v1/api/referralcodes/:code",
        GetReferralCode,
    },
    
    r.Route{
        "UseReferralCode",
        "PUT",
        "/v1/api/referralcodes/:code",
        UseReferralCode,
    },

    r.Route{
        "GenerateReferralCode",
        "POST",
        "/v1/api/referralcodes",
        GenerateReferralCode,
    },
    
    r.Route{
        "GetReferralCodes",
        "GET",
        "/v1/api/referralcodes",
        GetReferralCodes,
    },
    
    // OriginatorCode
    r.Route{
        "GetOriginatorCode",
        "GET",
        "/v1/api/stores/:code",
        GetOriginatorCode,
    },
    
    r.Route{
        "UseOriginatorCode",
        "PUT",
        "/v1/api/originatorcodes/:code",
        UseOriginatorCode,
    },
    
    r.Route{
        "GetOriginatorCodes",
        "GET",
        "/v1/api/originatorcodes",
        GetOriginatorCodes,
    },
    
    // User
    r.Route{
        "LogInUser",
        "POST",
        "/v1/api/users/login",
        LogInUser,
    },
    
    r.Route{
        "GetUser",
        "GET",
        "/v1/api/users/:useremail",
        GetUser,
    },
    
    r.Route{
        "PostUser",
        "POST",
        "/v1/api/users",
        PostUser,
    },
    
    // Store
    r.Route{
        "CreateStoreAccount",
        "POST",
        "/v1/api/createaccount",
        CreateStoreAccount,
    },
    
    r.Route{
        "LogInStore",
        "POST",
        "/v1/api/login",
        LogInStore,
    },

    r.Route{
        "GetStore",
        "GET",
        "/v1/api/stores/:store",
        GetStore,
    },

    
}

func GetRoutes() r.Routes {
    return routes
}