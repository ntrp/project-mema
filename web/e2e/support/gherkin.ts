import { readFileSync, readdirSync } from 'node:fs';
import { resolve } from 'node:path';

export interface Scenario {
	id: string;
	feature: string;
	name: string;
	tags: string[];
	steps: string[];
}

export function requireScenario(id: string, requiredTag?: string) {
	const scenario = loadScenarios().find((candidate) => candidate.id === id);
	if (!scenario) throw new Error(`Scenario ${id} was not found`);
	if (requiredTag && !hasTag(scenario, requiredTag)) {
		throw new Error(`Scenario ${id} is missing @${requiredTag}`);
	}
	return scenario;
}

function loadScenarios() {
	const dir = resolve(process.cwd(), '..', 'features', 'behavior');
	return readdirSync(dir)
		.filter((file) => file.endsWith('.feature'))
		.flatMap((file) => parseFeature(readFileSync(resolve(dir, file), 'utf8')));
}

function parseFeature(source: string) {
	let feature = '';
	let pendingTags: string[] = [];
	let current: Scenario | undefined;
	const scenarios: Scenario[] = [];
	for (const rawLine of source.split(/\r?\n/)) {
		const line = rawLine.trim();
		if (!line || line.startsWith('#')) continue;
		if (line.startsWith('Feature:')) {
			feature = line.replace('Feature:', '').trim();
		} else if (line.startsWith('@')) {
			pendingTags = line.split(/\s+/);
		} else if (line.startsWith('Scenario:')) {
			if (current) scenarios.push(current);
			current = {
				id: scenarioId(pendingTags),
				feature,
				name: line.replace('Scenario:', '').trim(),
				tags: pendingTags,
				steps: []
			};
			pendingTags = [];
		} else if (current && /^(Given|When|Then|And|But)\s/.test(line)) {
			current.steps.push(line);
		}
	}
	if (current) scenarios.push(current);
	return scenarios;
}

function hasTag(scenario: Scenario, tag: string) {
	const normalized = tag.replace(/^@/, '');
	return scenario.tags.some((candidate) => candidate.replace(/^@/, '') === normalized);
}

function scenarioId(tags: string[]) {
	return tags.map((tag) => tag.replace(/^@/, '')).find((tag) => tag.startsWith('SCN-')) ?? '';
}
