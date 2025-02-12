<template>
  <div class="conversation-detail h-100 d-flex flex-column">
    <!-- Conversation Header -->
    <div class="conversation-header d-flex justify-content-between align-items-center p-3 border-bottom">
      <div class="d-flex align-items-center">
        <img 
          :src="conversation.photo || '/default-avatar.png'" 
          class="rounded-circle me-3" 
          width="50" 
          height="50" 
          alt="Conversation Avatar"
        >
        <h5 class="mb-0">{{ conversation.name }}</h5>
      </div>
      <div class="conversation-actions">
        <div class="btn-group" v-if="conversation.isGroup">
          <button 
            class="btn btn-sm btn-outline-secondary" 
            @click="showAddMemberModal = true"
          >
            Add Member
          </button>
          <button 
            class="btn btn-sm btn-outline-danger" 
            @click="leaveGroup"
          >
            Leave Group
          </button>
        </div>
      </div>
    </div>

    <!-- Messages Container -->
    <div 
      ref="messagesContainer" 
      class="messages-container flex-grow-1 p-3 overflow-auto"
    >
      <LoadingSpinner :loading="loading">
        <ErrorMsg v-if="errorMsg" :msg="errorMsg" />
        
        <div 
          v-for="message in messages" 
          :key="message.messageId" 
          class="message mb-3"
          :class="{
            'message-sent': message.senderId === currentUserId,
            'message-received': message.senderId !== currentUserId
          }"
        >
          <div class="message-content">
            <div class="message-header d-flex justify-content-between">
              <small class="text-muted">
                {{ message.senderUsername }}
              </small>
              <small class="text-muted">
                {{ formatDate(message.sendTime) }}
              </small>
            </div>
            
            <!-- Message Body -->
            <div 
              v-if="message.text" 
              class="message-text p-2 rounded"
            >
              {{ message.text }}
            </div>
            
            <!-- Photo Message -->
            <img 
              v-if="message.photo" 
              :src="message.photo" 
              class="img-fluid rounded mt-2" 
              alt="Message Photo"
            >
            
            <!-- Message Actions -->
            <div class="message-actions mt-2 d-flex justify-content-between">
              <!-- Forward and Delete only for sent messages -->
              <div v-if="message.senderId === currentUserId">
                <button 
                  class="btn btn-sm btn-outline-secondary me-2"
                  @click="forwardMessage(message)"
                >
                  Forward
                </button>
                <button 
                  class="btn btn-sm btn-outline-danger"
                  @click="deleteMessage(message)"
                >
                  Delete
                </button>
              </div>
              
              <!-- Reactions -->
              <div class="reactions">
                <div 
                  v-for="reaction in message.comments" 
                  :key="reaction.userId" 
                  class="reaction badge bg-light text-dark me-1"
                >
                  {{ reaction.emoji }}
                  <span 
                    v-if="reaction.userId === currentUserId" 
                    @click="removeReaction(message, reaction)"
                    class="remove-reaction"
                  >
                    ✖
                  </span>
                </div>
                
                <!-- Add Reaction Button -->
                <button 
                  v-if="message.senderId !== currentUserId"
                  class="btn btn-sm btn-outline-secondary"
                  @click="showReactionModal(message)"
                >
                  React
                </button>
              </div>
            </div>
          </div>
        </div>
      </LoadingSpinner>
    </div>

    <!-- Message Input -->
    <div class="message-input p-3 border-top">
      <div class="input-group">
        <textarea 
          class="form-control" 
          placeholder="Type a message..." 
          v-model="newMessage"
          @keydown.enter.prevent="sendMessage"
          rows="3"
        ></textarea>
        <button 
          class="btn btn-primary" 
          @click="sendMessage"
        >
          Send
        </button>
      </div>
      <div class="d-flex mt-2">
        <input 
          type="file" 
          ref="photoInput" 
          @change="handlePhotoUpload" 
          accept="image/*" 
          class="d-none"
        >
        <button 
          class="btn btn-outline-secondary" 
          @click="$refs.photoInput.click()"
        >
          Add Photo
        </button>
      </div>
    </div>

    <!-- Reaction Picker Modal -->
    <ReactionPickerModal 
      v-if="currentReactionMessage" 
      @close="currentReactionMessage = null"
      @select-reaction="addReaction"
    />

    <!-- User Search Modal for Adding to Group -->
    <UserSearchModal 
      v-if="showAddMemberModal" 
      @close="showAddMemberModal = false"
      @start-conversation="addMemberToGroup"
    />
  </div>
</template>

<script>
import ReactionPickerModal from './ReactionPickerModal.vue'
import UserSearchModal from './UserSearchModal.vue'

