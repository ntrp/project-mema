-- +goose Up
alter table app.download_clients
    add column if not exists protocol text;

update app.download_clients
set protocol = case
    when type = 'sabnzbd' then 'usenet'
    when type = 'transmission' then 'torrent'
    else 'torrent'
end
where protocol is null
    or protocol not in ('torrent', 'usenet')
    or (type = 'sabnzbd' and protocol <> 'usenet')
    or (type = 'transmission' and protocol <> 'torrent');

alter table app.download_clients
    alter column protocol set default 'torrent',
    alter column protocol set not null;

alter table app.download_clients drop constraint if exists download_clients_protocol_check;
alter table app.download_clients
    add constraint download_clients_protocol_check check (protocol in ('torrent', 'usenet'));

alter table app.download_clients drop constraint if exists download_clients_type_protocol_check;
alter table app.download_clients
    add constraint download_clients_type_protocol_check
    check (
        (type = 'transmission' and protocol = 'torrent')
        or (type = 'sabnzbd' and protocol = 'usenet')
    );

-- +goose Down
alter table app.download_clients drop constraint if exists download_clients_type_protocol_check;
alter table app.download_clients drop constraint if exists download_clients_protocol_check;
alter table app.download_clients drop column if exists protocol;
