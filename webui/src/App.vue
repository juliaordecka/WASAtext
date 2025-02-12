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
            {{ username }}
          </button>
          <ul class="dropdown-menu" aria-labelledby="userDropdown">
            <li>
              <RouterLink to="/profile" class="dropdown-item">
                Profile
              </RouterLink>
            </li>
            <li>
              <a class="dropdown-item" href="#" @click.prevent="logout">
                Logout
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
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/" class="nav-link">
                Conversations
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

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'

export default {
  setup() {
    const router = useRouter()
    const username = ref('')
    
    const updateUsername = () => {
      const user = JSON.parse(localStorage.getItem('user'))
      username.value = user?.username || ''
    }

    const isLoggedIn = computed(() => {
      return localStorage.getItem('token') !== null
    })

    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      username.value = ''
      router.push('/login')
    }

    // Listen for storage changes
    const storageHandler = (event) => {
      if (event.key === 'user') {
        updateUsername()
      }
    }

    onMounted(() => {
      updateUsername()
      window.addEventListener('storage', storageHandler)
    })

    onUnmounted(() => {
      window.removeEventListener('storage', storageHandler)
    })

    return {
      username,
      isLoggedIn,
      logout
    }
  }
}
</script>