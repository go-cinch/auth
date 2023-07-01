-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO `action` (`id`, `name`, `code`, `word`, `resource`, `menu`, `btn`) VALUES
(8926844486248562689,'全部权限','SN2837AY','*','*','*','*'),
(8926844486248579238,'默认权限','KHXK5JVL','default','POST|/auth/pwd|/auth.v1.Auth/Pwd\nGET|/auth/idempotent|/auth.v1.Auth/Idempotent\n/auth.v1.Auth/CheckIdempotent','',''),
(8929298412088590337,'首页','2QKHTYEE','dashboard','','/dashboard/base',''),
(8929306305215070209,'用户查询','GRNA3NPV','user.read','GET|/auth/user|/auth.v1.Auth/FindUser','/system/user','system.user.read'),
(8929306391416406017,'用户新增','2LV9MDWB','user.create','POST|/auth/register|/auth.v1.Auth/Register','/system/user','system.user.create'),
(8929306434114420737,'用户修改','NME3CT5H','user.update','PUT,PATCH|/auth/user/*|/auth.v1.Auth/UpdateUser\nGET|/auth/action|/auth.v1.Auth/FindAction','/system/user','system.user.update'),
(8929306468524490753,'用户删除','EQH37R9C','user.delete','DELETE|/auth/user/*|/auth.v1.Auth/DeleteUser','/system/user','system.user.delete'),
(8929306872301748225,'用户组查询','V2HRXGW9','user.group.read','GET|/auth/user/group|/auth.v1.Auth/FindUserGroup','/system/group','system.user.group.read'),
(8929306977897545729,'用户组新增','GGKPXAL6','user.group.create','POST|/auth/user/group|/auth.v1.Auth/CreateUserGroup','/system/group','system.user.group.create'),
(8929307003835121665,'用户组修改','JM3TT968','user.group.update','PUT,PATCH|/auth/user/group/*|/auth.v1.Auth/UpdateUserGroup\nGET|/auth/user|/auth.v1.Auth/FindUser\nGET|/auth/action|/auth.v1.Auth/FindAction','/system/group','system.user.group.update'),
(8929307052153503745,'用户组删除','JE45TMPQ','user.group.delete','DELETE|/auth/user/group|/auth.v1.Auth/DeleteUserGroup','/system/group','system.user.group.delete'),
(8929716482426798081,'角色查询','AS2V9HND','role.read','GET|/auth/role|/auth.v1.Auth/FindRole','/system/role','system.role.read'),
(8929716516635541505,'角色创建','88BA22VF','role.create','POST|/auth/role|/auth.v1.Auth/CreateRole','/system/role','system.role.create'),
(8929716548663246849,'角色修改','GE5YBVDN','role.update','PUT,PATCH|/auth/role/*|/auth.v1.Auth/UpdateRole','/system/role','system.role.update'),
(8929716593156423681,'角色删除','AY6QE7QG','role.delete','DELETE|/auth/role/*|/auth.v1.Auth/DeleteRole','/system/role','system.role.delete'),
(8929717994339172353,'行为查询','42TMWNP3','action.read','GET|/auth/action|/auth.v1.Auth/FindAction','/system/action','system.action.read'),
(8929718040577179649,'行为新增','SXPYFM3K','action.create','POST|/auth/action|/auth.v1.Auth/CreateAction','/system/action','system.action.create'),
(8929718085473009665,'行为修改','8VCXMSCW','action.update','PUT,PATCH|/auth/action/*|/auth.v1.Auth/UpdateAction','/system/action','system.action.update'),
(8929718119446872065,'行为删除','86QSDSRL','action.delete','DELETE|/auth/action/*|/auth.v1.Auth/DeleteAction','/system/action','system.action.delete'),
(8946576676020551681,'白名单查询','ALX2LHB2','whitelist.read','GET|/auth/whitelist|/auth.v1.Auth/FindWhitelist','/system/whitelist','system.whitelist.read'),
(8946576685583564801,'白名单新增','ALCARRQ8','whitelist.create','POST|/auth/whitelist|/auth.v1.Auth/CreateWhitelist','/system/whitelist','system.whitelist.create'),
(8946576696086102017,'白名单修改','28FN73B3','whitelist.update','PUT,PATCH|/auth/whitelist/*|/auth.v1.Auth/UpdateWhitelist','/system/whitelist','system.whitelist.update'),
(8946576704491487233,'白名单删除','E8SN4T9K','whitelist.delete','DELETE|/auth/whitelist/*|/auth.v1.Auth/DeleteWhitelist','/system/whitelist','system.whitelist.delete');

