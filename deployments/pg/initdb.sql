CREATE TABLE public.url (
    id bigserial PRIMARY KEY,
    url text UNIQUE,
    url_compressed text UNIQUE
);