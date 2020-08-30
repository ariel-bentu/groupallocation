update pupils set name='יותם דבי' where task=15 and id=89

select * from pupils where name ='יותם דבי'
select * from pupils where task=15 order by id desc

update pupils set name = 'יותם דבי' where id=87

insert into pupils values ('ariel', 15, 90, 'מתן לב רן', 1, 0)

select * from taskResult

 
select * from taskResultLines where resultId = 31 and pupilId >=80
select * from pupils where task=15 and id >=80
select * from taskResult
update pupils set gender = 1 where id=78 and task=15

delete from pupils where name = 'יותם דבי' and inactive is null


update taskResultLines set groupId=1 where resultId = 31 and pupilId=87
delete from taskResultLines where resultId=31 and pupilId=89

insert into taskResultLines values (31, 87, 2)

insert into pupils values ('ariel', 15, 87, 'ארי שוורץ', 1, 0)
alter table pupils add  remarks VARCHAR

alter table subgroups add  garden int
alter table pupilPrefs add  active smallint
