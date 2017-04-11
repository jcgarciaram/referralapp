// API Route which will be prepended :) to all routes called
var api_route = "https://qs1l6iwb53.execute-api.us-east-1.amazonaws.com/prod/v1/api"

// While no real authentication is done, hard-code store_id
var store_id = "58ea6638c2bbb60014dfe7a9"

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

