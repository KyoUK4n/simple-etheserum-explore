create table transaction
(
    tx_hash      varchar(66)  not null comment '交易hash',
    tx_index     int          not null default 0,
    `from`       varchar(42)  not null default '' comment '发送者',
    `to`         varchar(42)  not null default '' comment '接收者',
    value        varchar(78)  not null default '0' comment '发送的eth(Wei),用varchar防止bigint溢出',
    gas_limit    bigint       not null default 0,
    gas_price    bigint       not null default 0,
    gas_used     bigint       not null default 0,
    nonce        int          not null default 0,
    status       tinyint      not null default 0,
    block_number bigint       not null default 0,
    block_hash   varchar(66)  not null default '',
    tx_timestamp datetime     null,
    PRIMARY KEY (tx_hash),
    KEY idx_from (`from`),
    KEY idx_to (`to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;