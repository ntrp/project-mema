import type { components } from '$lib/api/generated/schema';

type S = components['schemas'];

export type DLNASettings = S['DLNASettings'];
export type DLNASettingsRequest = S['DLNASettingsRequest'];
export type DLNAStatus = S['DLNAStatus'];
export type DLNAClientDiagnostic = S['DLNAClientDiagnostic'];
export type DLNAInterfaceDiagnostic = S['DLNAInterfaceDiagnostic'];
export type DLNAStreamDiagnostic = S['DLNAStreamDiagnostic'];
export type DLNARendererProfile = S['DLNARendererProfile'];
export type DLNARendererProfileRequest = S['DLNARendererProfileRequest'];
export type DLNARendererProfileCreateRequest = S['DLNARendererProfileCreateRequest'];
export type DLNARendererProfileCloneRequest = S['DLNARendererProfileCloneRequest'];
export type DLNARendererDeviceOverride = S['DLNARendererDeviceOverride'];
export type DLNARendererDeviceOverrideRequest = S['DLNARendererDeviceOverrideRequest'];
export type DLNAProfileMatchTraceRequest = S['DLNAProfileMatchTraceRequest'];
export type DLNAProfileMatchTraceResponse = S['DLNAProfileMatchTraceResponse'];
export type DLNADeliveryTraceRequest = S['DLNADeliveryTraceRequest'];
export type DLNADeliveryTraceResponse = S['DLNADeliveryTraceResponse'];
export type SystemJob = S['SystemJob'];
export type SystemJobListResponse = S['SystemJobListResponse'];
export type SystemJobsOverviewResponse = S['SystemJobsOverviewResponse'];
export type SystemJobSchedule = S['SystemJobSchedule'];
export type SystemJobExecution = S['SystemJobExecution'];
export type SystemJobExecutionListResponse = S['SystemJobExecutionListResponse'];
export type SystemJobExecutionLog = S['SystemJobExecutionLog'];
export type SystemJobHistorySettings = S['SystemJobHistorySettings'];
