DROP TABLE IF EXISTS `user_post_upvote`;

ALTER TABLE `user_comment_upvote`
DROP FOREIGN KEY FK_CommentPostID; 
ALTER TABLE `user_comment_upvote` DROP COLUMN `post_id`;