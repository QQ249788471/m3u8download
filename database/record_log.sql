create table `record_log` (
    `id` int unsigned not null primary key auto_increment,
    `record_node_id` int unsigned not null comment "收录处理表id",
    `level` tinyint unsigned not null default 1 comment "日志级别，1 INFO，2 WARN，3 ERROR",
    `content` text not null comment "内容",
    `created_at` timestamp not null default current_timestamp,

    key(`record_node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment "收录日志记录";