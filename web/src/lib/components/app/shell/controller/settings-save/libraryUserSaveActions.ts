import {
	saveLibraryFolder as saveLibraryFolderRequest,
	savePathMapping as savePathMappingRequest,
	saveUser as saveUserRequest
} from '$lib/settings/api';
import { emptyLibraryFolderForm, emptyPathMappingForm, emptyUserForm } from '$lib/settings/forms';
import { errorMessageFrom } from '../helpers';
import type { SettingsSaveContext } from './types';

export function createLibraryUserSaveActions({
	state,
	clearNotice,
	loadSettings
}: SettingsSaveContext) {
	async function saveLibraryFolder(event: SubmitEvent) {
		event.preventDefault();
		state.savingLibraryFolder = true;
		clearNotice();

		try {
			const result = await saveLibraryFolderRequest(state.libraryFolderForm);
			state.libraryFolderForm = emptyLibraryFolderForm();
			state.libraryFolders = [
				result.folder,
				...state.libraryFolders.filter((folder) => folder.id !== result.folder.id)
			];
			state.libraryScansByFolder = {
				...state.libraryScansByFolder,
				[result.folder.id]: result.scan
			};
			state.openLibraryFolderId = result.folder.id;
			state.message = `Library scan completed: ${result.scan.manualCount} pending`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not add library folder');
		} finally {
			state.savingLibraryFolder = false;
		}
	}

	async function savePathMapping(event: SubmitEvent) {
		event.preventDefault();
		state.savingPathMapping = true;
		clearNotice();

		try {
			const mapping = await savePathMappingRequest(state.pathMappingForm);
			state.pathMappingForm = emptyPathMappingForm();
			state.pathMappings = [
				mapping,
				...state.pathMappings.filter((item) => item.id !== mapping.id)
			];
			state.message = 'Path mapping saved';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save path mapping');
		} finally {
			state.savingPathMapping = false;
		}
	}

	async function saveUser(event: SubmitEvent) {
		event.preventDefault();
		state.savingUser = true;
		clearNotice();

		try {
			await saveUserRequest(state.userForm);
			state.userForm = emptyUserForm();
			state.message = 'User saved';
			await loadSettings();
			if (state.currentUser && state.users.some((user) => user.id === state.currentUser?.id)) {
				const updatedUser = state.users.find((user) => user.id === state.currentUser?.id);
				if (updatedUser) {
					state.currentUser = {
						id: updatedUser.id,
						username: updatedUser.username,
						role: updatedUser.role
					};
				}
			}
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save user');
		} finally {
			state.savingUser = false;
		}
	}

	return {
		saveLibraryFolder,
		savePathMapping,
		saveUser
	};
}
