select count(*) from post;

-- 100,000 DAU (daily active user)
-- 20% users submit 80% of the new posts
-- avg post submission count is 5 per day
-- 20,000 * 5 per day = 100,000 post 
-- week = 700,000
-- year = 36,500,000

use case:
-- showing top 20 comments of the day per post
```sql
select
  comment.id as id,
  count(user_comment_upvote.comment_id) as vote_count,
  comment.user_id,
  comment.post_id,
  post.created_at as posted_at,
  comment.content
  from comment
  left join user_comment_upvote on comment.id = user_comment_upvote.comment_id
  join post on post.id = comment.post_id
  where post.created_at >= '2025-12-25 00:00:00' and post.created_at < '2025-12-26 00:00:00'
  group by comment.id
  order by vote_count DESC;
  ```

-- showing top 20 posts per day (upvoted)
** same as comments

-- show new posts for the day
```sql
select * from post where created_at >= '2025-12-25 00:00:00' and created_at < '2025-12-26 00:00:00'
```


- exact search
```sql
select * from post where description = "ea"
```

- contains search
```sql
select * from post where description LIKE "%html%"
```

- full-text search

-- search post (full-text search)
we have to use a search index residing along our database

