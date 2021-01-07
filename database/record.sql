create table `record` (
    `id` int unsigned not null primary key auto_increment,
    `live` varchar(512) not null comment "直播地址",
    `start_time` int unsigned not null comment "开始收录时间，周期性节目只有时间，日期取 2020-01-01",
    `end_time` int unsigned not null comment "结束收录时间，周期性节目只有时间，日期取 2020-01-01",
    `is_repeat` tinyint unsigned not null default 1 comment "1 仅一次，2 周期性",
    `weekday` set("1","2","3","4","5","6","7") not null default "1" comment "1-7 周一到周日",
    `status` tinyint unsigned not null default 1 comment"开启状态，1 开启，2 关闭",
    `created_at` timestamp not null default current_timestamp,
    `updated_at` timestamp not null on update current_timestamp,
    `deleted_at` int unsigned not null default 0,

    key(`start_time`),
    key(`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment "收录任务表";