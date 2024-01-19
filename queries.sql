COPY (WITH cte as (select group_name, sum(going) as rsvp_23 from events where date > '2023-01-01' AND date <= '2023-12-31' GROUP BY group_name ORDER BY group_name ASC NULLS LAST),
cte2 as (select group_name, sum(going) as rsvp_22 from events where date > '2022-01-01' AND date <= '2022-12-31' GROUP BY group_name ORDER BY group_name ASC NULLS LAST)
select * from cte FULL OUTER JOIN cte2 ON cte.group_name = cte2.group_name) 
to 'meetup_rsvp.csv' (HEADER, DELIMITER ',');

COPy (WITH cte as (select group_name, sum(going) as rsvp_23, monthname(date) as mon from events where date > '2023-01-01' AND date <= '2023-12-31' GROUP by group_name, mon),
cte2 as (select group_name, sum(going) as rsvp_23, monthname(date) as mon from events where date > '2022-01-01' AND date <= '2022-12-31' GROUP by group_name, mon)
select * from cte FULL OUTER JOIN cte2 ON cte.group_name = cte2.group_name AND cte.mon = cte2.mon)
to 'meetup_rsvp_month.csv' (HEADER, DELIMITER ',');

