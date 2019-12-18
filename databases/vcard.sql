-- auto-generated definition
create table users
(
    id            bigint unsigned auto_increment
        primary key,
    username      varchar(32)                        not null,
    nickname      varchar(32)                        not null comment '昵称',
    password      varchar(255)                       not null,
    avatar        varchar(512)                       not null comment '用户头像',
    lastlogintime datetime                           null comment '最后登录时间',
    createtime    datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    constraint users_username_unique
        unique (username)
)
    collate = utf8mb4_unicode_ci;




-- auto-generated definition
create table user_cards
(
    id          varchar(32)       not null,
    userid      bigint            not null comment '用户id',
    name        varchar(32)       not null comment '姓名',
    phone       varchar(32)       null comment '电话',
    sex         tinyint default 0 not null comment '性别:1男，2:女',
    company     varchar(256)      null comment '公司',
    province    varchar(32)       null comment '省',
    city        varchar(32)       null comment '市',
    district    varchar(32)       null comment '区',
    height      int     default 0 not null comment '身高，单位cm',
    weight      int     default 0 not null comment '体重，单位kg',
    bust        int     default 0 not null comment '胸围，单位cm',
    waist       int     default 0 not null comment '腰围，单位cm',
    hipline     int     default 0 not null comment '臀围，单位cm',
    cover       varchar(256)      null comment '封面图',
    profile     varchar(512)      null comment '简介',
    video       varchar(128)      null comment '视频',
    isavailable tinyint default 1 not null comment '是否有效1:有效，0:无效',
    created_at  timestamp         null,
    updated_at  timestamp         null,
    constraint user_cards_id_unique
        unique (id)
)
    collate = utf8mb4_unicode_ci;

alter table user_cards
    add primary key (id);





INSERT INTO user_cards (id, userid, name, phone, sex, company, province, city, district, height, weight, bust, waist, hipline, cover, profile, video, isavailable, created_at, updated_at) VALUES ('5d42355020ba11ea8c6eacde48001122', 1, '陈s', '17611168388', 1, '', '', '', '', 0, 0, 0, 0, 0, '', '', '', 1, null, null);
INSERT INTO user_cards (id, userid, name, phone, sex, company, province, city, district, height, weight, bust, waist, hipline, cover, profile, video, isavailable, created_at, updated_at) VALUES ('769ae86d4b2d53bd88beb48287cbe39e', 2, 'yonds', '17611168388', 1, '', '', '', '', 0, 0, 0, 0, 0, '', '', '', 1, null, null);
INSERT INTO users (id, username, nickname, password, avatar, lastlogintime, createtime) VALUES (1, 'chen', '老七', 'xxxxxxxxxxx', '', null, '2019-12-04 13:37:31');
INSERT INTO users (id, username, nickname, password, avatar, lastlogintime, createtime) VALUES (2, 'yond', 'yond', 'xxxxxxxxx', '', null, '2019-12-04 13:45:16');