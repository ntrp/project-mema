import {
	saveLibraryFolder as saveLibraryFolderRequest,
	savePathMapping as savePathMappingRequest
} from '$lib/settings/api';
import { saveUser as saveUserRequest } from '$lib/settings/domains/users';
import { emptyLibraryFolderForm, emptyPathMappingForm, emptyUserForm } from '$lib/settings/forms';
import { errorMessageFrom } from '../helpers';
import type { SettingsSaveContext } from './types';

export function createLibraryUserSaveActions({
	state,
	clearNotice,
	loadSettings,
	users,
	upsertLibraryFolder,
	upsertPathMapping,
	upsertLibraryScan,
	runMutation = (command) => command()
}: SettingsSaveContext) {
	async function saveLibraryFolder(event: SubmitEvent) {
		event.preventDefault();
		state.savingLibraryFolder = true;
		clearNotice();

		try {
			const result = await runMutation(() => saveLibraryFolderRequest(state.libraryFolderForm));
			state.libraryFolderForm = emptyLibraryFolderForm();
			upsertLibraryFolder?.(result.folder);
			upsertLibraryScan?.(result.scan);
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
			const mapping = await runMutation(() => savePathMappingRequest(state.pathMappingForm));
			state.pathMappingForm = emptyPathMappingForm();
			upsertPathMapping?.(mapping);
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
			await runMutation(() => saveUserRequest(state.userForm));
			state.userForm = emptyUserForm();
			state.message = 'User saved';
			await loadSettings();
			if (state.currentUser && users) {
				const updatedUser = users().find((user) => user.id === state.currentUser?.id);
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
