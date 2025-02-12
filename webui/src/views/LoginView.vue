<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">Login to WASAText</div>
          <div class="card-body">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg" />
            <form @submit.prevent="login">
              <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="username" 
                  v-model="username" 
                  required 
                  maxlength="30"
                  pattern="^[a-zA-Z0-9_]+$"
                  title="Username can only contain letters, numbers, and underscores"
                >
                <small class="form-text text-muted">
                  Username must be 1-30 characters long, containing only letters, numbers, and underscores
                </small>
              </div>
              <button 
                type="submit" 
                class="btn btn-primary" 
                :disabled="!isValidUsername"
              >
                Login
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      errorMsg: null
    }
  },
  computed: {
    isValidUsername() {
      // Username validation: 1-30 characters, only letters, numbers, and underscores
      return /^[a-zA-Z0-9_]{1,30}$/.test(this.username)
    }
  },
  methods: {
async login() {
      if (!this.isValidUsername) {
        this.errorMsg = 'Invalid username format'
        return
      }

      try {
        // First do the login
        const response = await this.$axios.post('/session', { username: this.username })
        
        // Then get complete user data including photo
        const userData = {
          id: response.data.id,
          username: this.username,
          profilePhoto: localStorage.getItem('user') ? 
            JSON.parse(localStorage.getItem('user')).profilePhoto : null
        }
        
        // Store complete user data
        localStorage.setItem('user', JSON.stringify(userData))
        localStorage.setItem('token', response.data.id)
        
        // Navigate to conversations page
        this.$router.push('/')
      } catch (error) {
        this.errorMsg = error.response?.data || 'Login failed'
      }
}


  }
}
</script>

<style scoped>
.card {
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}
</style>