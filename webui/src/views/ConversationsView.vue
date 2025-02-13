<template>
	<div class="container-fluid">
		<div class="row">
			<!-- Conversations List -->
			<div class="col-md-4 border-end">
				<div
					class="d-flex justify-content-between align-items-center p-3 border-bottom"
				>
					<h4 class="mb-0">Conversations</h4>
					<div class="d-flex">
						<button
							class="btn btn-outline-primary me-2"
							@click="showCreateGroupModal = true"
						>
							Create Group
						</button>
						<button
							class="btn btn-outline-secondary"
							@click="showUserSearchModal = true"
						>
							Search Users
						</button>
					</div>
				</div>

				<LoadingSpinner :loading="loading">
					<ErrorMsg v-if="errorMsg" :msg="errorMsg" />

					<div
						v-for="conv in conversations"
						:key="conv.conversationId"
						class="conversation-item d-flex align-items-center p-3 border-bottom"
						@click="selectConversation(conv)"
						:class="{
							'active-conversation':
								selectedConversation && selectedConversation.conversationId ===
								conv.conversationId,
						}"
					>
						<div class="conversation-avatar me-3">
							<img
								:src="conv.photo || '/default-avatar.png'"
								class="rounded-circle"
								width="50"
								height="50"
								alt="Conversation Avatar"
							/>
						</div>
						<div class="conversation-details flex-grow-1">
							<div class="d-flex justify-content-between">
								<h6 class="mb-1">{{ conv.name }}</h6>
								<small class="text-muted">
									{{ formatDate(conv.lastMessageTime) }}
								</small>
							</div>
							<p class="text-muted mb-0 text-truncate">
								{{
									conv.isPhoto
										? "📸 Photo"
										: conv.lastMessageText
								}}
							</p>
						</div>
					</div>
				</LoadingSpinner>
			</div>

			<!-- Conversation Details -->
			<div class="col-md-8">
				<ConversationDetail
					v-if="selectedConversation"
					:conversation="selectedConversation"
					@refresh="fetchConversations"
				/>
				<div
					v-else
					class="d-flex justify-content-center align-items-center h-100 text-muted"
				>
					Select a conversation to view messages
				</div>
			</div>
		</div>

		<!-- User Search Modal -->
		<UserSearchModal
			v-if="showUserSearchModal"
			@close="showUserSearchModal = false"
			@start-conversation="startDirectConversation"
		/>

		<!-- Create Group Modal -->
		<CreateGroupModal
			v-if="showCreateGroupModal"
			@close="showCreateGroupModal = false"
			@group-created="onGroupCreated"
		/>
	</div>
</template>

<script>
import UserSearchModal from "../components/UserSearchModal.vue";
import CreateGroupModal from "../components/CreateGroupModal.vue";
import ConversationDetail from "../components/ConversationDetail.vue";

export default {
	components: {
		UserSearchModal,
		CreateGroupModal,
		ConversationDetail,
	},
	data() {
		return {
			conversations: [],
			selectedConversation: null,
			loading: false,
			errorMsg: null,
			showUserSearchModal: false,
			showCreateGroupModal: false,
			refreshInterval: null,
		};
	},
	created() {
		this.fetchConversations();
		// Auto-refresh every 5 seconds
		this.refreshInterval = setInterval(this.fetchConversations, 5000);
	},
	beforeUnmount() {
		if (this.refreshInterval) {
			clearInterval(this.refreshInterval);
		}
	},
	methods: {
		async fetchConversations() {
			this.loading = true;
			try {
				const response = await this.$axios.get("/conversations");
				this.conversations = response.data;
				this.loading = false;
			} catch (error) {
				console.error("Fetch conversations error:", error);
				this.errorMsg = "Failed to fetch conversations";
				this.loading = false;
			}
		},
		selectConversation(conversation) {
			this.selectedConversation = conversation;
		},
		formatDate(dateString) {
			return new Date(dateString).toLocaleString();
		},
		async startDirectConversation(user) {
			// Prompt user for the first message
			const firstMessage = prompt(
				`Send message to ${user.username}:`,
				""
			);

			// Check if user cancelled or entered an empty message
			if (firstMessage === null || firstMessage.trim() === "") {
				return;
			}

			try {
				console.log("Starting conversation with user:", user);

				// Send initial message with user's input
				const response = await this.$axios.post("/message", {
					text: firstMessage,
					recipientUsername: user.username,
				});

				// Refresh conversations list
				await this.fetchConversations();

				// Find and select the new conversation
				const newConversation = this.conversations.find(
					(conv) => !conv.isGroup && conv.name === user.username
				);
				if (newConversation) {
					this.selectConversation(newConversation);
				}

				this.showUserSearchModal = false;
			} catch (error) {
				console.error(
					"Start conversation error:",
					error.response || error
				);
				this.errorMsg =
					"Failed to start conversation: " +
					(error.response && error.response.data)
						? error.response.data
						: error.message;
			}
		},

		onGroupCreated(group) {
			this.showCreateGroupModal = false;
			this.fetchConversations();
			this.selectConversation(group);
		},
	},
};
</script>

<style scoped>
.conversation-item {
	cursor: pointer;
	transition: background-color 0.2s;
}

.conversation-item:hover {
	background-color: #f8f9fa;
}

.active-conversation {
	background-color: #e9ecef;
}

.conversation-avatar img {
	object-fit: cover;
}
</style>
