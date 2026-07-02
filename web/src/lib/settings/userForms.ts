import type { ManagedUser, UserCreateRequest, UserForm, UserUpdateRequest } from './types';

export function emptyUserForm(): UserForm {
	return {
		username: '',
		password: '',
		role: 'user'
	};
}

export function userFormFromUser(user: ManagedUser): UserForm {
	return {
		id: user.id,
		username: user.username,
		password: '',
		role: user.role
	};
}

export function normalizeUserCreateForm(form: UserForm): UserCreateRequest {
	return {
		username: form.username.trim(),
		password: form.password,
		role: form.role
	};
}

export function normalizeUserUpdateForm(form: UserForm): UserUpdateRequest {
	return {
		username: form.username.trim(),
		password: optionalString(form.password),
		role: form.role
	};
}

function optionalString(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}
