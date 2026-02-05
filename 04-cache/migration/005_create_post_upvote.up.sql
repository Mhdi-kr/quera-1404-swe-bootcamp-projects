CREATE TABLE IF NOT EXISTS `user_post_upvote` (
    `user_id` INT NOT NULL,
    `post_id` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`, `post_id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `post`(`id`) ON DELETE CASCADE,
    INDEX `idx_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `user_comment_upvote` ADD `post_id` INT NOT NULL; 
ALTER TABLE `user_comment_upvote`
ADD CONSTRAINT FK_CommentPostID
FOREIGN KEY (`post_id`) REFERENCES `post`(`id`) ON DELETE CASCADE;