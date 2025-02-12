<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAText</a>
    <div class="navbar-nav">
      <div class="nav-item text-nowrap" v-if="isLoggedIn">
        <div class="dropdown">
          <button 
            class="btn btn-dark dropdown-toggle" 
            type="button" 
            id="userDropdown" 
            data-bs-toggle="dropdown" 
            aria-expanded="false"
          >
            {{ currentUser.username }}
          </button>
          <ul class="dropdown-menu" aria-labelledby="userDropdown">
            <li>
              <RouterLink to="/profile" class="dropdown-item">
                <i class="feather feather-user me-2"></i>Profile
              </RouterLink>
            </li>
            <li>
              <a class="dropdown-item" href="#" @click.prevent="logout">
                <i class="feather feather-log-out me-2"></i>Logout
              </a>
            </li>
          </ul>
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
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>Conversations</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/" class="nav-link">
                <i class="feather feather-message-circle me-2"></i>
                All Conversations
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>
      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'

const router = useRouter()
const currentUser = ref(JSON.parse(localStorage.getItem('user')) || {})

const isLoggedIn = computed(() => {
  return localStorage.getItem('token') !== null
})

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style>
.sidebar .nav-link {
  display: flex;
  align-items: center;
}
.sidebar .nav-link i {
  margin-right: 8px;
}
</style>