-- =========================================================
-- 表名：device_state
-- 说明：
--   设备最新状态表（每台设备 1 行）
--   用于后台列表/看板/告警快速读取，不扫 heartbeat/event 大表
--
-- 设计原则：
--   1) 一设备一行，device_id 做主键
--   2) 状态由应用层在写入 heartbeat/event 时同步刷新（或异步任务刷新）
--   3) created_at / updated_at 由 GORM 维护（DB 不加 trigger）
-- =========================================================

create table if not exists device_state
(
    -- 与 device 一一对应
    device_id         bigint primary key,

    -- 最近一次收到“任何上报”的时间（心跳/事件都算）
    last_seen_at      timestamptz,

    -- 最近一次心跳/事件发生时间（发生时间）
    last_heartbeat_at timestamptz,
    last_heartbeat_id bigint,
    last_event_at     timestamptz,

    -- 当前关键状态（先保留少量字段，后续再加）
    door_open         boolean,
    signal_strength   smallint,
    battery_level     smallint,

    -- 当前称重相关（你后面做温漂 baseline 慢跟随会用到）
    -- 注意：具体单位/精度你定，numeric 更稳；如果不确定先放 payload
    weight            numeric(12, 3),
    baseline          numeric(12, 3),

    -- 最近异常（可选）
    last_error_code   text,
    last_error_at     timestamptz,

    -- 扩展状态（兜底，避免一开始想不全）
    payload           jsonb       not null default '{}'::jsonb,

    -- 审计字段（由 GORM 自动维护）
    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now()
);

-- 常用查询：找离线/久未上报（last_seen_at 越旧越靠前）
create index if not exists device_state_last_seen_at_idx
    on device_state (last_seen_at);
