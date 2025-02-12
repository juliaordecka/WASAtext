<template>
  <div class="modal" tabindex="-1" style="display: block;">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Search Users</h5>
          <button 
            type="button" 
            class="btn-close" 
            @click="$emit('close')"
          ></button>
        </div>
        <div class="modal-body">
          <div class="input-group mb-3">
            <input 
              type="text" 
              class="form-control" 
              placeholder="Search by username" 
              v-model="searchQuery"
              @input="searchUsers"
            >
          </div>

          <LoadingSpinner :loading="loading">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg" />

            <div 
              v-if="searchResults.length === 0 && searchQuery" 
              class="text-center text-muted"
            >
              No users found
            </div>

            <div class="list-group">
              <button 
                v-for="user in searchResults" 
                :key="user.id"
                class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
                @click="startConversation(user)"
              >
                <div class="d-flex align-items-center">
                  <img 
                    :src="user.profilePhoto || '/default-avatar.png'" 
                    class="rounded-circle me-3" 
                    width="40" 
                    height="40" 
                    alt="User Avatar"
                  >
                  <span>{{ user.username }}</span>
                </div>
                <span class="badge bg-primary rounded-pill">Start Chat</span>
              </button>
            </div>
          </LoadingSpinner>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      searchQuery: '',
      searchResults: [],
      loading: false,
      errorMsg: null,
      searchTimeout: null
    }
  },
  methods: {
    searchUsers() {
      if (this.searchTimeout) {
        clearTimeout(this.searchTimeout)
      }

      this.searchTimeout = setTimeout(async () => {
        if (!this.searchQuery) {
          this.searchResults = []
          return
        }

        this.loading = true
        this.errorMsg = null

        try {
          const response = await this.$axios.get('/users/search', {
            params: { username: this.searchQuery },
            headers: { 'Authorization': localStorage.getItem('token') }
          })
          
          this.searchResults = response.data
        } catch (error) {
          this.errorMsg = 'Failed to search users'
          console.error(error)
        } finally {
          this.loading = false
        }
      }, 300)
    },
    startConversation(user) {
      this.$emit('start-conversation', user)
    }
  }
}
</script>