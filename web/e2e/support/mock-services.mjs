import http from 'node:http';

const port = Number(process.env.TEST_MOCK_PORT ?? 18180);

const server = http.createServer((request, response) => {
	const url = new URL(request.url ?? '/', `http://${request.headers.host}`);
	if (url.pathname === '/health') return json(response, { status: 'ok' });
	if (url.pathname === '/tmdb/3/search/movie') return json(response, tmdbMovieSearch);
	if (url.pathname === '/tmdb/3/movie/936075') return json(response, tmdbMovieDetails);
	if (url.pathname === '/torznab/api') return torznab(response, url);
	response.writeHead(404).end();
});

server.listen(port, '0.0.0.0');

function json(response, body) {
	response.writeHead(200, { 'content-type': 'application/json' });
	response.end(JSON.stringify(body));
}

function torznab(response, url) {
	response.writeHead(200, { 'content-type': 'application/xml' });
	response.end(url.searchParams.get('t') === 'caps' ? torznabCaps : torznabSearch);
}

const tmdbMovieSearch = {
	page: 1,
	results: [
		{
			id: 936075,
			title: 'Example Movie',
			release_date: '2026-02-14',
			overview: 'A realistic local metadata result.'
		}
	],
	total_pages: 1,
	total_results: 1
};

const tmdbMovieDetails = {
	id: 936075,
	title: 'Example Movie',
	release_date: '2026-02-14',
	overview: 'A realistic local metadata detail response.'
};

const torznabCaps = `<caps>
  <server title="Local Torznab Mock" version="1.0"/>
  <limits max="100" default="50"/>
  <categories>
    <category id="2000" name="Movies"><subcat id="2040" name="HD"/></category>
    <category id="5000" name="TV"><subcat id="5070" name="Anime"/></category>
  </categories>
</caps>`;

const torznabSearch = `<rss><channel>
  <title>Local releases</title>
  <item><title>Example.Movie.2026.1080p.WEB-DL</title></item>
</channel></rss>`;
