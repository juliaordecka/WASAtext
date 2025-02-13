<template>
	<div class="conversation-detail h-100 d-flex flex-column">
		<!-- Conversation Header -->
		<div
			class="conversation-header d-flex justify-content-between align-items-center p-3 border-bottom"
		>
			<div class="d-flex align-items-center">
				<img
					:src="conversation.photo || '/default-pic.jpg'"
					class="rounded-circle me-3"
					width="50"
					height="50"
					alt="Conversation Avatar"
				/>
				<h5 class="mb-0">{{ conversationName }}</h5>
			</div>
			<div class="conversation-actions">
				<div class="btn-group" v-if="conversation.isGroup">
					<button
						class="btn btn-sm btn-outline-secondary"
						@click="showGroupPhotoInput = true"
					>
						Set Photo
					</button>
					<button
						class="btn btn-sm btn-outline-secondary"
						@click="showGroupNameInput = true"
					>
						Set Name
					</button>

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

		<div
			v-if="showGroupPhotoInput && conversation.isGroup"
			class="p-3 border-bottom"
		>
			<div class="input-group">
				<input
					type="file"
					class="form-control"
					@change="handleGroupPhotoUpload"
					accept="image/*"
				/>
				<button
					class="btn btn-outline-secondary"
					@click="showGroupPhotoInput = false"
				>
					Cancel
				</button>
			</div>
		</div>

		<div
			v-if="showGroupNameInput && conversation.isGroup"
			class="p-3 border-bottom"
		>
			<div class="input-group">
				<input
					type="text"
					class="form-control"
					v-model="newGroupName"
					placeholder="Enter new group name"
				/>
				<button class="btn btn-primary" @click="updateGroupName">
					Save
				</button>
				<button
					class="btn btn-outline-secondary"
					@click="showGroupNameInput = false"
				>
					Cancel
				</button>
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
					v-for="message in sortedMessages"
					:key="message.messageId"
					class="message mb-3"
					:class="{
						'message-sent': message.senderId === currentUserId,
						'message-received': message.senderId !== currentUserId,
					}"
				>
					<div class="message-content">
						<div
							class="message-header d-flex justify-content-between"
						>
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
							style="white-space: pre-wrap"
						>
							{{ message.text }}
						</div>

						<!-- Photo Message -->
						<img
							v-if="message.photo"
							:src="'data:image/jpeg;base64,' + message.photo"
							class="img-fluid rounded mt-2"
							alt="Message Photo"
						/>

						<!-- Message Status (for sent messages) -->
						<div
							v-if="message.senderId === currentUserId"
							class="message-status"
						>
							<span v-if="message.status === 'Sent'">✓</span>
							<span v-else-if="message.status === 'Read'"
								>✓✓</span
							>
						</div>

						<!-- Message Actions -->
						<div
							class="message-actions mt-2 d-flex justify-content-between"
						>
							<div>
								<button
									class="btn btn-sm btn-outline-secondary me-2"
									@click="forwardMessage(message)"
								>
									Forward
								</button>
								<button
									v-if="message.senderId === currentUserId"
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
										@click="
											removeReaction(message, reaction)
										"
										class="remove-reaction"
									>
										✖
									</span>
								</div>

								<!-- Add Reaction Button -->
								<button
									v-if="
										message.senderId !== currentUserId &&
										!hasUserReacted(message)
									"
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
					@keydown.enter.exact.prevent="sendMessage"
					@keydown.enter.shift.exact.prevent="newMessage += '\n'"
					rows="3"
				></textarea>
				<button class="btn btn-primary" @click="sendMessage">
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
				/>
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

		<!-- Add this new modal for forwarding -->
		<div
			v-if="showForwardModal"
			class="modal"
			tabindex="-1"
			style="display: block"
		>
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">Forward Message</h5>
						<button
							type="button"
							class="btn-close"
							@click="showForwardModal = false"
						></button>
					</div>
					<div class="modal-body">
						<h6>Select a conversation to forward to:</h6>
						<div class="list-group">
							<button
								v-for="destination in forwardDestinations"
								:key="destination.conversationId"
								class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
								@click="handleForward(destination)"
							>
								<div>
									<img
										:src="
											destination.photo ||
											'/default-avatar.png'
										"
										class="rounded-circle me-2"
										width="30"
										height="30"
										alt="Conversation Avatar"
									/>
									{{ destination.name }}
								</div>
								<span class="badge bg-primary rounded-pill">
									{{ destination.isGroup ? "Group" : "Chat" }}
								</span>
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import ReactionPickerModal from "./ReactionPickerModal.vue";
import UserSearchModal from "./UserSearchModal.vue";

