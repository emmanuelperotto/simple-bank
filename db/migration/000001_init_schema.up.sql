-- auto-generated definition
create table accounts
(
    id         bigserial
        primary key,
    owner      varchar                 not null,
    balance    bigint                  not null,
    currency   varchar                 not null,
    created_at timestamp default now() not null
);

alter table accounts
    owner to root;

create index accounts_owner_idx
    on accounts (owner);

-- auto-generated definition
create table entries
(
    id         bigserial
        primary key,
    account_id bigint                  not null
        references accounts,
    amount     bigint                  not null,
    created_at timestamp default now() not null
);

comment on column entries.amount is 'can be negative or positive';

alter table entries
    owner to root;

create index entries_account_id_idx
    on entries (account_id);

-- auto-generated definition
create table transfers
(
    id              bigserial
        primary key,
    from_account_id bigint                  not null
        references accounts,
    to_account_id   bigint                  not null
        references accounts,
    amount          bigint                  not null,
    created_at      timestamp default now() not null
);

comment on column transfers.amount is 'must be positive';

alter table transfers
    owner to root;

create index transfers_from_account_id_idx
    on transfers (from_account_id);

create index transfers_to_account_id_idx
    on transfers (to_account_id);

create index transfers_from_account_id_to_account_id_idx
    on transfers (from_account_id, to_account_id);



