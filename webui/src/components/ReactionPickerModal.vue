<template>
  <div class="modal" tabindex="-1" style="display: block;">
    <div class="modal-dialog modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">React to Message</h5>
          <button 
            type="button" 
            class="btn-close" 
            @click="$emit('close')"
          ></button>
        </div>
        <div class="modal-body">
          <div class="emoji-picker d-flex flex-wrap justify-content-center">
            <button 
              v-for="emoji in commonEmojis" 
              :key="emoji"
              class="btn btn-outline-secondary m-1"
              @click="selectReaction(emoji)"
            >
              {{ emoji }}
            </button>
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
      commonEmojis: [
        '👍', '❤️', '😂', '😮', '😢', '🔥', 
        '👏', '😍', '🤔', '🙌', '🎉', '👀'
      ]
    }
  },
  methods: {
    selectReaction(emoji) {
      this.$emit('select-reaction', emoji)
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.modal {
  background-color: rgba(0,0,0,0.5);
}
.emoji-picker .btn {
  font-size: 1.5rem;
  margin: 5px;
  width: 50px;
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>