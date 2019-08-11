<template>
    <div style="margin-top: 10%">



                    <div class="columns is-vcentered is-centered">
                        <div class="column is-half-desktop ">
                                 <div class="notification is-danger" v-if="errors.length">
                                  <b>Please correct the following error(s):</b>
                                  <ul>
                                    <li v-for="(error,id) in errors" v-bind:key=id>{{ error }}</li>
                                  </ul>
                                </div>
                            <div class="box ">
                                <h4 class="title is-4 has-text-centered">MarchCMS</h4>

                                <hr>

                                <div class="field">
                                    <p class="control has-icons-left has-icons-right">
                                      <input class="input" type="email"
                                       v-model="email"
                                       v-on:keyup.enter="login"
                                       :disabled=is_processing
                                       placeholder="Email">
                                      <span class="icon is-small is-left">
                                        <i class="fa fas fa-envelope"></i>
                                      </span>
                                    </p>
                                </div>

                                <div class="field">
                                    <p class="control has-icons-left">
                                      <input class="input" type="password"
                                      v-model="password"
                                      v-on:keyup.enter="login"
                                      :disabled=is_processing
                                      placeholder="Password">
                                      <span class="icon is-small is-left">
                                        <i class="fa fas fa-lock"></i>
                                      </span>
                                    </p>
                                </div>

                                <div class="field is-pulled-right">
                                    <button class="button is-primary"
                                      @click="login"
                                      :class="[ is_processing ? ' is-loading' : '']"
                                      :disabled=is_processing>
                                        <span class="icon is-small">
                                            <i class="fa fas fa-sign-in"></i>
                                        </span>
                                        &nbspLogin
                                    </button>
                                    <p> <a href="#">Forget Password</a></p>
                                </div>

                                <div style="clear:both"></div>
                            </div>
                        </div>
                    </div>
                </div>
</template>

<script>
export default {

    data(){
      return{
        email:'',
        password:'',
        is_processing:false,
        errors:[],
      }
    },

    methods:{
        login(){
          let vm = this;
          vm.is_processing = true
          if (vm.checkForm()){
          axios.post('/login', {
              Email: vm.email,
              Password: vm.password,
                })
                .then(function (response) {
                  if (response.data.error){
                      vm.errors.push(response.data.error);
                  }else if(response.data.success){
                    alert('redirecting')
                     
                    }
                   vm.is_processing = false
                })
                .catch(function (error) {
                  console.log(error);
                    if (response.data.error){
                      this.errors.push(response.data.error);
                    }
                   vm.is_processing = false
                });
          }else{
             vm.is_processing = false
          }
        },
         checkForm: function() {
            this.errors = [];

            if (!this.email) {
              this.errors.push('Email required.');
            } else if (!this.validEmail(this.email)) {
              this.errors.push('Valid email required.');
            }

            if (!this.password){
                this.errors.push('Password required.');
            }

            if (!this.errors.length) {
              return true;
            }

          },

          validEmail: function (email) {
            var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
          }
    }
}
</script>
