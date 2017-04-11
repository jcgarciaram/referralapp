package main

import (

    api "github.com/jcgarciaram/general-api"
    "github.com/jcgarciaram/general-api/routes"
    ra "github.com/jcgarciaram/referralapp/referralapp_api"
    
)


func main() {
	r := routes.Routes{}
    
    // Append general_museum routes
    r.AppendRoutes(ra.GetRoutes())
    

    router := api.NewRouter(r)


    router.Listen()
    // router.Gateway()
    
}
