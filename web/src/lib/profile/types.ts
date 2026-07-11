import type { components } from '$lib/api/generated/schema';

type Schemas = components['schemas'];

export type UserProfile = Schemas['UserProfile'];
export type UserProfileUpdateRequest = Schemas['UserProfileUpdateRequest'];
