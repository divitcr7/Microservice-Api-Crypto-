create database if not exists cryptocompare /*!40100 DEFAULT CHARACTER SET utf8 */;
use cryptocompare;

drop table if exists data;
create table data
(
    _id             int auto_increment primary key,
    fromSym         int        not null,
    toSym           int        not null,
    change24hour    double     null,
    changepct24hour double     null,
    open24hour      double     null,
    volume24hour    double     null,
    low24hour       double     null,
    high24hour      double     null,
    price           double     null,
    supply          double     null,
    mktcap          double     null,
    lastupdate      mediumtext not null,
    displaydataraw  text       null
)
    collate = utf8_general_ci;

drop table if exists symbols;
create table symbols
(
    _id int auto_increment primary key,
    symbol varchar(64) collate latin1_swedish_ci default '' not null,
    unicode char null,
    constraint symbols_symbol_uindex unique (symbol)
) collate = utf8_general_ci;

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
    _id       int auto_increment primary key,
    task_name varchar(64) not null default '',
    `interval`  integer default 60 not null
);

create unique index session_task_name_uindex
    on session (task_name);
