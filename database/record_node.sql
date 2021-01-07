create table `record_node` (
    `id` int unsigned not null primary key auto_increment,
    `record_id` int unsigned not null comment "收录任务id",
    `record_start_time` int unsigned not null comment "收录任务的开始时间",
    `node_name` varchar (128) not null comment "节点名称",
    `node_ip` varchar(20) not null comment "节点ip",
    `status` tinyint unsigned not null default 1 comment "任务状态，1 进行中，2 已完成",
    `created_at` timestamp not null default current_timestamp,
    `updated_at` timestamp not null on update current_timestamp,
    `deleted_at` int not null default 0,

    unique(`record_id`, `record_start_time`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment "收录处理表";