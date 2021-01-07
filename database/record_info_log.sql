create table `record_info_log` (
    `id` int unsigned not null primary key auto_increment,
    `record_info_id` int unsigned not null comment "收录详情表id",
    `content` text not null default "" comment "详细内容",
    `node_name` varchar (128) not null comment "节点名称",
    `node_ip` varchar(20) not null comment "节点ip",

    key(`record_info_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment "收录详情日志表";