INSERT INTO `role` (`id`, `name`, `word`, `action`) VALUES
(8929718176338411521,'管理员','admin','SN2837AY'),
(8929721534264639489,'访客','guest','2QKHTYEE');

INSERT INTO `user_group` (`id`, `name`, `word`, `action`) VALUES
(8929306707314606081,'只读','readonly','GRNA3NPV,V2HRXGW9,AS2V9HND,42TMWNP3'),
(8929306725685657601,'读写','write','GRNA3NPV,2LV9MDWB,NME3CT5H,EQH37R9C,V2HRXGW9,GGKPXAL6,JM3TT968,JE45TMPQ,AS2V9HND,88BA22VF,GE5YBVDN,AY6QE7QG,42TMWNP3,SXPYFM3K,8VCXMSCW,86QSDSRL'),
(8929717758803836929,'不能删除','nodelete','GRNA3NPV,2LV9MDWB,NME3CT5H,V2HRXGW9,GGKPXAL6,JM3TT968,AS2V9HND,88BA22VF,GE5YBVDN,42TMWNP3,SXPYFM3K,8VCXMSCW');

INSERT INTO `user` (`id`, `created_at`, `updated_at`, `role_id`, `username`, `code`, `password`, `platform`) VALUES
(8929281526625992705,NOW(),NOW(),8929718176338411521,'super','89HEK28Y','$2a$10$TRT9yIpxi3LLgBnVrvktDOpxYUeSpq4cKDhuSDU8n16iXRPWkvmxG','pc'),
(8929298014988664833,NOW()+1,NOW()+1,8929721534264639489,'guest','4VPNKE6M','$2a$10$er8ILElzUu9m7n6DLWZaPeG8h6R2hyySGawvx4y7E/CXKYfvxKifW','pc'),
(8929306627069181953,NOW()+2,NOW()+2,0,'readonly','EXP78RGH','$2a$10$a5pNKJGB3X1BScsEUkA6Yub184Q99SiNbxbftJsOG88liuIKlnxcW','pc'),
(8929306650406289409,NOW()+3,NOW()+3,0,'write','6SHWH93V','$2a$10$C.9Zfx/D0n9tep8zXP4jUekz58ClC6Zrx.vMjwxHCNPB6Rblib//S','pc'),
(8929717570412478465,NOW()+4,NOW()+4,0,'nodelete','JJHWJ9YJ','$2a$10$8SPpr/z.ukV4IvSVUIHVQOhKzY3Xfp9QJla5poW4/HgBeMxSviQ22','pc');

INSERT INTO `user_user_group_relation` (`user_id`, `user_group_id`) VALUES
(8929306627069181953,8929306707314606081),
(8929306650406289409,8929306725685657601),
(8929717570412478465,8929717758803836929);

INSERT INTO `whitelist` (`id`, `category`, `resource`) VALUES
(8946576552976449537, 0, 'POST|/auth/login|/auth.v1.Auth/Login\nGET|/auth/status|/auth.v1.Auth/Status\nPOST|/auth/logout|/auth.v1.Auth/Logout\nGET|/auth/captcha|/auth.v1.Auth/Captcha\nPOST|/auth/register|/auth.v1.Auth/Register\nGET|/auth/info|/auth.v1.Auth/Info\n/grpc.health.v1.Health/Check\n/grpc.health.v1.Health/Watch'),
(8946576560039657473, 1, 'POST|/auth/login|/auth.v1.Auth/Login\nGET|/auth/status|/auth.v1.Auth/Status\nPOST|/auth/logout|/auth.v1.Auth/Logout\nGET|/auth/captcha|/auth.v1.Auth/Captcha\nPOST|/auth/register|/auth.v1.Auth/Register\n/auth.v1.Auth/GetUserByCode\n/auth.v1.Auth/CheckIdempotent\n/grpc.health.v1.Health/Check\n/grpc.health.v1.Health/Watch'),
(8946576567690067969, 2, 'POST|/auth/action|/auth.v1.Auth/CreateAction\nPOST|/auth/role|/auth.v1.Auth/CreateRole\nPOST|/auth/user/group|/auth.v1.Auth/CreateUserGroup');

-- +migrate Down
TRUNCATE TABLE `action`;
TRUNCATE TABLE `role`;
TRUNCATE TABLE `user_group`;
TRUNCATE TABLE `user`;
TRUNCATE TABLE `user_user_group_relation`;
TRUNCATE TABLE `whitelist`;
