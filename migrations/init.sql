select
    'create database auth_svc'
where
    not exists (
        select
        from
            pg_database
        where
            datname = 'auth_svc'
    );
