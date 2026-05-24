-- Polymorphic pivot: users/admins -> roles
CREATE TABLE IF NOT EXISTS `model_roles` (
  `model_id` bigint unsigned NOT NULL,
  `model_type` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`model_id`, `model_type`, `role_id`),
  KEY `idx_model_roles_model_id` (`model_id`),
  KEY `idx_model_roles_model_type` (`model_type`),
  KEY `idx_model_roles_role_id` (`role_id`),
  CONSTRAINT `fk_model_roles_role_id`
    FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
