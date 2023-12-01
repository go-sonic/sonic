create table if not exists attachment
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    file_key    varchar(2047)            not null,
    height      int           default 0  not null,
    media_type  varchar(127)  default '' not null,
    name        varchar(255)             not null,
    path        varchar(1023)            not null,
    size        bigint                   not null,
    suffix      varchar(50)   default '' not null,
    thumb_path  varchar(1023) default '' not null,
    type        int           default 0  not null,
    width       int           default 0  not null,
    index attachment_create_time (create_time),
    index attachment_media_type (media_type)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists category
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    description varchar(100)  default '' not null,
    type        tinyint       default 0  not null,
    name        varchar(255)             not null,
    parent_id   int           default 0  not null,
    password    varchar(255)  default '' not null,
    slug        varchar(255)  default '' not null,
    thumbnail   varchar(1023) default '' not null,
    priority    int           default 0  not null,
    unique index uniq_category_slug (slug),
    index category_name (name),
    index category_parent_id (parent_id)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists comment_black
(
    id          int auto_increment primary key,
    create_time datetime(6)  not null,
    update_time datetime(6)  null,
    ban_time    datetime(6)  not null,
    ip_address  varchar(127) not null
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists comment
(
    id                 int auto_increment primary key,
    type               int          default 0  not null,
    create_time        datetime(6)             not null,
    update_time        datetime(6)             null,
    allow_notification tinyint(1)   default 1  not null,
    author             varchar(50)             not null,
    author_url         varchar(511) default '' not null,
    content            varchar(1023)           not null,
    email              varchar(255)            not null,
    gravatar_md5       varchar(127) default '' not null,
    ip_address         varchar(127) default '' not null,
    is_admin           tinyint(1)   default 0  not null,
    parent_id          int       default 0  not null,
    post_id            int                     not null,
    status             int          default 0  not null,
    top_priority       int          default 0  not null,
    user_agent         varchar(511) default '' not null,
    likes              int          default 0 not null ,
    index comment_parent_id (parent_id),
    index comment_post_id (post_id),
    index comment_type_status (type, status)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists flyway_schema_history
(
    installed_rank int                                 not null
        primary key,
    version        varchar(50)                         null,
    description    varchar(200)                        not null,
    type           varchar(20)                         not null,
    script         varchar(1000)                       not null,
    checksum       int                                 null,
    installed_by   varchar(100)                        not null,
    installed_on   timestamp default CURRENT_TIMESTAMP not null,
    execution_time int                                 not null,
    success        tinyint(1)                          not null,
    index flyway_schema_history_s_idx (success)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists journal
(
    id             int auto_increment primary key,
    create_time    datetime(6)      not null,
    update_time    datetime(6)      null,
    content        text             not null,
    likes          bigint default 0 not null,
    source_content longtext         not null,
    type           int    default 0 not null
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists link
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    description varchar(255)  default '' not null,
    logo        varchar(1023) default '' not null,
    name        varchar(255)             not null,
    priority    int           default 0  not null,
    team        varchar(255)  default '' not null,
    url         varchar(1023)            not null,
    index link_name (name)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists log
(
    id          bigint auto_increment primary key,
    create_time datetime(6)             not null,
    update_time datetime(6)             null,
    content     varchar(1023)           not null,
    ip_address  varchar(127) default '' not null,
    log_key     varchar(1023)           not null,
    type        int                     not null,
    index log_create_time (create_time)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists menu
(
    id          int auto_increment primary key,
    create_time datetime(6)                  not null,
    update_time datetime(6)                  null,
    icon        varchar(50)  default ''      not null,
    name        varchar(50)                  not null,
    parent_id   int          default 0       not null,
    priority    int          default 0       not null,
    target      varchar(20)  default '_self' not null,
    team        varchar(255) default ''      not null,
    url         varchar(1023)                not null,
    index menu_name (name),
    index menu_parent_id (parent_id)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists meta
(
    id          int auto_increment primary key,
    type        int default 0 not null,
    create_time datetime(6)   not null,
    update_time datetime(6)   null,
    meta_key    varchar(255)  not null,
    post_id     int           not null,
    meta_value  varchar(1023) not null
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists `option`
(
    id           int auto_increment primary key,
    create_time  datetime(6)   not null,
    update_time  datetime(6)   null,
    option_key   varchar(100)  not null,
    type         int default 0 not null,
    option_value longtext      not null
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists photo
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    description varchar(255)  default '' not null,
    location    varchar(255)  default '' not null,
    name        varchar(255)             not null,
    take_time   datetime(6)              null,
    team        varchar(255)  default '' not null,
    thumbnail   varchar(1023) default '' not null,
    url         varchar(1023)            not null,
    likes       bigint        default 0  not null,
    index photo_create_time (create_time),
    index photo_team (team)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists post_category
(
    id          int auto_increment primary key,
    create_time datetime(6) not null,
    update_time datetime(6) null,
    category_id int         not null,
    post_id     int         not null,
    index post_category_category_id (category_id),
    index post_category_post_id (post_id)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists post_tag
(
    id          int auto_increment primary key,
    create_time datetime(6) not null,
    update_time datetime(6) null,
    post_id     int         not null,
    tag_id      int         not null,
    index post_tag_post_id (post_id),
    index post_tag_tag_id (tag_id)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists post
(
    id               int auto_increment primary key,
    type             int           default 0  not null,
    create_time      datetime(6)              not null,
    update_time      datetime(6)              null,
    disallow_comment tinyint(1)    default 0  not null,
    edit_time        datetime(6)              null,
    editor_type      int           default 0  not null,
    format_content   longtext                 not null,
    likes            bigint        default 0  not null,
    meta_description varchar(1023) default '' not null,
    meta_keywords    varchar(511)  default '' not null,
    original_content longtext                 not null,
    password         varchar(255)  default '' not null,
    slug             varchar(255)             not null,
    status           int           default 1  not null,
    summary          longtext                 not null,
    template         varchar(255)  default '' not null,
    thumbnail        varchar(1023) default '' not null,
    title            varchar(255)             not null,
    top_priority     int           default 0  not null,
    visits           bigint        default 0  not null,
    word_count       bigint        default 0  not null,
    unique index uniq_post_slug (slug),
    index post_create_time (create_time),
    index post_type_status (type, status)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

create table if not exists tag
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    name        varchar(255)             not null,
    slug        varchar(50)              not null,
    thumbnail   varchar(1023) default '' not null,
    color       varchar(25)   default '' not null,
    unique index uniq_tag_slug (slug),
    index tag_name (name)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists theme_setting
(
    id            int auto_increment primary key,
    create_time   datetime(6)  not null,
    update_time   datetime(6)  null,
    setting_key   varchar(255) not null,
    theme_id      varchar(255) not null,
    setting_value longtext     not null,
    index theme_setting_setting_key (setting_key),
    index theme_setting_theme_id (theme_id)
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;


create table if not exists user
(
    id          int auto_increment primary key,
    create_time datetime(6)              not null,
    update_time datetime(6)              null,
    avatar      varchar(1023) default '' not null,
    description varchar(1023) default '' not null,
    email       varchar(127)  default '' not null,
    expire_time datetime(6)              null,
    mfa_key     varchar(64)   default '' not null,
    mfa_type    int           default 0  not null,
    nickname    varchar(255)             not null,
    password    varchar(255)             not null,
    username    varchar(50)              not null
) ENGINE = INNODB
  DEFAULT charset = utf8mb4;

