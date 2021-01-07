create table `record_info` (
    `id` int unsigned not null primary key auto_increment,
    `record_node_id` int unsigned not null comment "收录处理表id",
    `uri` varchar(512) not null comment "ts uri",
    `path` varchar(512) not null comment "文件保存路径",
    `status` tinyint unsigned not null default 0 comment "任务状态，0 等待下载，1 下载中，2 下载完成，3 下载失败",
    `retry` tinyint unsigned not null default 0 comment "重试次数",
    `size` int not null default 0 comment "文件大小，单位字节",
    `time` varchar(10) not null default "0" comment "下载耗时，单位秒",
    `created_at` timestamp not null default current_timestamp,
    `updated_at` timestamp not null on update current_timestamp,

    key(`record_node_id`),
    key(`created_at`),
    key(`retry`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment "收录详情表";