<template>
  <div class="modal" tabindex="-1" style="display: block;">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Create New Group</h5>
          <button 
            type="button" 
            class="btn-close" 
            @click="$emit('close')"
          ></button>
        </div>
        <div class="modal-body">
          <div class="mb-3">
            <label class="form-label">Group Name</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="groupName" 
              placeholder="Enter group name"
              required
            >
          </div>

          <div class="mb-3">
            <label class="form-label">Add Members</label>
            <div class="input-group mb-3">
              <input 
                type="text" 
                class="form-control" 
                placeholder="Search users" 
                v-model="searchQuery"
                @input="searchUsers"
              >
            </div>

            <div class="selected-members mb-3">
              <span 
                v-for="member in selectedMembers" 
                :key="member.id" 
                class="badge bg-primary me-2 mb-2"
              >
                {{ member.username }}
                <button 
                  class="btn-close btn-close-white" 
                  @click="removeMember(member)"
                ></button>
              </span>
            </div>

            <LoadingSpinner :loading="loading">
              <ErrorMsg v-if="errorMsg" :msg="errorMsg" />

              <div class="list-group">
                <button 
                  v-for="user in searchResults" 
                  :key="user.id"
                  class="list-group-item list-group-item-action"
                  @click="addMember(user)"
                  :disabled="isMemberSelected(user)"
                >
                  {{ user.username }}
                </button>
              </div>
            </LoadingSpinner>
          </div>
        </div>
        <div class="modal-footer">
          <button 
            type="button" 
            class="btn btn-secondary" 
            @click="$emit('close')"
          >
            Cancel
          </button>
          <button 
            type="button" 
            class="btn btn-primary" 
            @click="createGroup"
            :disabled="!isFormValid"
          >
            Create Group
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      groupName: '',
      searchQuery: '',
      searchResults: [],
      selectedMembers: [],
      loading: false,
      errorMsg: null,
      searchTimeout: null
    }
  },
  computed: {
    isFormValid() {
      return this.groupName.trim().length > 0 && this.selectedMembers.length > 0
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
          
          this.searchResults = response.data.filter(
            user => !this.selectedMembers.some(m => m.id === user.id)
          )
        } catch (error) {
          this.errorMsg = 'Failed to search users'
          console.error(error)
        } finally {
          this.loading = false
        }
      }, 300)
    },
    addMember(user) {
      if (!this.isMemberSelected(user)) {
        this.selectedMembers.push(user)
        this.searchQuery = ''
        this.searchResults = []
      }
    },
    removeMember(user) {
      this.selectedMembers = this.selectedMembers.filter(m => m.id !== user.id)
    },
    isMemberSelected(user) {
      return this.selectedMembers.some(m => m.id === user.id)
    },
    async createGroup() {
      if (!this.isFormValid) {
        this.errorMsg = 'Please enter a group name and select at least one member'
        return
      }

      try {
        const response = await this.$axios.post('/group', {
          name: this.groupName,
          usernames: this.selectedMembers.map(m => m.username)
        }, {
          headers: { 'Authorization': localStorage.getItem('token') }
        })
        
        this.$emit('group-created', response.data)
      } catch (error) {
        this.errorMsg = 'Failed to create group'
        console.error(error)
      }
    }
  }
}
</script>

<style scoped>
.modal {
  background-color: rgba(0,0,0,0.5);
}
.selected-members .badge {
  display: inline-flex;
  align-items: center;
}
.selected-members .btn-close {
  margin-left: 5px;
  font-size: 0.5rem;
}
</style>