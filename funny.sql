select band.name, string_agg(album.name, ', ') from band inner join album on album.band_id = band.id group by band.id;

select album.name, string_agg(track.name, ', ') from album inner join track on track.album_id = album.id group by album.id;
