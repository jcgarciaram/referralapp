// API Route which will be prepended :) to all routes called
var api_route = "https://qs1l6iwb53.execute-api.us-east-1.amazonaws.com/prod/v1/api"

// While no real authentication is done, hard-code store_id
var store_id = "58ea6638c2bbb60014dfe7a9"

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
    
    Vue.component('loginpage', {
    template: '\
        <div>\
            <label for="userEmail">User Email</label>\
            <input type="email" class="form-control" id="userEmail" placeholder="user@domain.com"  v-model="loginuser.user_email" v-on:input="updateLoginUser">\
            <label for="userPassword">Password</label>\
            <input type="password" class="form-control" id="userPassword" placeholder="" v-model="loginuser.password" v-on:input="updateLoginUser">\
            <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="login">Log In</button>\
            <div>\
              <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#postUser">\
              create account\
              </button>\
            </div>\
            <div class="modal fade" id="postUser" tabindex="-1" role="dialog" aria-labelledby="postUserLabel" aria-hidden="true">\
              <div class="modal-dialog" role="document">\
                <div class="modal-content">\
                  <div class="modal-header">\
                    <div class="column-left"> \
                    </div>\
                    <div class="column-center">\
                      <h4 class="modal-title" id="postUserLabel">create account</h4>\
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
                        <label for="createUserEmail">User Email</label>\
                        <input type="email" class="form-control" id="createUserEmail" placeholder="user@domain.com"  v-model="createuser.user_email" v-on:input="updateCreateUser">\
                      </div>\
                      <div class="form-group">\
                        <label for="createUserPass">Password</label>\
                        <input type="password" class="form-control" id="createUserPass" placeholder=""  v-model="createuser.password" v-on:input="updateCreateUser">\
                      </div>\
                      <div class="form-group">\
                        <label for="createUserFirstName">First name</label>\
                        <input type="text" class="form-control" id="createUserFirstName" placeholder="Enter first name"  v-model="createuser.first_name" v-on:input="updateCreateUser">\
                      </div>\
                      <div class="form-group">\
                        <label for="createUserMiddleName">Middle name</label>\
                        <input type="text" class="form-control" id="createUserMiddleName" placeholder="Enter middle name"  v-model="createuser.middle_name" v-on:input="updateCreateUser">\
                      </div>\
                      <div class="form-group">\
                        <label for="createUserLastName">Last name</label>\
                        <input type="text" class="form-control" id="createUserLastName" placeholder="Enter last name"  v-model="createuser.last_name" v-on:input="updateCreateUser">\
                      </div>\
                      <div class="form-group">\
                        <label for="createUserSecondLastName">Second last name</label>\
                        <input type="text" class="form-control" id="createUserSecondLastName" placeholder="Enter second last name"  v-model="createuser.second_last_name" v-on:input="updateCreateUser">\
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
        loginuser: {
            type: Object,
            default: function () {
                return { 
                    user_email: '',
                    password: ''
                }
            }
        },
        createuser: {
            type: Object,
            default: function () {
                return { 
                    user_email: '',
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
        updateLoginUser: function () {
          this.$emit('input', this.loginuser)
        },
        updateCreateUser: function () {
          this.$emit('input', this.createuser)
        }
    }
    });
    
    
    // Vue component for the cake-table. This will contain all functions related to cakes
    new Vue({
        el: '#auth-div',
        delimiters: ['$$$', '$$$'],
        data: {
            user: {}
        },
      
        // Methods for all API calls
        methods: {
        
        
            // Create new cake.
            loginUser: function() {
        
                // Not sure what this does :). Copied and pasted from website. Please investigate
                if (!$.trim(this.user)) {
                    this.user = {};
                    return
                }
                
                // Post the new cake to the /cakes route using the $http client
                this.$http.post(api_route + '/users/login', this.user).then((response) => {
        
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
                this.$http.post(api_route + '/users', this.user).then((response) => {
                    
                    if (response.status == 200) {
                        console.log("account has been created");
                    }
                    
                    // Clear out newCake
                    this.user = {};
                    
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

