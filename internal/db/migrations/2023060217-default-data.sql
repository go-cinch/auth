-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO "t_action" (id, name, code, word, resource, menu, btn) VALUES
(2972307337314305,'All Permissions','SN2837AY','*','*','*','*'),
(2972307337379841,'Default Permissions','KHXK5JVL','default','POST|/auth/logout|/auth.v1.Auth/Logout
GET|/auth/info|/auth.v1.Auth/Info
POST|/auth/pwd|/auth.v1.Auth/Pwd','/dashboard/base
/user/index',''),
(2972307337445377,'Dashboard','2QKHTYEE','dashboard','','/dashboard/base',''),
(2972307337510913,'User Query','GRNA3NPV','user.read','GET|/auth/user/list|/auth.v1.Auth/FindUser','/system/user','system.user.read'),
(2972307337576449,'User Create','2LV9MDWB','user.create','POST|/auth/register|/auth.v1.Auth/Register','/system/user','system.user.create'),
(2972307337641985,'User Update','NME3CT5H','user.update','PUT,PATCH|/auth/user/update|/auth.v1.Auth/UpdateUser
GET|/auth/action|/auth.v1.Auth/FindAction','/system/user','system.user.update'),
(2972307337707521,'User Delete','EQH37R9C','user.delete','DELETE|/auth/user/delete|/auth.v1.Auth/DeleteUser','/system/user','system.user.delete'),
(2972307337773057,'User Group Query','V2HRXGW9','user.group.read','GET|/auth/user/group/list|/auth.v1.Auth/FindUserGroup','/system/group','system.user.group.read'),
(2972307337838593,'User Group Create','GGKPXAL6','user.group.create','POST|/auth/user/group/create|/auth.v1.Auth/CreateUserGroup','/system/group','system.user.group.create'),
(2972307337904129,'User Group Update','JM3TT968','user.group.update','PUT,PATCH|/auth/user/group/update|/auth.v1.Auth/UpdateUserGroup
GET|/auth/user/list|/auth.v1.Auth/FindUser
GET|/auth/action/list|/auth.v1.Auth/FindAction','/system/group','system.user.group.update'),
(2972307337969665,'User Group Delete','JE45TMPQ','user.group.delete','DELETE|/auth/user/group/delete|/auth.v1.Auth/DeleteUserGroup','/system/group','system.user.group.delete'),
(2972307338035201,'Role Query','AS2V9HND','role.read','GET|/auth/role/list|/auth.v1.Auth/FindRole','/system/role','system.role.read'),
(2972307338100737,'Role Create','88BA22VF','role.create','POST|/auth/role/create|/auth.v1.Auth/CreateRole','/system/role','system.role.create'),
(2972307338166273,'Role Update','GE5YBVDN','role.update','PUT,PATCH|/auth/role/update|/auth.v1.Auth/UpdateRole','/system/role','system.role.update'),
(2972307338231809,'Role Delete','AY6QE7QG','role.delete','DELETE|/auth/role/delete|/auth.v1.Auth/DeleteRole','/system/role','system.role.delete'),
(2972307338297345,'Action Query','42TMWNP3','action.read','GET|/auth/action/list|/auth.v1.Auth/FindAction','/system/action','system.action.read'),
(2972307338362881,'Action Create','SXPYFM3K','action.create','POST|/auth/action/create|/auth.v1.Auth/CreateAction','/system/action','system.action.create'),
(2972307338428417,'Action Update','8VCXMSCW','action.update','PUT,PATCH|/auth/action/update|/auth.v1.Auth/UpdateAction','/system/action','system.action.update'),
(2972307338493953,'Action Delete','86QSDSRL','action.delete','DELETE|/auth/action/delete|/auth.v1.Auth/DeleteAction','/system/action','system.action.delete'),
(2972307338559489,'Whitelist Query','ALX2LHB2','whitelist.read','GET|/auth/whitelist/list|/auth.v1.Auth/FindWhitelist','/system/whitelist','system.whitelist.read'),
(2972307338625025,'Whitelist Create','ALCARRQ8','whitelist.create','POST|/auth/whitelist/create|/auth.v1.Auth/CreateWhitelist','/system/whitelist','system.whitelist.create'),
(2972307338690561,'Whitelist Update','28FN73B3','whitelist.update','PUT,PATCH|/auth/whitelist/update|/auth.v1.Auth/UpdateWhitelist','/system/whitelist','system.whitelist.update'),
(2972307338756097,'Whitelist Delete','E8SN4T9K','whitelist.delete','DELETE|/auth/whitelist/delete|/auth.v1.Auth/DeleteWhitelist','/system/whitelist','system.whitelist.delete');

INSERT INTO "t_role" (id, name, word, action) VALUES
(2972307338821633,'Admin','admin','SN2837AY'),
(2972307338887169,'Guest','guest','2QKHTYEE');

INSERT INTO "t_user_group" (id, name, word, action) VALUES
(2972307338952705,'Read Only','readonly','GRNA3NPV,V2HRXGW9,AS2V9HND,42TMWNP3'),
(2972307339018241,'Read Write','write','GRNA3NPV,2LV9MDWB,NME3CT5H,EQH37R9C,V2HRXGW9,GGKPXAL6,JM3TT968,JE45TMPQ,AS2V9HND,88BA22VF,GE5YBVDN,AY6QE7QG,42TMWNP3,SXPYFM3K,8VCXMSCW,86QSDSRL'),
(2972307339083777,'No Delete','nodelete','GRNA3NPV,2LV9MDWB,NME3CT5H,V2HRXGW9,GGKPXAL6,JM3TT968,AS2V9HND,88BA22VF,GE5YBVDN,42TMWNP3,SXPYFM3K,8VCXMSCW');

INSERT INTO "t_user" (id, created_at, updated_at, role_id, username, code, password, platform) VALUES
(2972307339149313,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,2972307338821633,'super','89HEK28Y','$2a$10$TRT9yIpxi3LLgBnVrvktDOpxYUeSpq4cKDhuSDU8n16iXRPWkvmxG','pc'),
(2972307339214849,CURRENT_TIMESTAMP + INTERVAL '1 second',CURRENT_TIMESTAMP + INTERVAL '1 second',2972307338887169,'guest','4VPNKE6M','$2a$10$er8ILElzUu9m7n6DLWZaPeG8h6R2hyySGawvx4y7E/CXKYfvxKifW','pc'),
(2972307339280385,CURRENT_TIMESTAMP + INTERVAL '2 seconds',CURRENT_TIMESTAMP + INTERVAL '2 seconds',0,'readonly','EXP78RGH','$2a$10$a5pNKJGB3X1BScsEUkA6Yub184Q99SiNbxbftJsOG88liuIKlnxcW','pc'),
(2972307339345921,CURRENT_TIMESTAMP + INTERVAL '3 seconds',CURRENT_TIMESTAMP + INTERVAL '3 seconds',0,'write','6SHWH93V','$2a$10$C.9Zfx/D0n9tep8zXP4jUekz58ClC6Zrx.vMjwxHCNPB6Rblib//S','pc'),
(2972307339411457,CURRENT_TIMESTAMP + INTERVAL '4 seconds',CURRENT_TIMESTAMP + INTERVAL '4 seconds',0,'nodelete','JJHWJ9YJ','$2a$10$8SPpr/z.ukV4IvSVUIHVQOhKzY3Xfp9QJla5poW4/HgBeMxSviQ22','pc');

INSERT INTO "t_user_user_group_relation" (user_id, user_group_id) VALUES
(2972307339280385,2972307338952705),
(2972307339345921,2972307339018241),
(2972307339411457,2972307339083777);

INSERT INTO "t_whitelist" (id, category, resource) VALUES
(2972307339476993, 0, '/grpc.health.v1.Health/Check
/grpc.health.v1.Health/Watch'),
(2972307339542529, 1, '/grpc.health.v1.Health/Check
/grpc.health.v1.Health/Watch');

-- +migrate Down
TRUNCATE TABLE "t_action" CASCADE;
TRUNCATE TABLE "t_role" CASCADE;
TRUNCATE TABLE "t_user_group" CASCADE;
TRUNCATE TABLE "t_user" CASCADE;
TRUNCATE TABLE "t_user_user_group_relation" CASCADE;
TRUNCATE TABLE "t_whitelist" CASCADE;
