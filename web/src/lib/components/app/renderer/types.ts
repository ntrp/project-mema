import type { components } from '$lib/api/generated/schema';

type Schemas = components['schemas'];

export type Renderer = Schemas['Renderer'];
export type RendererStatus = Schemas['RendererStatus'];
export type RendererTransportRequest = Schemas['RendererTransportRequest'];
