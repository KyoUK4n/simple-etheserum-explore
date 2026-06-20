create table event_log
(
    id           bigint auto_increment comment '自增主键',
    tx_hash      varchar(66)  not null default '' comment '交易hash，0x开头，方便转换',
    address      varchar(42)  not null default '' comment '合约地址',
    event_name   varchar(50)  not null default '' comment '事件名称',
    block_number bigint       not null default 0  comment '区块号',
    log_index    int          not null default 0  comment '日志索引',
    topics       json         null comment '索引字段',
    data         json         null comment '非索引字段',
    tx_timestamp datetime     not null default CURRENT_TIMESTAMP comment '交易时间',
    PRIMARY KEY (id),
    KEY idx_address (address),
    KEY idx_tx_hash (tx_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;