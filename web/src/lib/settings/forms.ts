import type {
	LibraryFolderForm,
	LibraryFolderRequest,
	PathMappingForm,
	PathMappingRequest
} from './types';

export {
	customFormatFormFromFormat,
	emptyCustomFormatForm,
	normalizeCustomFormatForm
} from './customFormatForms';
export {
	emptyLanguageForm,
	languageFormFromLanguage,
	normalizeLanguageForm,
	normalizeLanguageUpdateForm
} from './languageForms';
export {
	defaultAudioTarget,
	defaultSubtitleTarget,
	emptyMediaProfileForm,
	mediaProfileFormFromProfile,
	normalizeMediaProfileForm
} from './mediaProfileForms';
export {
	emptyUserForm,
	normalizeUserCreateForm,
	normalizeUserUpdateForm,
	userFormFromUser
} from './userForms';
export {
	downloadClientFormFromClient,
	downloadClientProtocolForType,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyMetadataProviderForm,
	emptySubtitleProviderForm,
	indexerFormFromIndexer,
	metadataProviderFormFromProvider,
	normalizeDownloadClientForm,
	normalizeIndexerForm,
	normalizeMetadataProviderForm,
	normalizeSubtitleProviderForm,
	subtitleProviderFormFromProvider
} from './integrationForms';

export function emptyLibraryFolderForm(): LibraryFolderForm {
	return {
		path: '',
		kind: 'movie'
	};
}

export function emptyPathMappingForm(): PathMappingForm {
	return {
		clientPath: '',
		appPath: ''
	};
}

export function normalizeLibraryFolderForm(form: LibraryFolderForm): LibraryFolderRequest {
	return {
		path: form.path.trim(),
		kind: form.kind
	};
}

export function normalizePathMappingForm(form: PathMappingForm): PathMappingRequest {
	return {
		clientPath: form.clientPath.trim(),
		appPath: form.appPath.trim()
	};
}
