<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAText</a>
    <div class="navbar-nav">
      <div class="nav-item text-nowrap" v-if="isLoggedIn">
        <div class="d-flex align-items-center">
          <span class="text-white me-3">{{ username }}</span>
          <div class="btn-group">
            <RouterLink 
              to="/profile" 
              class="btn btn-sm btn-outline-light me-2"
            >
              Profile
            </RouterLink>
            <button 
              class="btn btn-sm btn-outline-danger" 
              @click="handleLogout"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
  <div class="container-fluid">
    <div class="row">
      <nav 
        v-if="isLoggedIn" 
        id="sidebarMenu" 
        class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse"
      >
        <div class="position-sticky pt-3 sidebar-sticky">
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/conversations" class="nav-link">
                Conversations
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>
      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView @login-success="updateUserData" />
      </main>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      isLoggedIn: false
    }
  },
  created() {
    this.checkLoginStatus()
  },
  methods: {
    checkLoginStatus() {
      const token = localStorage.getItem('token')
      const userData = JSON.parse(localStorage.getItem('user') || '{}')
      this.isLoggedIn = !!token
      this.username = userData.username || ''
    },
    updateUserData() {
      this.checkLoginStatus()
    },
    handleLogout() {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      this.isLoggedIn = false
      this.username = ''
      this.$router.push('/login')
    }
  },
  watch: {
    '$route'(to) {
      this.checkLoginStatus()
    }
  }
}
</script>