export default {
  components: {
    ReactionPickerModal,
    UserSearchModal
  },
  props: {
    conversation: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      messages: [],
      newMessage: '',
      loading: true,
      errorMsg: null,
      currentReactionMessage: null,
      showAddMemberModal: false,
      currentUserId: parseInt(localStorage.getItem('token'))
    }
  },
  mounted() {
    this.fetchConversationDetails()
  },
  methods: {
    async fetchConversationDetails() {
      this.loading = true
      this.errorMsg = null

      try {
        const response = await this.$axios.get(`/conversation/${this.conversation.conversationId}`)
        this.messages = response.data.messages
        this.loading = false
      } catch (error) {
        this.errorMsg = 'Failed to load conversation'
        this.loading = false
        console.error(error)
      }
    },
    async sendMessage() {
      if (!this.newMessage.trim() && !this.$refs.photoInput?.files[0]) {
        return
      }

      try {
        const messageData = {
          conversationId: this.conversation.conversationId,
          text: this.newMessage
        }

        const response = await this.$axios.post('/message', messageData)
        
        // Clear input after sending
        this.newMessage = ''
        
        // Refresh conversation
        await this.fetchConversationDetails()
      } catch (error) {
        this.errorMsg = 'Failed to send message'
        console.error(error)
      }
    },
    async handlePhotoUpload(event) {
      const file = event.target.files[0]
      if (!file) return

      const reader = new FileReader()
      reader.onload = async (e) => {
        try {
          const base64Photo = e.target.result.split(',')[1]
          
          const messageData = {
            conversationId: this.conversation.conversationId,
            photo: base64Photo
          }

          await this.$axios.post('/message', messageData)
          
          // Refresh conversation
          await this.fetchConversationDetails()
        } catch (error) {
          this.errorMsg = 'Failed to upload photo'
          console.error(error)
        }
      }
      reader.readAsDataURL(file)
    },
    async forwardMessage(message) {
      try {
        await this.$axios.post(`/message/${message.messageId}/forward`, {
          conversationId: this.conversation.conversationId
        })
        
        // Refresh conversation
        await this.fetchConversationDetails()
      } catch (error) {
        this.errorMsg = 'Failed to forward message'
        console.error(error)
      }
    },
    async deleteMessage(message) {
      try {
        await this.$axios.delete(`/message/${message.messageId}`)
        
        // Refresh conversation
        await this.fetchConversationDetails()
      } catch (error) {
        this.errorMsg = 'Failed to delete message'
        console.error(error)
      }
    },
    showReactionModal(message) {
      this.currentReactionMessage = message
    },
    async addReaction(emoji) {
      if (!this.currentReactionMessage) return

      try {
        await this.$axios.post(`/message/${this.currentReactionMessage.messageId}/comment`, {
          emoji: emoji
        })
        
        // Refresh conversation
        await this.fetchConversationDetails()
        
        this.currentReactionMessage = null
      } catch (error) {
        this.errorMsg = 'Failed to add reaction'
        console.error(error)
      }
    },
    async removeReaction(message, reaction) {
      try {
        await this.$axios.delete(`/message/${message.messageId}/uncomment`)
        
        // Refresh conversation
        await this.fetchConversationDetails()
      } catch (error) {
        this.errorMsg = 'Failed to remove reaction'
        console.error(error)
      }
    },
    async leaveGroup() {
      try {
        await this.$axios.delete(`/group/${this.conversation.conversationId}/leave`)
        
        // Navigate back to conversations
        this.$router.push('/')
      } catch (error) {
        this.errorMsg = 'Failed to leave group'
        console.error(error)
      }
    },
    async addMemberToGroup(user) {
      try {
        await this.$axios.post(`/group/${this.conversation.conversationId}/add`, {
          username: user.username
        })
        
        // Close modal
        this.showAddMemberModal = false
        
        // Refresh conversation
        await this.fetchConversationDetails()
      } catch (error) {
        this.errorMsg = 'Failed to add member to group'
        console.error(error)
      }
    },
    formatDate(dateString) {
      return new Date(dateString).toLocaleString()
    }
  }
}
</script>

<style scoped>
.message-sent {
  text-align: right;
}

.message-received {
  text-align: left;
}

.message-content {
  display: inline-block;
  max-width: 80%;
  padding: 10px;
  border-radius: 10px;
  background-color: #f1f0f0;
}

.message-sent .message-content {
  background-color: #dcf8c6;
}

.reactions .remove-reaction {
  cursor: pointer;
  margin-left: 5px;
  color: red;
}
</style>