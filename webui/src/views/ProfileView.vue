<template>
	<div class="container mt-4">
		<div class="row justify-content-center">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h3>Profile Settings</h3>
					</div>
					<div class="card-body">
						<div class="mb-3">
							<label class="form-label">Update Username</label>
							<div class="input-group">
								<input
									type="text"
									class="form-control"
									v-model="newUsername"
									placeholder="Enter new username"
								/>
								<button
									class="btn btn-primary"
									@click="updateUsername"
								>
									Update Username
								</button>
							</div>
						</div>
						<div class="mb-3">
							<label class="form-label"
								>Update Profile Photo</label
							>
							<div class="text-center mb-3">
								<img
									:src="
										currentUser.profilePhoto ||
										'/default-pic.jpg'
									"
									class="rounded-circle mb-3"
									width="150"
									height="150"
									alt="Profile Photo"
								/>
							</div>
							<div class="input-group">
								<input
									type="file"
									class="form-control"
									@change="handlePhotoUpload"
									accept="image/*"
								/>
							</div>
						</div>
						<ErrorMsg v-if="errorMsg" :msg="errorMsg" />
						<div v-if="successMsg" class="alert alert-success">
							{{ successMsg }}
						</div>
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
			currentUser: JSON.parse(localStorage.getItem("user")) || {},
			newUsername: "",
			errorMsg: null,
			successMsg: null,
		};
	},
	methods: {
		async updateUsername() {
			if (!this.newUsername) {
				this.errorMsg = "Username cannot be empty";
				return;
			}
			try {
				const response = await this.$axios.put(
					`/user/${this.currentUser.username}/setmyusername`,
					{
						username: this.newUsername,
					}
				);
				const updatedUser = response.data;
				localStorage.setItem("user", JSON.stringify(updatedUser));
				this.currentUser = updatedUser;
				this.successMsg = "Username updated successfully";
				this.errorMsg = null;
				this.newUsername = "";
			} catch (error) {
				this.errorMsg =
					error.response && error.response.data
						? error.response.data
						: "Failed to update username";
				this.successMsg = null;
			}
		},
		async handlePhotoUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			const reader = new FileReader();
			reader.onload = async (e) => {
				try {
					const base64Photo = e.target.result.split(",")[1];

					await this.$axios.put(
						`/user/${this.currentUser.username}/photo`,
						{
							photo: base64Photo,
						}
					);

					const updatedUser = {
						...this.currentUser,
						profilePhoto: e.target.result, // Store the complete data URL
					};
					localStorage.setItem("user", JSON.stringify(updatedUser));
					this.currentUser = updatedUser;
					this.successMsg = "Profile photo updated successfully";
					this.errorMsg = null;
				} catch (error) {
					console.error("Update photo error:", error);
					this.errorMsg =
						error.response?.data ||
						"Failed to update profile photo";
					this.successMsg = null;
				}
			};
			reader.readAsDataURL(file);
		},
	},
};
</script>
