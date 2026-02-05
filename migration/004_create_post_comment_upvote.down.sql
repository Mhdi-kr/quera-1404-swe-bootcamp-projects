DROP TABLE IF EXISTS `user_comment_upvote`;
ALTER TABLE post
ADD `vote_count` INT DEFAULT 0;
 