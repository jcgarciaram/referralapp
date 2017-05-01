// API Route which will be prepended :) to all routes called
var api_route = "https://qs1l6iwb53.execute-api.us-east-1.amazonaws.com/prod/v1/api"




// Load when index.html is first open
window.onload = function () {

    

    // While no cookie set up is yet implemented, allow all requests
    Vue.http.headers.common['Authorization'] = 'allow';
    
    // Custom component for datepicker.
    Vue.component('datepicker', {
    template: '\
        <div class="input-group date"><input class="form-control datetimepicker"\
            ref="input"\
            v-bind:value="value"\
            data-date-format="LLL"\
            data-date-end-date="0d"\
            placeholder=""\
            type="text"  />\
            <span class="input-group-addon">\
              <span class="glyphicon glyphicon-calendar"></span>\
            </span></div>\
    ',
    props: {
        value: {
            type: String,
            default: moment().format('LLL')
        }
    },
    mounted: function() {
        let self = this;
        this.$nextTick(function() {
            $(this.$el).datetimepicker().on('dp.change', function(e) {
                //this.value = moment(e.date).format('MMMM Do, YYYY h:mm a')
                self.updateValue(moment(e.date).format('LLL'));
            });
        });
    },
    methods: {
        updateValue: function (value) {
            
            this.$emit('input', value);
        },
    }
    });
    
    
    var GenerateReferralCodeComponent = Vue.component('generatereferralcode', {
    template: '\
        <div>\
            <label for="phoneNumber">Phone Number</label>\
            <input type="tel" class="form-control" id="phoneNumber" placeholder="(555) 555-5555"  v-model="refercode.generated_phone" v-on:input="updateNewReferCode">\
            <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="createreferralcode">Generate personalized QR code</button>\
        </div>\
    ',
    props: {
        refercode: {
            type: Object,
            default: function () {
                return { 
                    generated_phone: ''
                }
            }
        }
    },
    methods: {
        createreferralcode: function () {
            this.$emit('createreferralcode')
        },
        updateNewReferCode: function () {
          this.$emit('input', this.refercode)
        },
    }
    });
    
    
    // Vue component for the cake-table. This will contain all functions related to cakes
    new Vue({
        el: '#referral-code-table',
        delimiters: ['$$$', '$$$'],
        data: {
            referralCodes: [],
            lenReferralCodes: 0,
            newReferralCode: {},
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
            createReferralCode: function() {
        
                // Not sure what this does :). Copied and pasted from website. Please investigate
                if (!$.trim(this.newReferralCode)) {
                    this.newReferralCode = {};
                    return
                }

                phoneString = String(this.newReferralCode.generated_phone)
                
                phoneString = phoneString.replace("(", "")
                phoneString = phoneString.replace(")", "")
                phoneString = phoneString.replace(" ", "")
                phoneString = phoneString.replace("-", "")
                
                if (phoneString.charAt(0) != '+') {
                    this.newReferralCode.generated_phone = '+1' + phoneString
                } else {
                    this.newReferralCode.generated_phone = phoneString
                }
                
                // Post the new cake to the /cakes route using the $http client
                this.$http.post(api_route + '/referralcodes', this.newReferralCode).then((response) => {
                    
                    if (response.status == 403) {
                        window.location.href = 'authenticate.html';
                    }
                    
                    // If API returns with OK status, add the cake to the cakes array
                    if (response.status == 200) {
                        console.log("referralcode has been added!");
                        
                        this.newReferralCode.s3URL = response.body
                    }
                    
                    // Clear out newCake
                    this.newReferralCode = {};
                    
                    // Update lenCakes
                    this.lenReferralCodes = this.referralCodes.length;
        
                // Investigate
                }, (response) => {
                
                    
                    console.log(response.status);
        
                });
            
            
            },
        
            // Function to mark a cake as decorated
            useReferralCode: function(referralCodeId) {
            
                // Make API request
                this.$http.put(api_route + '/referralcodes/' + referralCodeId).then((response) => {
                    
                    if (response.status == 403) {
                        window.location.href = 'authenticate.html';
                    }
                    
                    if (response.status == 200) {
                        console.log("referral code has been used!");
                    }
                    
                    
                }, (response) => {
                
                    console.log(response)
                    
                })
            
            },
            
            isLoggedIn: function()  {
                
                if (window.localStorage.getItem('loggedIn')) {
                    return true;
                } else {
                    return false;
                }
                    
            }

        }
    });

}

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


