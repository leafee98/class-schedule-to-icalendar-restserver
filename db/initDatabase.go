package db

import (
	"fmt"
	"strings"

	// for the side effect of driver register
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
)

const initSql = `
drop database if exists %s;
create database %s character set='utf8mb4' collate='utf8mb4_unicode_ci';

use %s;
create table t_user (
	c_id integer primary key AUTO_INCREMENT,
	c_email varchar(64),
	c_nickname varchar(32),
	c_username varchar(32),
	c_password binary(32),
	c_bio varchar(300) default '',
	c_join_time datetime not null default now()
);

create table t_config (
	c_id integer primary key AUTO_INCREMENT,
	c_type tinyint,                 # 1-global, 2-lesson
	c_name varchar(64),
	c_content varchar(1024),
	c_format tinyint,				# 1-json, 2-toml
	c_owner_id integer,
	c_remark varchar(300),
	c_create_time datetime not null default now(), # default create time is now()
	c_modify_time datetime not null default now(), # default modify time is now()
	c_deleted bool default false,
	
	constraint foreign key (c_owner_id) references t_user (c_id)
);

create table t_config_share (
	c_id integer primary key AUTO_INCREMENT,
	# c_access_uid varchar(64) unique,
	c_config_id integer,
	c_create_time datetime not null default now(),
	c_remark varchar(300),
	c_deleted bool default false,
	
	constraint foreign key (c_config_id) references t_config (c_id)
);

create table t_user_favourite_config (
	c_id integer primary key AUTO_INCREMENT,
	c_user_id integer,
	c_config_share_id integer,
	c_create_time datetime not null default now(), # default create time is now()
	
	constraint foreign key (c_user_id) references t_user (c_id),
	constraint foreign key (c_config_share_id) references t_config_share (c_id)
);

create table t_plan (
	c_id integer primary key AUTO_INCREMENT,
	c_name varchar(64),
	c_owner_id integer,
	c_remark varchar(300),
	c_create_time datetime not null default now(), # default create time is now()
	c_modify_time datetime not null default now(), # default modify time is now()
	c_deleted bool default false,
	
	constraint foreign key (c_owner_id) references t_user (c_id)
);

create table t_plan_config_relation (
	c_id integer primary key AUTO_INCREMENT,
	c_plan_id integer,
	c_config_id integer,
	
	constraint foreign key (c_plan_id) references t_plan (c_id),
	constraint foreign key (c_config_id) references t_config (c_id)
);

create table t_plan_config_share_relation (
	c_id integer primary key AUTO_INCREMENT,
	c_plan_id integer,
	c_config_share_id integer,
	
	constraint foreign key (c_plan_id) references t_plan (c_id),
	constraint foreign key (c_config_share_id) references t_config_share(c_id)
);

create table t_plan_share (
	c_id integer primary key AUTO_INCREMENT,
	# c_access_uid varchar(64) unique, # length of UUID
	c_plan_id integer,
	c_create_time datetime not null default now(),
	c_remark varchar(300),
	c_deleted bool default false,
	
	constraint foreign key (c_plan_id) references t_plan (c_id)
);

create table t_user_favourite_plan (
	c_id integer primary key AUTO_INCREMENT,
	c_user_id integer,
	c_plan_share_id integer,
	c_create_time datetime not null default now(), # default create time is now()
	
	constraint foreign key (c_user_id) references t_user (c_id),
	constraint foreign key (c_plan_share_id) references t_plan_share (c_id)
);

create table t_login_token (
	c_id integer primary key AUTO_INCREMENT,
	c_user_id integer,
	c_token varchar(32), # token will be uuid string removed dashes
	c_expire_time datetime not null default date_add(now(), interval 3 day),
	
	constraint foreign key (c_user_id) references t_user (c_id)
);

create table t_plan_token (
	c_id integer primary key AUTO_INCREMENT,
	c_token varchar(32) unique, # token will be uuid string removed dashes
	c_plan_id integer,
	c_create_time datetime not null default now(),
	
	constraint foreign key (c_plan_id) references t_plan (c_id)
);


# auto delete expired token
create event auto_remove_expired_token 
	on schedule every 4 hour
	comment 'auto delete expired token'
	do
		delete from t_login_token where c_expire_time > now();

# update the c_modify_time to now() on table t_config
create trigger t_config_update_modify_time
	before update on t_config
	for each row
	set new.c_modify_time = now();

# update the c_modify_time to now() on table t_plan
create trigger t_plan_update_modify_time
	before update on t_plan
	for each row
	set new.c_modify_time = now();

# update the c_modify_time to now() on table t_plan when
# 1. add config to plan
# 2. remove config from plan
# 3. add config shared to plan
# 4. remove config shared from plan
create trigger t_plan_update_modify_time_relationship_insert
	after insert on t_plan_config_relation 
	for each row
	update t_plan set c_modify_time = now() where c_id = new.c_plan_id;
create trigger t_plan_update_modify_time_relationship_delete
	after delete on t_plan_config_relation 
	for each row
	update t_plan set c_modify_time = now() where c_id = old.c_plan_id;
create trigger t_plan_update_modify_time_share_relationship_insert
	after insert on t_plan_config_share_relation
	for each row
	update t_plan set c_modify_time = now() where c_id = new.c_plan_id;
create trigger t_plan_update_modify_time_share_relationship_delete
	after delete on t_plan_config_share_relation 
	for each row
	update t_plan set c_modify_time = now() where c_id = old.c_plan_id;
`

func getSqlCommand() string {
	return fmt.Sprintf(initSql, config.DatabaseName, config.DatabaseName, config.DatabaseName)
}

func InitDatabase() error {
	DB, err := sqlx.Connect("mysql", dsn(config.DatabaseUsername, config.DatabasePassword, config.DatabaseHost, ""))
	if err != nil {
		return err
	}

	trans, err := DB.Begin()
	for _, command := range strings.Split(getSqlCommand(), ";") {
		if len(strings.Trim(command, " \t\n;")) == 0 {
			continue
		}
		if _, err = trans.Exec(command); err != nil {
			logrus.Info(command)
			return err
		}
	}
	err = trans.Commit()
	return err
}
