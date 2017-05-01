// API Route which will be prepended :) to all routes called
var api_route = "https://qs1l6iwb53.execute-api.us-east-1.amazonaws.com/prod/v1/api"

// While no real authentication is done, hard-code store_id
// var store_id = "59021b5369c9e30013d404ca"

function addLoadEvent(func) {
    var oldonload = window.onload;
    if (typeof window.onload != 'function') {
        window.onload = func;
    } else {
        window.onload = function() {
            if (oldonload) {
                oldonload();
            }
            func();
        }
    }
}

addLoadEvent(function() {

    // While no cookie set up is yet implemented, allow all requests
    Vue.http.headers.common['Authorization'] = 'allow';
    
    var LoginComponent = Vue.component('loginpage', {
    template: '\
        <div>\
            <label for="storeEmail">Store Email</label>\
            <input type="email" class="form-control" id="storeEmail" placeholder="store@domain.com"  v-model="loginstore.store_email" v-on:input="updateLoginStore">\
            <label for="storePassword">Password</label>\
            <input type="password" class="form-control" id="storePassword" placeholder="" v-model="loginstore.password" v-on:input="updateLoginStore">\
            <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="login">Log In</button>\
            <div>\
              <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#postStore">\
              create account\
              </button>\
            </div>\
            <div class="modal fade" id="postStore" tabindex="-1" role="dialog" aria-labelledby="postStoreLabel" aria-hidden="true">\
              <div class="modal-dialog" role="document">\
                <div class="modal-content">\
                  <div class="modal-header">\
                    <div class="column-left"> \
                    </div>\
                    <div class="column-center">\
                      <h4 class="modal-title" id="postStoreLabel">create account</h4>\
                    </div>\
                    <div id="modal-title-col" class="column-right">\
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">\
                      <span aria-hidden="true">&times;</span>\
                      </button>\
                    </div>\
                  </div>\
                  <div class="modal-body">\
                    <form name="createAccountForm">\
                      <div class="form-group">\
                        <label for="createStoreEmail">Store Email</label>\
                        <input type="email" class="form-control" id="createStoreEmail" placeholder="store@domain.com"  v-model="createstore.store_email" v-on:input="updateCreateStore">\
                      </div>\
                      <div class="form-group">\
                        <label for="createStorePass">Password</label>\
                        <input type="password" class="form-control" id="createStorePass" placeholder=""  v-model="createstore.password" v-on:input="updateCreateStore">\
                      </div>\
                      <div class="form-group">\
                        <label for="createStoreFirstName">First name</label>\
                        <input type="text" class="form-control" id="createStoreFirstName" placeholder="Enter first name"  v-model="createstore.first_name" v-on:input="updateCreateStore">\
                      </div>\
                      <div class="form-group">\
                        <label for="createStoreMiddleName">Middle name</label>\
                        <input type="text" class="form-control" id="createStoreMiddleName" placeholder="Enter middle name"  v-model="createstore.middle_name" v-on:input="updateCreateStore">\
                      </div>\
                      <div class="form-group">\
                        <label for="createStoreLastName">Last name</label>\
                        <input type="text" class="form-control" id="createStoreLastName" placeholder="Enter last name"  v-model="createstore.last_name" v-on:input="updateCreateStore">\
                      </div>\
                      <div class="form-group">\
                        <label for="createStoreSecondLastName">Second last name</label>\
                        <input type="text" class="form-control" id="createStoreSecondLastName" placeholder="Enter second last name"  v-model="createstore.second_last_name" v-on:input="updateCreateStore">\
                      </div>\
                    </form>\
                  </div>\
                  <div class="modal-footer">\
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button>\
                    <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="createaccount">Create account</button>\
                  </div>\
                </div>\
              </div>\
            </div>\
        </div>\
    ',
    props: {
        loginstore: {
            type: Object,
            default: function () {
                return { 
                    store_email: '',
                    password: ''
                }
            }
        },
        createstore: {
            type: Object,
            default: function () {
                return { 
                    store_email: '',
                    password: '',
                    first_name: '',
                    middle_name: '',
                    last_name: '',
                    second_last_name: '',
                }
            }
        }
    },
    methods: {
        createaccount: function () {
            this.$emit('createaccount')
        },
        login: function () {
            this.$emit('login')
        },
        updateLoginStore: function () {
          this.$emit('input', this.loginstore)
        },
        updateCreateStore: function () {
          this.$emit('input', this.createstore)
        }
    }
    });
    
    
    // Vue component for the cake-table. This will contain all functions related to cakes
    new Vue({
        el: '#auth-div',
        delimiters: ['$$$', '$$$'],
        data: {
            store: {},
            loggedIn: false
        },

        // This is run whenever the page is loaded to make sure we have a current cake list
        created: function() {
      
            if (window.localStorage.getItem('loggedIn')) {
                    this.loggedIn = true;
                    
            } else {
                    this.loggedIn = false;
            }
        },
      
        // Methods for all API calls
        methods: {
        
        
            // Create new cake.
            loginStore: function() {
        
                // Not sure what this does :). Copied and pasted from website. Please investigate
                if (!$.trim(this.store)) {
                    this.store = {};
                    return
                }
                
                // Post the new cake to the /cakes route using the $http client
                this.$http.post(api_route + '/login', this.store).then((response) => {
        
                    // If API returns with OK status, add the cake to the cakes array
                    if (response.status == 200) {
                        console.log("logged in!");
                        console.log(response.body);
                        
                        window.localStorage.setItem('loggedIn', true);
                        // window.location.href = 'index.html';
                        
                    }

        
                // Investigate
                }, (response) => {
                
                    
                    console.log(response.status, response.body);
        
                });
            
            
            },
        
            // Function to mark a cake as decorated
            createAccount: function() {
            
                // Make API request
                this.$http.post(api_route + '/createaccount', this.store).then((response) => {
                    
                    if (response.status == 200) {
                        console.log("account has been created");
                    }
                    
                    // Clear out newCake
                    this.store = {};
                    
                }, (response) => {
                
                    console.log(response)
                    
                })
                
                // Reset the createAccount to be blank
                var frm = document.getElementsByName('createAccountForm')[0];
                frm.reset();
            
            },

        }
    });

})

// insertCake inserts a new cake into the cakes array into a correct sorted position based on the pickup_timestamp
function insertCake(element, array) {
    array.splice(locationOfCake(element, array), 0, element);
    return array;
}

// locationOfCake figures out the correct location of where to insert the cake based on the pickup_timestamp
function locationOfCake(element, array) {
    
    // If array is empty
    if (array.length == 0) {
        return 0
    } 
    
    // Define pivot as middle point
    var pivot = parseInt(array.length / 2);
    
    var elPickupTimestamp = moment(element.pickup_timestamp, "LLL")
    var pivotPickupTimestamp = moment(array[pivot].pickup_timestamp, "LLL")
    
    // If element is the same as the middle, insert right after
    if (pivotPickupTimestamp.isSame(elPickupTimestamp)) {
        return pivot+1;
    }
    
    // If element is less than middle element
    if (elPickupTimestamp.isBefore(pivotPickupTimestamp)) {
    
        return locationOfCake(element, array.slice(0, pivot));
    
    // If element is greater than middle element
    } else {
        return (pivot+1) + locationOfCake(element, array.slice(pivot+1, array.length));
    }
}