export default {
	components: {
		ReactionPickerModal,
		UserSearchModal,
	},
	props: {
		conversation: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			messages: [],
			newMessage: "",
			loading: true,
			errorMsg: null,
			currentReactionMessage: null,
			showAddMemberModal: false,
			currentUserId: parseInt(localStorage.getItem("token")),
			refreshInterval: null,
			showGroupPhotoInput: false,
			showGroupNameInput: false,
			newGroupName: "",
			showForwardModal: false,
			messageToForward: null,
			forwardDestinations: [],
		};
	},
	computed: {
		conversationName() {
			// Use the conversation prop, which should update when the conversation changes
			return this.conversation.name;
		},

		sortedMessages() {
			return [...this.messages].sort(
				(a, b) => new Date(a.sendTime) - new Date(b.sendTime)
			);
		},
	},
	mounted() {
		this.fetchConversationDetails();
		this.refreshInterval = setInterval(this.fetchConversationDetails, 5000);
	},
	beforeUnmount() {
		if (this.refreshInterval) {
			clearInterval(this.refreshInterval);
		}
	},
	methods: {
		async fetchConversationDetails() {
			// Reset messages before fetching
			this.messages = [];
			this.loading = true;
			this.errorMsg = null;

			try {
				const response = await this.$axios.get(
					`/conversation/${this.conversation.conversationId}`
				);

				// Only update if we're still on the same conversation
				if (
					this.conversation.conversationId ===
					response.data.conversationId
				) {
					// Explicitly set messages, even if it's an empty array
					this.messages = response.data.messages || [];
					this.loading = false;
				}
			} catch (error) {
				console.error("Fetch conversation error:", error);
				this.errorMsg = "Failed to load conversation";
				this.messages = []; // Ensure messages are cleared on error
				this.loading = false;
			}
		},
		async sendMessage() {
			if (!this.newMessage.trim()) {
				return;
			}

			try {
				await this.$axios.post("/message", {
					conversationId: this.conversation.conversationId,
					text: this.newMessage,
				});

				this.newMessage = "";
				await this.fetchConversationDetails();
			} catch (error) {
				console.error("Send message error:", error);
				this.errorMsg = "Failed to send message";
			}
		},
		async handlePhotoUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			// Add size check
			if (file.size > 5 * 1024 * 1024) {
				// 5MB limit
				this.errorMsg =
					"File too large. Please choose a file under 5MB.";
				return;
			}

			const reader = new FileReader();
			reader.onload = async (e) => {
				try {
					const base64Photo = e.target.result.split(",")[1];
					console.log("Sending photo with size:", base64Photo.length);

					const response = await this.$axios.post("/message", {
						conversationId: this.conversation.conversationId,
						photo: base64Photo,
						text: "Photo message", // Changed from empty text to actual text
					});

					console.log("Photo upload response:", response);
					await this.fetchConversationDetails();
					event.target.value = "";
				} catch (error) {
					console.error(
						"Upload photo error details:",
						error.response || error
					);
					this.errorMsg = `Failed to upload photo: ${
						error.response && error.response.data
							? error.response.data
							: error.message
					}`;
				}
			};

			reader.onerror = (error) => {
				console.error("FileReader error:", error);
				this.errorMsg = "Error reading file";
			};

			reader.readAsDataURL(file);
		},

		async forwardMessage(message) {
			try {
				// Fetch all conversations you're a part of
				const response = await this.$axios.get("/conversations");

				// Create a list of possible forward destinations
				this.forwardDestinations = response.data;

				this.messageToForward = message;
				this.showForwardModal = true;
			} catch (error) {
				console.error("Error fetching conversations:", error);
				this.errorMsg = "Failed to fetch forward destinations";
			}
		},

		async handleForward(destination) {
			try {
				await this.$axios.post(
					`/message/${this.messageToForward.messageId}/forward`,
					{
						conversationId: destination.conversationId,
					}
				);

				this.showForwardModal = false;
				this.messageToForward = null;
				await this.fetchConversationDetails();
			} catch (error) {
				console.error("Forward message error:", error);
				this.errorMsg = `Failed to forward message: ${
					error.response && error.response.data
						? error.response.data
						: error.message
				}`;
			}
		},

		async deleteMessage(message) {
			try {
				await this.$axios.delete(`/message/${message.messageId}`);
				await this.fetchConversationDetails();
			} catch (error) {
				console.error("Delete message error:", error);
				this.errorMsg = "Failed to delete message";
			}
		},
		showReactionModal(message) {
			this.currentReactionMessage = message;
		},
		async addReaction(emoji) {
			if (!this.currentReactionMessage) return;

			try {
				await this.$axios.post(
					`/message/${this.currentReactionMessage.messageId}/comment`,
					{
						emoji: emoji,
					}
				);

				await this.fetchConversationDetails();

				this.currentReactionMessage = null;
			} catch (error) {
				console.error("Add reaction error:", error);
				this.errorMsg = "Failed to add reaction";
			}
		},
		async removeReaction(message) {
			try {
				await this.$axios.delete(
					`/message/${message.messageId}/uncomment`
				);
				await this.fetchConversationDetails();
			} catch (error) {
				console.error("Remove reaction error:", error);
				this.errorMsg = "Failed to remove reaction";
			}
		},
		async leaveGroup() {
			try {
				await this.$axios.delete(
					`/group/${this.conversation.conversationId}/leave`
				);
				this.$emit("refresh");
			} catch (error) {
				console.error("Leave group error:", error);
				this.errorMsg = "Failed to leave group";
			}
		},

		async handleGroupPhotoUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			const reader = new FileReader();
			reader.onload = async (e) => {
				try {
					const base64Photo = e.target.result.split(",")[1];

					await this.$axios.put(
						`/group/${this.conversation.conversationId}/photo`,
						{
							photo: base64Photo,
						}
					);

					await this.fetchConversationDetails();
					this.showGroupPhotoInput = false;
				} catch (error) {
					console.error("Upload group photo error:", error);
					this.errorMsg = "Failed to update group photo";
				}
			};
			reader.readAsDataURL(file);
		},

		async updateGroupName() {
			if (!this.newGroupName.trim()) {
				this.errorMsg = "Group name cannot be empty";
				return;
			}

			try {
				const response = await this.$axios.put(
					`/group/${this.conversation.conversationId}/name`,
					{
						name: this.newGroupName,
					}
				);

				// Update the conversation prop directly
				this.conversation.name = this.newGroupName;

				await this.fetchConversationDetails();
				this.showGroupNameInput = false;
				this.newGroupName = "";
			} catch (error) {
				console.error("Update group name error:", error);
				this.errorMsg = "Failed to update group name";
			}
		},

		async addMemberToGroup(user) {
			try {
				await this.$axios.post(
					`/group/${this.conversation.conversationId}/add`,
					{
						username: user.username,
					}
				);

				this.showAddMemberModal = false;
				await this.fetchConversationDetails();
			} catch (error) {
				console.error("Add member error:", error);
				this.errorMsg = "Failed to add member to group";
			}
		},
		formatDate(dateString) {
			return new Date(dateString).toLocaleString();
		},
		hasUserReacted(message) {
			return (
				message.comments &&
				message.comments.some(
					(comment) => comment.userId === this.currentUserId
				)
			);
		},
	},
};
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

.message-text {
	white-space: pre-wrap;
	word-break: break-word;
}

.message-status {
	font-size: 0.8em;
	color: #666;
	margin-top: 2px;
}

.reactions .remove-reaction {
	cursor: pointer;
	margin-left: 5px;
	color: red;
}

.messages-container {
	height: calc(100vh - 200px);
}
</style>
