create database cryptocompare;

drop table if exists data;
create table data
(
    _id             serial not null primary key,
    fromsym         bigint not null,
    tosym           bigint not null,
    change24hour    double precision,
    changepct24hour double precision,
    open24hour      double precision,
    volume24hour    double precision,
    low24hour       double precision,
    high24hour      double precision,
    price           double precision,
    supply          double precision,
    mktcap          double precision,
    lastupdate      text   not null,
    displaydataraw  text
);

drop table if exists symbols;
create table symbols
(
    _id serial not null constraint symbols_pk primary key,
    symbol varchar(64) default '' not null constraint symbols_symbol_uindex unique,
    unicode char
);

insert into symbols(symbol, unicode)
values ('USDT','₮'),
       ('BTC','₿'),
       ('ETH','⟠'),
       ('USD','$'),
       ('XRP','✕'),
       ('LTC','Ł'),
       ('EUR','€'),
       ('GBP','£'),
       ('JPY','¥');

drop table if exists session cascade;
create table session
(
    _id serial not null primary key,
    task_name varchar(64) not null,
    interval integer default 60 not null
);

create unique index session_task_name_uindex
    on session (task_name);
