<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAText</a>
    <div class="navbar-nav">
      <div class="nav-item text-nowrap" v-if="isLoggedIn">
        <div class="d-flex align-items-center">
          <span class="text-white me-3">{{ currentUsername }}</span>
          <div class="btn-group">
            <RouterLink 
              to="/profile" 
              class="btn btn-sm btn-outline-light me-2"
            >
              Profile
            </RouterLink>
            <button 
              class="btn btn-sm btn-outline-danger" 
              @click="logout"
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'

export default {
  setup() {
    const router = useRouter()
    const currentUsername = ref('')
    
    const updateUsername = () => {
      const user = JSON.parse(localStorage.getItem('user'))
      currentUsername.value = user?.username || ''
    }

    const isLoggedIn = computed(() => {
      return localStorage.getItem('token') !== null
    })

    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      currentUsername.value = ''
      router.push('/login')
    }

    // Listen for storage changes
    const storageHandler = (event) => {
      if (event.key === 'user') {
        updateUsername()
      }
    }

    // Watch for route changes and update username
    watch(
      () => router.currentRoute.value,
      () => {
        updateUsername()
      },
      { immediate: true }
    )

    onMounted(() => {
      updateUsername()
      window.addEventListener('storage', storageHandler)
    })

    onUnmounted(() => {
      window.removeEventListener('storage', storageHandler)
    })

    return {
      currentUsername,
      isLoggedIn,
      logout
    }
  }
}
</script>
