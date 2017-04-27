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
        "/v1/api/stores/:store/referralcodes/:code",
        GetReferralCode,
    },
    
    r.Route{
        "UseReferralCode",
        "PUT",
        "/v1/api/stores/:store/referralcodes/:code",
        UseReferralCode,
    },

    r.Route{
        "PostReferralCode",
        "POST",
        "/v1/api/stores/:store/referralcodes",
        PostReferralCode,
    },
    
    r.Route{
        "GetReferralCodes",
        "GET",
        "/v1/api/stores/:store/referralcodes",
        GetReferralCodes,
    },
    
    // OriginatorCode
    r.Route{
        "GetOriginatorCode",
        "GET",
        "/v1/api/stores/:store/originatorcodes/:code",
        GetOriginatorCode,
    },
    
    r.Route{
        "UseOriginatorCode",
        "PUT",
        "/v1/api/stores/:store/originatorcodes/:code",
        UseOriginatorCode,
    },
    
    r.Route{
        "GetOriginatorCodes",
        "GET",
        "/v1/api/stores/:store/originatorcodes",
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
        "PostStore",
        "POST",
        "/v1/api/stores",
        PostStore,
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

